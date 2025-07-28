package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/scraper"
	scraper_types "canada-hires/scraper-types"
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

// ScraperConfig holds configuration for scraping operation
type ScraperConfig struct {
	JobTitle    string
	Province    string
	Pages       int
	SaveToAPI   bool
}

type ScraperService interface {
	RunScraper(numberOfPages int) (*models.JobScrapingRun, error)
	RunScraperWithConfig(ctx context.Context, config ScraperConfig) (*models.JobScrapingRun, error)
}

type scraperService struct {
	jobRepo repos.JobBankRepository
	logger  *log.Logger
}

func NewScraperService(jobRepo repos.JobBankRepository, logger *log.Logger) ScraperService {
	return &scraperService{
		jobRepo: jobRepo,
		logger:  logger,
	}
}

func (s *scraperService) RunScraper(numberOfPages int) (*models.JobScrapingRun, error) {
	s.logger.Info("Starting job scraper", "pages", numberOfPages)

	// Create scraping run record
	scrapingRun := &models.JobScrapingRun{
		ID:        uuid.New().String(),
		Status:    "running",
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.jobRepo.CreateScrapingRun(scrapingRun); err != nil {
		return nil, fmt.Errorf("failed to create scraping run: %w", err)
	}

	// Initialize scraper
	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Failed to initialize scraper: %v", err))
		return nil, fmt.Errorf("failed to create scraper: %w", err)
	}
	defer scraperInstance.Close()

	// Scrape jobs
	jobs, err := scraperInstance.ScrapeLMIAJobs(numberOfPages)
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Scraping failed: %v", err))
		return nil, fmt.Errorf("failed to scrape jobs: %w", err)
	}

	s.logger.Info("Scraping completed", "jobs_found", len(jobs))

	// Convert scraper data to models
	scraperData := make([]models.ScraperJobData, len(jobs))
	for i, job := range jobs {
		scraperData[i] = convertToScraperJobData(job)
	}

	// Save jobs to database
	savedJobs, err := s.jobRepo.CreateJobPostingsFromScraperData(scraperData, scrapingRun.ID)
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Failed to save jobs: %v", err))
		return nil, fmt.Errorf("failed to save jobs to database: %w", err)
	}

	// Clean up orphaned jobs (jobs that existed in previous scrapes but not in current scrape)
	currentJobBankIDs := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if job.JobBankID != "" {
			currentJobBankIDs = append(currentJobBankIDs, job.JobBankID)
		}
	}

	deletedCount, err := s.jobRepo.DeleteJobPostingsNotInScrapeRun(scrapingRun.ID, currentJobBankIDs)
	if err != nil {
		s.logger.Error("Failed to clean up orphaned job postings", "error", err)
		// Don't fail the entire operation if cleanup fails, just log the error
	} else {
		s.logger.Info("Cleaned up orphaned job postings", "deleted_count", deletedCount)
	}

	// Update scraping run as completed
	if err := s.jobRepo.UpdateScrapingRunCompleted(scrapingRun.ID, numberOfPages, len(jobs), len(savedJobs)); err != nil {
		s.logger.Error("Failed to update scraping run completion", "error", err)
	}

	s.logger.Info("Scraping run completed successfully", 
		"run_id", scrapingRun.ID,
		"jobs_scraped", len(jobs),
		"jobs_saved", len(savedJobs),
		"orphaned_jobs_deleted", deletedCount)

	return scrapingRun, nil
}

func (s *scraperService) updateScrapingRunError(runID, errorMessage string) {
	if err := s.jobRepo.UpdateScrapingRunStatus(runID, "failed", &errorMessage); err != nil {
		s.logger.Error("Failed to update scraping run status", "error", err)
	}
}

