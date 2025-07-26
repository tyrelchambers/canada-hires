package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"context"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/charmbracelet/log"
)

type JobBankConcurrentService interface {
	ScrapeTFWJobsConcurrently() error
}

type jobBankConcurrentService struct {
	repo repos.JobBankRepository
}

type PageJob struct {
	PageNum       int
	ScrapingRunID string
}

type PageResult struct {
	PageNum int
	Jobs    []*models.JobPosting
	Error   error
}

func NewJobBankConcurrentService(repo repos.JobBankRepository) JobBankConcurrentService {
	return &jobBankConcurrentService{
		repo: repo,
	}
}

func (s *jobBankConcurrentService) ScrapeTFWJobsConcurrently() error {
	log.Info("Starting concurrent TFW job scraping")

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

	// Get total job count first
	totalJobs, err := s.getTotalJobCountConcurrent()
	if err != nil {
		errorMsg := fmt.Sprintf("failed to get total job count: %v", err)
		s.repo.UpdateScrapingRunStatus(scrapingRun.ID, "failed", &errorMsg)
		return fmt.Errorf("failed to get total job count: %w", err)
	}

	totalPages := (totalJobs + jobsPerPage - 1) / jobsPerPage
	log.Info("Total jobs and pages calculated", "total_jobs", totalJobs, "total_pages", totalPages)

	// Determine number of workers (use about 60% of cores, minimum 2, maximum 8)
	numCores := runtime.NumCPU()
	numWorkers := max(2, min(8, numCores*6/10))
	log.Info("Starting concurrent scraping", "workers", numWorkers, "cores", numCores)

	// Create channels
	pageJobs := make(chan PageJob, totalPages)
	results := make(chan PageResult, totalPages)

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go s.worker(i+1, pageJobs, results, &wg)
	}

	// Send all page jobs
	go func() {
		defer close(pageJobs)
		for page := 1; page <= totalPages; page++ {
			pageJobs <- PageJob{
				PageNum:       page,
				ScrapingRunID: scrapingRun.ID,
			}
		}
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	// Accumulate all results
	var allJobs []*models.JobPosting
	completedPages := 0
	errorCount := 0
	jobsScraped := 0

	fmt.Printf("=== CONCURRENT SCRAPING PROGRESS ===\n")
	
	for result := range results {
		completedPages++
		
		if result.Error != nil {
			log.Error("Page scraping failed", "page", result.PageNum, "error", result.Error)
			errorCount++
		} else {
			allJobs = append(allJobs, result.Jobs...)
			jobsScraped += len(result.Jobs)
		}

		// Progress update
		progress := float64(completedPages) / float64(totalPages) * 100
		fmt.Printf("\rProgress: %d/%d pages (%.1f%%) | Jobs: %d | Errors: %d", 
			completedPages, totalPages, progress, len(allJobs), errorCount)
		
		// Update database progress every 10 pages
		if completedPages%10 == 0 {
			s.repo.UpdateScrapingRunProgress(scrapingRun.ID, totalPages, jobsScraped, 0, completedPages)
		}
	}
	fmt.Printf("\n")

	// Now save all data to database in batches
	fmt.Printf("\n=== SAVING TO DATABASE ===\n")
	jobsStored := 0
	batchSize := 100
	
	for i := 0; i < len(allJobs); i += batchSize {
		end := min(i+batchSize, len(allJobs))
		batch := allJobs[i:end]
		
		err = s.repo.CreateJobPostingsBatch(batch)
		if err != nil {
			log.Error("Failed to store job batch", "batch_start", i, "error", err)
		} else {
			jobsStored += len(batch)
		}
		
		// Progress for saving
		saveProgress := float64(jobsStored) / float64(len(allJobs)) * 100
		fmt.Printf("\rSaving: %d/%d jobs (%.1f%%)", jobsStored, len(allJobs), saveProgress)
	}
	fmt.Printf("\n")

	// Update scraping run as completed
	err = s.repo.UpdateScrapingRunCompleted(scrapingRun.ID, totalPages, jobsScraped, jobsStored)
	if err != nil {
		log.Error("Failed to update scraping run as completed", "error", err)
	}

	// Print final summary
	s.printFinalSummary(scrapingRun.ID, totalPages, jobsScraped, jobsStored, errorCount)

	log.Info("Concurrent TFW job scraping completed", 
		"total_pages", totalPages, 
		"jobs_scraped", jobsScraped, 
		"jobs_stored", jobsStored,
		"errors", errorCount,
		"workers", numWorkers)
	
	return nil
}

