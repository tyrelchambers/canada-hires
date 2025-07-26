package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"

	"github.com/charmbracelet/log"
)

type JobService interface {
	ProcessScraperData(scraperData []models.ScraperJobData, scrapingRunID string) error
	ValidateScraperData(scraperData []models.ScraperJobData) error
	GetJobStatistics() (*JobStatistics, error)
}

type JobStatistics struct {
	TotalJobs      int                      `json:"total_jobs"`
	TotalEmployers int                      `json:"total_employers"`
	TopEmployers   []map[string]interface{} `json:"top_employers"`
}

type jobService struct {
	jobBankRepo repos.JobBankRepository
}

func NewJobService(jobBankRepo repos.JobBankRepository) JobService {
	return &jobService{
		jobBankRepo: jobBankRepo,
	}
}

// ProcessScraperData processes and stores job data from the scraper
func (js *jobService) ProcessScraperData(scraperData []models.ScraperJobData, scrapingRunID string) error {
	if len(scraperData) == 0 {
		return fmt.Errorf("no scraper data provided")
	}

	// Validate the data first
	if err := js.ValidateScraperData(scraperData); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Store the job postings using the repository
	if err := js.jobBankRepo.CreateJobPostingsFromScraperData(scraperData, scrapingRunID); err != nil {
		return fmt.Errorf("failed to store job postings: %w", err)
	}

	log.Info("Successfully processed scraper data", "job_count", len(scraperData), "scraping_run_id", scrapingRunID)
	return nil
}

// ValidateScraperData validates the incoming scraper data
func (js *jobService) ValidateScraperData(scraperData []models.ScraperJobData) error {
	for i, data := range scraperData {
		if data.JobTitle == "" {
			return fmt.Errorf("job %d: job title is required", i)
		}
		if data.Business == "" {
			return fmt.Errorf("job %d: business name is required", i)
		}
		if data.JobUrl == "" {
			return fmt.Errorf("job %d: job URL is required", i)
		}
		if data.Location == "" {
			return fmt.Errorf("job %d: location is required", i)
		}
	}
	return nil
}

// GetJobStatistics returns comprehensive job statistics
func (js *jobService) GetJobStatistics() (*JobStatistics, error) {
	totalJobs, err := js.jobBankRepo.GetJobPostingsCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get job count: %w", err)
	}

	totalEmployers, err := js.jobBankRepo.GetDistinctEmployersCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get employer count: %w", err)
	}

	topEmployers, err := js.jobBankRepo.GetEmployerJobCounts(10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top employers: %w", err)
	}

	return &JobStatistics{
		TotalJobs:      totalJobs,
		TotalEmployers: totalEmployers,
		TopEmployers:   topEmployers,
	}, nil
}