package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/log"
)

type JobBankService interface {
	ScrapeTFWJobs() error
	ScrapeJobsFromPage(pageNum int, scrapingRunID string) ([]*models.JobPosting, error)
	ParseJobDetails(jobURL string) (*JobDetails, error)
	GetTotalJobCount() (int, error)
	GetScrapingStatus() (*models.JobScrapingRun, error)
}

type jobBankService struct {
	repo   repos.JobBankRepository
	client *http.Client
}

type JobDetails struct {
	Title       string
	Employer    string
	Location    string
	Province    string
	City        string
	SalaryMin   *float64
	SalaryMax   *float64
	SalaryType  string
	PostingDate *time.Time
	Description string
}

const (
	// Job Bank TFW search URL - fsrc=32 is the key parameter for TFW jobs
	baseURL           = "https://www.jobbank.gc.ca/jobsearch/jobsearch"
	tfwSourceParam    = "32"
	requestDelay      = 2 * time.Second // Rate limiting
	maxRetries        = 3
	jobsPerPage       = 25 // Job Bank shows 25 jobs per page
)

func NewJobBankService(repo repos.JobBankRepository) JobBankService {
	return &jobBankService{
		repo: repo,
		client: &http.Client{
			Timeout: 30 * time.Second,
			// Add User-Agent to be respectful
			Transport: &customTransport{
				Transport: http.DefaultTransport,
			},
		},
	}
}

type customTransport struct {
	Transport http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Canada-Hires Research Bot 1.0 - TFW Transparency Platform")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	return t.Transport.RoundTrip(req)
}