func (s *jobBankConcurrentService) worker(workerID int, pageJobs <-chan PageJob, results chan<- PageResult, wg *sync.WaitGroup) {
	defer wg.Done()
	
	log.Debug("Worker started", "worker_id", workerID)

	for job := range pageJobs {
		// Rate limiting per worker - stagger workers
		time.Sleep(time.Duration(workerID) * 500 * time.Millisecond)
		
		log.Debug("Worker processing page", "worker_id", workerID, "page", job.PageNum)
		
		// Create fresh context for each page to avoid context pollution
		allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
			chromedp.Flag("headless", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocatorCtx)
		
		jobs, err := s.scrapePageConcurrent(ctx, job.PageNum, job.ScrapingRunID)
		
		// Always cancel context after use
		cancel()
		
		results <- PageResult{
			PageNum: job.PageNum,
			Jobs:    jobs,
			Error:   err,
		}
		
		// Small delay between requests from same worker
		time.Sleep(1 * time.Second)
	}
	
	log.Debug("Worker finished", "worker_id", workerID)
}

func (s *jobBankConcurrentService) getTotalJobCountConcurrent() (int, error) {
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	defer cancel()

	// Create chromedp context
	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var resultText string
	url := fmt.Sprintf("%s?fsrc=%s&page=1&sort=M", baseURL, tfwSourceParam)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`
			let resultText = '';
			const resultsSummary = document.querySelector('.results-summary');
			if (resultsSummary) {
				resultText = resultsSummary.textContent.trim();
			} else {
				const h2Elements = document.querySelectorAll('h2');
				for (let h2 of h2Elements) {
					if (h2.textContent.toLowerCase().includes('results')) {
						resultText = h2.textContent.trim();
						break;
					}
				}
			}
			resultText;
		`, &resultText),
	)

	if err != nil {
		return 0, fmt.Errorf("failed to get page content: %w", err)
	}

	if resultText == "" {
		return 0, fmt.Errorf("could not find job count text on page")
	}

	// Parse the number
	re := regexp.MustCompile(`([\d,]+)\s*results?`)
	matches := re.FindStringSubmatch(resultText)
	if len(matches) < 2 {
		re = regexp.MustCompile(`([\d,]+)`)
		matches = re.FindStringSubmatch(resultText)
		if len(matches) < 2 {
			return 0, fmt.Errorf("could not parse job count from text: %s", resultText)
		}
	}

	countStr := strings.ReplaceAll(matches[1], ",", "")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert count to integer: %w", err)
	}

	log.Info("Found total job count", "count", count)
	return count, nil
}

