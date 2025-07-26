package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

type JobBankBrowserService interface {
	ScrapeTFWJobsWithBrowser() error
}

type jobBankBrowserService struct {
	repo repos.JobBankRepository
}

func NewJobBankBrowserService(repo repos.JobBankRepository) JobBankBrowserService {
	return &jobBankBrowserService{
		repo: repo,
	}
}

func (s *jobBankBrowserService) ScrapeTFWJobsWithBrowser() error {
	log.Info("Starting TFW job scraping with browser automation")

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

	// Create chromedp context
	allocatorCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute) // 5 minute timeout for the whole process
	defer cancel()

	url := fmt.Sprintf("%s?fsrc=%s&sort=M", baseURL, tfwSourceParam)

	var jobsData []map[string]interface{}

	log.Info("Navigating to URL", "url", url)
		err = chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Info("Page loaded, waiting for initial content")
			return nil
		}),
		chromedp.Sleep(3*time.Second), // Wait for initial page load

		// Loop to click the "load more" button
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Info("Starting 'load more' loop")
			for {
				// Check if the button exists and is visible
				var buttonNodes []*cdp.Node
				err := chromedp.Nodes("#moreresultbutton", &buttonNodes, chromedp.AtLeast(0)).Do(ctx)
				if err != nil {
					return err
				}

				if len(buttonNodes) == 0 {
					log.Info("No 'more results' button found, assuming all jobs are loaded.")
					break
				}

				// Click the button
				err = chromedp.Click("#moreresultbutton", chromedp.NodeVisible).Do(ctx)
				if err != nil {
					// If the button is not clickable, it might be hidden or gone
					log.Warn("Could not click 'more results' button, assuming all jobs are loaded.", "error", err)
					break
				}

				log.Info("Clicked 'more results' button, waiting for new content...")
				// Wait for a bit for new results to load
				time.Sleep(2 * time.Second)
			}
			log.Info("Finished 'load more' loop")
			return nil
		}),

		// Scrape all the jobs now that they are loaded
		chromedp.ActionFunc(func(ctx context.Context) error {
			log.Info("Starting final scrape of all loaded jobs")
			return nil
		}),
		chromedp.Evaluate(`
			const jobs = [];
			const articles = document.querySelectorAll('article[id^="article-"]');

			articles.forEach((article) => {
				try {
					const jobId = article.id.replace('article-', '');
					const titleLink = article.querySelector('a[href*="/jobposting/"]');
					if (!titleLink) return;

					const url = titleLink.href;
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

	if err != nil {
		errorMsg := fmt.Sprintf("failed to scrape jobs: %v", err)
		s.repo.UpdateScrapingRunStatus(scrapingRun.ID, "failed", &errorMsg)
		return fmt.Errorf(errorMsg)
	}

	log.Info("Successfully scraped jobs", "count", len(jobsData))

	var allJobs []*models.JobPosting
	for _, jobData := range jobsData {
		job := convertBrowserJobToModel(jobData, scrapingRun.ID)
		if job != nil {
			allJobs = append(allJobs, job)
		}
	}

	jobsScraped := len(allJobs)
	jobsStored := 0

	// Store jobs in batches
	batchSize := 100
	for i := 0; i < len(allJobs); i += batchSize {
		end := i + batchSize
		if end > len(allJobs) {
			end = len(allJobs)
		}
		batch := allJobs[i:end]
		err = s.repo.CreateJobPostingsBatch(batch)
		if err != nil {
			log.Error("Failed to store job batch", "error", err)
		} else {
			jobsStored += len(batch)
		}
	}

	// Update scraping run as completed
	err = s.repo.UpdateScrapingRunCompleted(scrapingRun.ID, 1, jobsScraped, jobsStored)
	if err != nil {
		log.Error("Failed to update scraping run as completed", "error", err)
	}

	return nil
}

func convertBrowserJobToModel(jobData map[string]interface{}, scrapingRunID string) *models.JobPosting {
	title, _ := jobData["title"].(string)
	employer, _ := jobData["employer"].(string)
	location, _ := jobData["location"].(string)
	salaryText, _ := jobData["salaryText"].(string)
	jobURL, _ := jobData["url"].(string)
	jobID, _ := jobData["jobId"].(string)
	postedDate, _ := jobData["postedDate"].(string)
	lmiaFlag, _ := jobData["lmiaFlag"].(string)

	if title == "" || jobID == "" {
		return nil
	}

	// Parse location
	province, city := parseLocation(location)

	// Parse salary
	salaryMin, salaryMax, salaryType := parseSalary(salaryText)

	// Parse posted date
	var parsedDate *time.Time
	if postedDate != "" {
		if parsed, err := parsePostedDate(postedDate); err == nil {
			parsedDate = &parsed
		}
	}

	// Check if LMIA flag is present
	hasLMIA := lmiaFlag != ""

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
		PostingDate:   parsedDate,
		URL:           jobURL,
		IsTFW:         true,
		HasLMIA:       hasLMIA,
		ScrapingRunID: scrapingRunID,
	}
}