func (s *jobBankService) ScrapeTFWJobs() error {
	log.Info("Starting TFW job scraping process")

	// Create scraping run record
	scrapingRun := &models.JobScrapingRun{
		Status:    "running",
		StartedAt: time.Now(),
	}

	err := s.repo.CreateScrapingRun(scrapingRun)
	if err != nil {
		return fmt.Errorf("failed to create scraping run: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			errorMsg := fmt.Sprintf("Panic occurred: %v", r)
			s.repo.UpdateScrapingRunStatus(scrapingRun.ID, "failed", &errorMsg)
		}
	}()

	// First, get total job count to determine how many pages to scrape
	totalJobs, err := s.GetTotalJobCount()
	if err != nil {
		errorMsg := fmt.Sprintf("failed to get total job count: %v", err)
		s.repo.UpdateScrapingRunStatus(scrapingRun.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to get total job count: %w", err)
	}

	totalPages := (totalJobs + jobsPerPage - 1) / jobsPerPage
	log.Info("Total jobs and pages calculated", "total_jobs", totalJobs, "total_pages", totalPages)

	var allJobs []*models.JobPosting
	jobsScraped := 0
	jobsStored := 0

	// Scrape each page
	for page := 1; page <= totalPages; page++ {
		log.Info("Scraping page", "page", page, "total_pages", totalPages)

		// Rate limiting - be respectful to the server
		if page > 1 {
			time.Sleep(requestDelay)
		}

		jobs, err := s.ScrapeJobsFromPage(page, scrapingRun.ID)
		if err != nil {
			log.Error("Failed to scrape page", "page", page, "error", err)
			continue
		}

		jobsScraped += len(jobs)
		allJobs = append(allJobs, jobs...)

		// Store jobs in batches to avoid memory issues
		if len(allJobs) >= 100 {
			err = s.repo.CreateJobPostingsBatch(allJobs)
			if err != nil {
				log.Error("Failed to store job batch", "error", err)
			} else {
				jobsStored += len(allJobs)
			}
			allJobs = nil // Clear the batch
		}

		// Update progress
		err = s.repo.UpdateScrapingRunProgress(scrapingRun.ID, totalPages, jobsScraped, jobsStored, page)
		if err != nil {
			log.Error("Failed to update scraping progress", "error", err)
		}

		log.Info("Page completed", "page", page, "jobs_found", len(jobs), "total_scraped", jobsScraped)
	}

	// Store any remaining jobs
	if len(allJobs) > 0 {
		err = s.repo.CreateJobPostingsBatch(allJobs)
		if err != nil {
			log.Error("Failed to store final job batch", "error", err)
		} else {
			jobsStored += len(allJobs)
		}
	}

	// Update scraping run as completed
	err = s.repo.UpdateScrapingRunCompleted(scrapingRun.ID, totalPages, jobsScraped, jobsStored)
	if err != nil {
		log.Error("Failed to update scraping run as completed", "error", err)
	}

	log.Info("TFW job scraping completed", "total_pages", totalPages, "jobs_scraped", jobsScraped, "jobs_stored", jobsStored)
	return nil
}

func (s *jobBankService) GetTotalJobCount() (int, error) {
	// Build URL for first page to get total count
	params := url.Values{}
	params.Set("fsrc", tfwSourceParam)
	params.Set("page", "1")
	params.Set("sort", "M") // Sort by most recent

	fullURL := baseURL + "?" + params.Encode()

	resp, err := s.client.Get(fullURL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch job search page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Look for result count text - typically something like "4,394 jobs found"
	resultText := doc.Find(".results-count, .job-count, .search-results-count").First().Text()
	if resultText == "" {
		// Try alternative selectors
		resultText = doc.Find("h1").Text()
		if resultText == "" {
			// Look for any text containing "jobs"
			doc.Find("*").Each(func(i int, s *goquery.Selection) {
				text := s.Text()
				if strings.Contains(strings.ToLower(text), "jobs") && strings.Contains(text, "found") {
					resultText = text
					return
				}
			})
		}
	}

	// Parse the number from the result text
	re := regexp.MustCompile(`([\d,]+)\s*jobs?`)
	matches := re.FindStringSubmatch(resultText)
	if len(matches) < 2 {
		return 0, fmt.Errorf("could not find job count in page, result text: %s", resultText)
	}

	countStr := strings.ReplaceAll(matches[1], ",", "")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse job count: %w", err)
	}

	return count, nil
}

func (s *jobBankService) ScrapeJobsFromPage(pageNum int, scrapingRunID string) ([]*models.JobPosting, error) {
	// Build URL for specific page
	params := url.Values{}
	params.Set("fsrc", tfwSourceParam)
	params.Set("page", strconv.Itoa(pageNum))
	params.Set("sort", "M") // Sort by most recent

	fullURL := baseURL + "?" + params.Encode()

	var resp *http.Response
	var err error

	// Retry logic for resilience
	for retry := 0; retry < maxRetries; retry++ {
		resp, err = s.client.Get(fullURL)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}

		if retry < maxRetries-1 {
			time.Sleep(time.Duration(retry+1) * time.Second)
			log.Warn("Retrying page request", "page", pageNum, "retry", retry+1)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch page %d after %d retries: %w", pageNum, maxRetries, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code for page %d: %d", pageNum, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML for page %d: %w", pageNum, err)
	}

	var jobs []*models.JobPosting

	// Find job listings - Job Bank typically uses specific selectors for job cards
	doc.Find(".job-item, .job-posting, .job-card, article[role='article']").Each(func(i int, s *goquery.Selection) {
		job := parseJobFromListing(s, scrapingRunID)
		if job != nil {
			jobs = append(jobs, job)
		}
	})

	// If no jobs found with those selectors, try a more generic approach
	if len(jobs) == 0 {
		doc.Find("div").Each(func(i int, s *goquery.Selection) {
			// Look for divs that contain job-like structure
			if s.Find("a").Length() > 0 && (strings.Contains(s.Text(), "$") || strings.Contains(s.Text(), "/hour")) {
				job := parseJobFromListing(s, scrapingRunID)
				if job != nil {
					jobs = append(jobs, job)
				}
			}
		})
	}

	return jobs, nil
}

// Helper function to parse job from a selection
func parseJobFromListing(s *goquery.Selection, scrapingRunID string) *models.JobPosting {
	// Extract job title and link
	titleLink := s.Find("a, .job-title a, h3 a, h2 a").First()
	title := strings.TrimSpace(titleLink.Text())
	jobURL, exists := titleLink.Attr("href")

	if title == "" || !exists {
		return nil
	}

	// Make sure URL is absolute
	if strings.HasPrefix(jobURL, "/") {
		jobURL = "https://www.jobbank.gc.ca" + jobURL
	}

	// Extract job ID from URL
	jobID := extractJobIDFromURL(jobURL)
	if jobID == "" {
		return nil
	}

	// Extract employer
	employer := strings.TrimSpace(s.Find(".employer, .company, .job-employer").First().Text())
	if employer == "" {
		// Try to find employer in the text
		text := s.Text()
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != title && len(line) > 0 && !strings.Contains(line, "$") && !strings.Contains(line, "/") {
				employer = line
				break
			}
		}
	}

	// Extract location
	location := strings.TrimSpace(s.Find(".location, .job-location, .city").First().Text())

	// Parse location into city/province
	province, city := parseLocation(location)

	// Extract salary information
	salaryText := s.Find(".salary, .wage, .pay").First().Text()
	salaryMin, salaryMax, salaryType := parseSalary(salaryText)

	return &models.JobPosting{
		JobBankID:     &jobID,
		Title:         title,
		Employer:      employer,
		Location:      location,
		Province:      &province,
		City:          &city,
		SalaryMin:     salaryMin,
		SalaryMax:     salaryMax,
		SalaryType:    &salaryType,
		URL:           jobURL,
		IsTFW:         true,
		ScrapingRunID: scrapingRunID,
	}
}