func (s *jobBankConcurrentService) scrapePageConcurrent(ctx context.Context, pageNum int, scrapingRunID string) ([]*models.JobPosting, error) {
	// Create timeout context for this page
	pageCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	url := fmt.Sprintf("%s?fsrc=%s&page=%d&sort=M", baseURL, tfwSourceParam, pageNum)

	var jobsData []map[string]interface{}
	var pageTitle string

	// Add retry logic for failed pages
	maxRetries := 2
	var lastErr error
	
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := chromedp.Run(pageCtx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`article[id^="article-"]`, chromedp.ByQuery),
			chromedp.Sleep(1*time.Second), // Small sleep after visible
			
			// Get page title for debugging
			chromedp.Title(&pageTitle),
			
			chromedp.Evaluate(`
				const jobs = [];
				const articles = document.querySelectorAll('article[id^="article-"]');
				
				articles.forEach((article) => {
					try {
						// Get job ID from article ID
						const jobId = article.id.replace('article-', '');
						
						// Find title link - it's the main job posting link
						const titleLink = article.querySelector('a[href*="/jobposting/"]');
						if (!titleLink) return;
						
						const url = titleLink.href;
						
						// Extract job details using class names
						const titleEl = article.querySelector('span.no-wrap[property="title"]');
						const title = titleEl ? titleEl.textContent.trim() : '';
						
						const businessEl = article.querySelector('li.employer');
						const employer = businessEl ? businessEl.textContent.trim() : '';
						
						const locationEl = article.querySelector('li.location');
						const location = locationEl ? locationEl.textContent.trim() : '';
						
						const salaryEl = article.querySelector('li.salary');
						const salaryText = salaryEl ? salaryEl.textContent.trim() : '';
						
						const dateEl = article.querySelector('.date');
						const postedDate = dateEl ? dateEl.textContent.trim() : '';
						
						const lmiaEl = article.querySelector('.jobLMIAflag');
						const lmiaFlag = lmiaEl ? lmiaEl.textContent.trim() : '';
						
						if (title && jobId && title.length > 2) {
							jobs.push({
								jobId: jobId,
								title: title,
								employer: employer || 'Unknown',
								location: location || '',
								salaryText: salaryText || '',
								postedDate: postedDate || '',
								lmiaFlag: lmiaFlag || '',
								url: url
							});
						}
					} catch (e) {
						console.log('Error processing article:', e);
					}
				});
				
				jobs;
			`, &jobsData),
		)

		if err == nil && len(jobsData) > 0 {
			lastErr = nil
			break // Success
		}

		lastErr = err
		if attempt < maxRetries {
			var retryReason error
			if err != nil {
				retryReason = err
			} else {
				retryReason = fmt.Errorf("no jobs found on page (title: %s)", pageTitle)
			}
			log.Warn("Page scraping failed, retrying", "page", pageNum, "attempt", attempt, "reason", retryReason)
			time.Sleep(time.Duration(attempt) * 2 * time.Second) // Exponential backoff
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("failed to scrape page %d after %d attempts (title: %s): %w", pageNum, maxRetries, pageTitle, lastErr)
	}

	if len(jobsData) == 0 {
		log.Warn("No jobs found on page", "page", pageNum, "title", pageTitle)
		return []*models.JobPosting{}, nil // Return empty slice instead of error
	}

	// Convert to job models
	var jobs []*models.JobPosting
	for _, jobData := range jobsData {
		job := convertBrowserJobToModel(jobData, scrapingRunID)
		if job != nil {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (s *jobBankConcurrentService) printFinalSummary(scrapingRunID string, totalPages, jobsScraped, jobsStored, errorCount int) {
	fmt.Printf("\n=== CONCURRENT SCRAPING COMPLETED ===\n")
	fmt.Printf("Scraping Run ID: %s\n", scrapingRunID)
	fmt.Printf("Total Pages Scraped: %d\n", totalPages)
	fmt.Printf("Total Jobs Found: %d\n", jobsScraped)
	fmt.Printf("Total Jobs Stored: %d\n", jobsStored)
	fmt.Printf("Errors: %d\n", errorCount)
	fmt.Printf("Success Rate: %.1f%%\n", float64(jobsStored)/float64(jobsScraped)*100)
	
	// Get sample data
	recentJobs, err := s.repo.GetJobPostingsByScrapingRun(scrapingRunID)
	if err == nil && len(recentJobs) > 0 {
		fmt.Printf("\n=== SAMPLE JOBS (First 3) ===\n")
		for i, job := range recentJobs[:min(3, len(recentJobs))] {
			fmt.Printf("\nJob %d:\n", i+1)
			fmt.Printf("  ID: %s\n", job.JobBankID)
			fmt.Printf("  Title: %s\n", job.Title)
			fmt.Printf("  Employer: %s\n", job.Employer)
			fmt.Printf("  Location: %s\n", job.Location)
			if job.SalaryMin != nil && job.SalaryMax != nil {
				fmt.Printf("  Salary: $%.2f - $%.2f", *job.SalaryMin, *job.SalaryMax)
				if job.SalaryType != nil {
					fmt.Printf(" (%s)", *job.SalaryType)
				}
				fmt.Printf("\n")
			}
		}
	}

	// Get top employers
	employerCounts, err := s.repo.GetEmployerJobCounts(5)
	if err == nil && len(employerCounts) > 0 {
		fmt.Printf("\n=== TOP EMPLOYERS ===\n")
		for i, emp := range employerCounts {
			fmt.Printf("%d. %s: %v jobs\n", i+1, emp["employer"], emp["job_count"])
		}
	}

	fmt.Printf("===================================\n\n")
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}