func (s *scraperService) RunScraperWithConfig(ctx context.Context, config ScraperConfig) (*models.JobScrapingRun, error) {
	s.logger.Info("Starting job scraper with config", 
		"title", config.JobTitle, 
		"province", config.Province,
		"pages", config.Pages,
		"save_api", config.SaveToAPI)

	// If not saving to API, just run the basic scraper logic without database operations
	if !config.SaveToAPI {
		return s.runScraperSimple(ctx, config)
	}

	// Create scraping run record
	scrapingRun := &models.JobScrapingRun{
		ID:        uuid.New().String(),
		Status:    "running",
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.jobRepo.CreateScrapingRun(scrapingRun); err != nil {
		return nil, fmt.Errorf("failed to create scraping run: %w", err)
	}

	// Initialize scraper
	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Failed to initialize scraper: %v", err))
		return nil, fmt.Errorf("failed to create scraper: %w", err)
	}
	defer scraperInstance.Close()

	// Scrape jobs with config parameters (Note: scraper currently only supports numberOfPages)
	jobs, err := scraperInstance.ScrapeLMIAJobs(config.Pages)
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Scraping failed: %v", err))
		return nil, fmt.Errorf("failed to scrape jobs: %w", err)
	}

	s.logger.Info("Scraping completed", "jobs_found", len(jobs))

	// Convert scraper data to models
	scraperData := make([]models.ScraperJobData, len(jobs))
	for i, job := range jobs {
		scraperData[i] = convertToScraperJobData(job)
	}

	// Save jobs to database
	savedJobs, err := s.jobRepo.CreateJobPostingsFromScraperData(scraperData, scrapingRun.ID)
	if err != nil {
		s.updateScrapingRunError(scrapingRun.ID, fmt.Sprintf("Failed to save jobs: %v", err))
		return nil, fmt.Errorf("failed to save jobs to database: %w", err)
	}

	// Clean up orphaned jobs (jobs that existed in previous scrapes but not in current scrape)
	currentJobBankIDs := make([]string, 0, len(jobs))
	for _, job := range jobs {
		if job.JobBankID != "" {
			currentJobBankIDs = append(currentJobBankIDs, job.JobBankID)
		}
	}

	deletedCount, err := s.jobRepo.DeleteJobPostingsNotInScrapeRun(scrapingRun.ID, currentJobBankIDs)
	if err != nil {
		s.logger.Error("Failed to clean up orphaned job postings", "error", err)
		// Don't fail the entire operation if cleanup fails, just log the error
	} else {
		s.logger.Info("Cleaned up orphaned job postings", "deleted_count", deletedCount)
	}

	// Update scraping run as completed
	if err := s.jobRepo.UpdateScrapingRunCompleted(scrapingRun.ID, config.Pages, len(jobs), len(savedJobs)); err != nil {
		s.logger.Error("Failed to update scraping run completion", "error", err)
	}

	s.logger.Info("Scraping run completed successfully", 
		"run_id", scrapingRun.ID,
		"jobs_scraped", len(jobs),
		"jobs_saved", len(savedJobs),
		"orphaned_jobs_deleted", deletedCount)

	return scrapingRun, nil
}

func (s *scraperService) runScraperSimple(ctx context.Context, config ScraperConfig) (*models.JobScrapingRun, error) {
	s.logger.Info("Running scraper in simple mode (no database save)")

	// Initialize scraper
	scraperInstance, err := scraper.NewScraper()
	if err != nil {
		return nil, fmt.Errorf("failed to create scraper: %w", err)
	}
	defer scraperInstance.Close()

	// Scrape jobs with config parameters (Note: scraper currently only supports numberOfPages)
	jobs, err := scraperInstance.ScrapeLMIAJobs(config.Pages)
	if err != nil {
		return nil, fmt.Errorf("failed to scrape jobs: %w", err)
	}

	s.logger.Info("Scraping completed", "jobs_found", len(jobs))

	// Create a mock scraping run for return
	scrapingRun := &models.JobScrapingRun{
		ID:        uuid.New().String(),
		Status:    "completed",
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	return scrapingRun, nil
}

func convertToScraperJobData(job scraper_types.JobData) models.ScraperJobData {
	return models.ScraperJobData{
		JobTitle:  job.JobTitle,
		Business:  job.Business,
		Salary:    job.Salary,
		Location:  job.Location,
		JobUrl:    job.JobURL,
		Date:      job.Date,
		JobBankID: &job.JobBankID,
	}
}