func extractJobIDFromURL(jobURL string) string {
	// Job Bank URLs typically have job IDs in them
	re := regexp.MustCompile(`/jobposting/(\d+)`)
	matches := re.FindStringSubmatch(jobURL)
	if len(matches) > 1 {
		return matches[1]
	}

	// Alternative pattern
	re = regexp.MustCompile(`jobid=(\d+)`)
	matches = re.FindStringSubmatch(jobURL)
	if len(matches) > 1 {
		return matches[1]
	}

	// If no ID found, use hash of URL
	return fmt.Sprintf("url_%x", jobURL)
}

func parseLocation(location string) (province, city string) {
	if location == "" {
		return "", ""
	}

	// Canadian locations are typically "City, Province" or "City, XX"
	parts := strings.Split(location, ",")
	if len(parts) >= 2 {
		city = strings.TrimSpace(parts[0])
		province = strings.TrimSpace(parts[len(parts)-1])
	} else {
		city = strings.TrimSpace(location)
	}

	return province, city
}

func parseSalary(salaryText string) (min *float64, max *float64, salaryType string) {
	if salaryText == "" {
		return nil, nil, ""
	}

	// Determine salary type
	salaryType = "hourly" // default
	if strings.Contains(strings.ToLower(salaryText), "year") {
		salaryType = "yearly"
	} else if strings.Contains(strings.ToLower(salaryText), "month") {
		salaryType = "monthly"
	} else if strings.Contains(strings.ToLower(salaryText), "week") {
		salaryType = "weekly"
	}

	// Extract numeric values
	re := regexp.MustCompile(`\$?([\d,]+(?:\.\d{2})?)`)
	matches := re.FindAllStringSubmatch(salaryText, -1)

	var amounts []float64
	for _, match := range matches {
		if len(match) > 1 {
			amountStr := strings.ReplaceAll(match[1], ",", "")
			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
				amounts = append(amounts, amount)
			}
		}
	}

	if len(amounts) == 1 {
		min = &amounts[0]
		max = &amounts[0]
	} else if len(amounts) >= 2 {
		min = &amounts[0]
		max = &amounts[1]
		// Ensure min <= max
		if *min > *max {
			min, max = max, min
		}
	}

	return min, max, salaryType
}

// parsePostedDate parses posted date strings from Job Bank
func parsePostedDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	
	// Common Job Bank date formats
	formats := []string{
		"January 2, 2006",
		"Jan 2, 2006", 
		"January 02, 2006",
		"Jan 02, 2006",
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
	}
	
	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}
	
	// Try to parse relative dates like "Today", "Yesterday", etc.
	now := time.Now()
	switch strings.ToLower(dateStr) {
	case "today":
		return now, nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	}
	
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

func (s *jobBankService) ParseJobDetails(jobURL string) (*JobDetails, error) {
	// This could be used for getting more detailed information from individual job pages
	// For now, we'll implement basic structure
	resp, err := s.client.Get(jobURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch job details: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse job details HTML: %w", err)
	}

	// Parse detailed job information
	details := &JobDetails{
		Title:       strings.TrimSpace(doc.Find("h1, .job-title").First().Text()),
		Employer:    strings.TrimSpace(doc.Find(".employer, .company-name").First().Text()),
		Location:    strings.TrimSpace(doc.Find(".location, .job-location").First().Text()),
		Description: strings.TrimSpace(doc.Find(".job-description, .description").First().Text()),
	}

	details.Province, details.City = parseLocation(details.Location)

	return details, nil
}

func (s *jobBankService) GetScrapingStatus() (*models.JobScrapingRun, error) {
	return s.repo.GetLatestScrapingRun()
}