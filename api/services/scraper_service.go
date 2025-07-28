package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/scraper"
	scraper_types "canada-hires/scraper-types"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type ScraperService interface {
	RunScraper(numberOfPages int) (*models.JobScrapingRun, error)
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

	// Update scraping run as completed
	if err := s.jobRepo.UpdateScrapingRunCompleted(scrapingRun.ID, numberOfPages, len(jobs), len(savedJobs)); err != nil {
		s.logger.Error("Failed to update scraping run completion", "error", err)
	}

	s.logger.Info("Scraping run completed successfully", 
		"run_id", scrapingRun.ID,
		"jobs_scraped", len(jobs),
		"jobs_saved", len(savedJobs))

	return scrapingRun, nil
}

func (s *scraperService) updateScrapingRunError(runID, errorMessage string) {
	if err := s.jobRepo.UpdateScrapingRunStatus(runID, "failed", &errorMessage); err != nil {
		s.logger.Error("Failed to update scraping run status", "error", err)
	}
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