package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

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
	jobBankRepo   repos.JobBankRepository
	redditService RedditService
}

func NewJobService(jobBankRepo repos.JobBankRepository, redditService RedditService) JobService {
	return &jobService{
		jobBankRepo:   jobBankRepo,
		redditService: redditService,
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
	newJobPostings, err := js.jobBankRepo.CreateJobPostingsFromScraperData(scraperData, scrapingRunID)
	if err != nil {
		return fmt.Errorf("failed to store job postings: %w", err)
	}

	log.Info("Successfully processed scraper data", 
		"total_scraped", len(scraperData), 
		"new_jobs", len(newJobPostings),
		"scraping_run_id", scrapingRunID)

	// Check if automatic Reddit posting is enabled via environment variable
	if js.isAutoPostingEnabled() && js.redditService != nil && len(newJobPostings) > 0 {
		log.Info("Auto-posting enabled - posting new jobs to Reddit", "new_jobs", len(newJobPostings))
		go js.postJobsToReddit(newJobPostings)
	} else if len(newJobPostings) > 0 {
		log.Info("New jobs created and ready for admin review (auto-posting disabled)", "new_jobs", len(newJobPostings))
	} else {
		log.Info("No new jobs created - all jobs were updates to existing postings")
	}
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

// isAutoPostingEnabled checks if automatic Reddit posting is enabled via environment variable
// Defaults to false (disabled) for safety - set REDDIT_AUTO_POST=true to enable
func (js *jobService) isAutoPostingEnabled() bool {
	autoPost := os.Getenv("REDDIT_AUTO_POST")
	if autoPost == "" {
		return false // Default to disabled
	}
	
	enabled, err := strconv.ParseBool(autoPost)
	if err != nil {
		log.Warn("Invalid REDDIT_AUTO_POST value, defaulting to disabled", "value", autoPost)
		return false
	}
	
	return enabled
}

// postJobsToReddit posts new job postings to Reddit
func (js *jobService) postJobsToReddit(jobPostings []*models.JobPosting) {
	if len(jobPostings) == 0 {
		return
	}

	log.Info("Posting new jobs to Reddit", "job_count", len(jobPostings))

	for _, job := range jobPostings {
		if job == nil {
			continue
		}

		// Use async posting to avoid blocking
		go func(j *models.JobPosting) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if err := js.redditService.PostJob(ctx, j); err != nil {
				log.Error("Failed to post job to Reddit",
					"error", err,
					"job_id", j.ID,
					"job_title", j.Title,
					"employer", j.Employer,
				)
			} else {
				log.Info("Successfully posted job to Reddit",
					"job_id", j.ID,
					"job_title", j.Title,
					"employer", j.Employer,
				)
			}
		}(job)
	}
}