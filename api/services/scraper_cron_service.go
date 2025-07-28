package services

import (
	"canada-hires/repos"
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type ScraperCronService struct {
	cron              *cron.Cron
	logger            *log.Logger
	scraperService    ScraperService
	scraperJobRepo    repos.ScraperJobRepository
	jobType           string
}

func NewScraperCronService(logger *log.Logger, scraperService ScraperService, scraperJobRepo repos.ScraperJobRepository) *ScraperCronService {
	c := cron.New(cron.WithLocation(time.UTC))

	return &ScraperCronService{
		cron:           c,
		logger:         logger,
		scraperService: scraperService,
		scraperJobRepo: scraperJobRepo,
		jobType:        "lmia_scraper",
	}
}

func (scs *ScraperCronService) Start(ctx context.Context) error {
	// Check for missed execution on startup
	if err := scs.checkMissedExecution(); err != nil {
		scs.logger.Error("Failed to check missed execution", "error", err)
	}

	// Schedule daily execution at midnight UTC
	_, err := scs.cron.AddFunc("0 0 * * *", scs.runScraper)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	scs.cron.Start()
	scs.logger.Info("Scraper cron service started - scheduled for daily execution at midnight UTC")

	// Keep the service running until context is cancelled
	<-ctx.Done()
	scs.Stop()
	return nil
}

func (scs *ScraperCronService) Stop() {
	if scs.cron != nil {
		scs.cron.Stop()
		scs.logger.Info("Scraper cron service stopped")
	}
}

func (scs *ScraperCronService) runScraper() {
	scs.logger.Info("Starting scheduled scraper execution")

	// Update status to running
	if err := scs.scraperJobRepo.UpdateStatus(scs.jobType, "running"); err != nil {
		scs.logger.Error("Failed to update scraper status to running", "error", err)
	}

	if err := scs.executeScraper(); err != nil {
		scs.logger.Error("Scraper execution failed", "error", err)
		// Update status to failed
		if updateErr := scs.scraperJobRepo.UpdateStatus(scs.jobType, "failed"); updateErr != nil {
			scs.logger.Error("Failed to update scraper status to failed", "error", updateErr)
		}
		return
	}

	// Update last run time and status
	now := time.Now()
	if err := scs.scraperJobRepo.UpdateLastRunTime(scs.jobType, now); err != nil {
		scs.logger.Error("Failed to update last run time", "error", err)
	}
	
	if err := scs.scraperJobRepo.UpdateStatus(scs.jobType, "completed"); err != nil {
		scs.logger.Error("Failed to update scraper status to completed", "error", err)
	}

	// Set next scheduled run for tomorrow
	nextRun := now.Add(24 * time.Hour)
	if err := scs.scraperJobRepo.UpdateNextScheduledRun(scs.jobType, nextRun); err != nil {
		scs.logger.Error("Failed to update next scheduled run", "error", err)
	}

	scs.logger.Info("Scraper execution completed successfully", "timestamp", now)
}

func (scs *ScraperCronService) executeScraper() error {
	// Run the integrated scraper service
	scrapingRun, err := scs.scraperService.RunScraper(-1) // -1 means scrape all pages
	if err != nil {
		return fmt.Errorf("scraper execution failed: %w", err)
	}

	scs.logger.Info("Scraper execution completed",
		"run_id", scrapingRun.ID,
		"jobs_scraped", scrapingRun.JobsScraped,
		"jobs_stored", scrapingRun.JobsStored)
	return nil
}

func (scs *ScraperCronService) checkMissedExecution() error {
	// Get scraper job from database
	scraperJob, err := scs.scraperJobRepo.GetScraperJobByType(scs.jobType)
	if err != nil {
		scs.logger.Error("Failed to get scraper job from database", "error", err)
		return err
	}

	// Check if we should run based on database record
	if scraperJob.ShouldRun() {
		if scraperJob.IsOverdue() {
			scs.logger.Info("Detected overdue scraper job, running catch-up scraper")
		} else {
			scs.logger.Info("Detected scraper job should run, starting execution")
		}
		go scs.runScraper() // Run asynchronously to not block startup
	} else {
		scs.logger.Info("Scraper job does not need to run yet", 
			"last_run", scraperJob.LastRunAt,
			"should_run_at", scraperJob.LastRunAt.Add(24*time.Hour))
	}

	return nil
}

// GetLastRunTime returns the timestamp of the last successful scraper execution from database
func (scs *ScraperCronService) GetLastRunTime() (*time.Time, error) {
	scraperJob, err := scs.scraperJobRepo.GetScraperJobByType(scs.jobType)
	if err != nil {
		return nil, err
	}
	return scraperJob.LastRunAt, nil
}

// RunNow manually triggers the scraper execution (useful for testing/admin)
func (scs *ScraperCronService) RunNow() error {
	scs.logger.Info("Manual scraper execution triggered")
	return scs.executeScraper()
}
