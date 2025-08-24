package services

import (
	"canada-hires/repos"
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type NonCompliantCronService struct {
	cron                    *cron.Cron
	logger                  *log.Logger
	nonCompliantService     NonCompliantService
	scraperJobRepo          repos.ScraperJobRepository
	jobType                 string
}

func NewNonCompliantCronService(logger *log.Logger, nonCompliantService NonCompliantService, scraperJobRepo repos.ScraperJobRepository) *NonCompliantCronService {
	c := cron.New(cron.WithLocation(time.UTC))

	return &NonCompliantCronService{
		cron:                c,
		logger:              logger,
		nonCompliantService: nonCompliantService,
		scraperJobRepo:      scraperJobRepo,
		jobType:             "non_compliant_scraper",
	}
}

func (ncs *NonCompliantCronService) Start(ctx context.Context) error {
	// Check for missed execution on startup
	if err := ncs.checkMissedExecution(); err != nil {
		ncs.logger.Error("Failed to check missed execution", "error", err)
	}

	// Schedule daily execution at midnight UTC
	_, err := ncs.cron.AddFunc("0 0 * * *", ncs.runScraper)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	ncs.cron.Start()
	ncs.logger.Info("Non-compliant scraper cron service started - scheduled for daily execution at midnight UTC")

	// Keep the service running until context is cancelled
	<-ctx.Done()
	ncs.Stop()
	return nil
}

func (ncs *NonCompliantCronService) Stop() {
	if ncs.cron != nil {
		ncs.cron.Stop()
		ncs.logger.Info("Non-compliant scraper cron service stopped")
	}
}

func (ncs *NonCompliantCronService) runScraper() {
	ncs.logger.Info("Starting scheduled non-compliant scraper execution")

	// Update status to running
	if err := ncs.scraperJobRepo.UpdateStatus(ncs.jobType, "running"); err != nil {
		ncs.logger.Error("Failed to update non-compliant scraper status to running", "error", err)
	}

	if err := ncs.executeScraper(); err != nil {
		ncs.logger.Error("Non-compliant scraper execution failed", "error", err)
		// Update status to failed
		if updateErr := ncs.scraperJobRepo.UpdateStatus(ncs.jobType, "failed"); updateErr != nil {
			ncs.logger.Error("Failed to update non-compliant scraper status to failed", "error", updateErr)
		}
		return
	}

	// Update last run time and status
	now := time.Now()
	if err := ncs.scraperJobRepo.UpdateLastRunTime(ncs.jobType, now); err != nil {
		ncs.logger.Error("Failed to update last run time", "error", err)
	}
	
	if err := ncs.scraperJobRepo.UpdateStatus(ncs.jobType, "completed"); err != nil {
		ncs.logger.Error("Failed to update non-compliant scraper status to completed", "error", err)
	}

	// Set next scheduled run for tomorrow
	nextRun := now.Add(24 * time.Hour)
	if err := ncs.scraperJobRepo.UpdateNextScheduledRun(ncs.jobType, nextRun); err != nil {
		ncs.logger.Error("Failed to update next scheduled run", "error", err)
	}

	ncs.logger.Info("Non-compliant scraper execution completed successfully", "timestamp", now)
}

func (ncs *NonCompliantCronService) executeScraper() error {
	// Run the non-compliant scraper service
	cronJob, err := ncs.nonCompliantService.ScrapeAndStoreNonCompliantEmployers()
	if err != nil {
		return fmt.Errorf("non-compliant scraper execution failed: %w", err)
	}

	if cronJob.Status != "completed" {
		errorMsg := "Unknown error"
		if cronJob.ErrorMessage != nil {
			errorMsg = *cronJob.ErrorMessage
		}
		return fmt.Errorf("scraper failed with status %s: %s", cronJob.Status, errorMsg)
	}

	ncs.logger.Info("Non-compliant scraper execution completed",
		"status", cronJob.Status,
		"records_processed", cronJob.RecordsProcessed)
	return nil
}

func (ncs *NonCompliantCronService) checkMissedExecution() error {
	// Get scraper job from database
	scraperJob, err := ncs.scraperJobRepo.GetScraperJobByType(ncs.jobType)
	if err != nil {
		ncs.logger.Error("Failed to get non-compliant scraper job from database", "error", err)
		return err
	}

	// Check if we should run based on database record
	if scraperJob.ShouldRun() {
		if scraperJob.IsOverdue() {
			ncs.logger.Info("Detected overdue non-compliant scraper job, running catch-up scraper")
		} else {
			ncs.logger.Info("Detected non-compliant scraper job should run, starting execution")
		}
		go ncs.runScraper() // Run asynchronously to not block startup
	} else {
		ncs.logger.Info("Non-compliant scraper job does not need to run yet", 
			"last_run", scraperJob.LastRunAt,
			"should_run_at", scraperJob.LastRunAt.Add(24*time.Hour))
	}

	return nil
}

// GetLastRunTime returns the timestamp of the last successful scraper execution from database
func (ncs *NonCompliantCronService) GetLastRunTime() (*time.Time, error) {
	scraperJob, err := ncs.scraperJobRepo.GetScraperJobByType(ncs.jobType)
	if err != nil {
		return nil, err
	}
	return scraperJob.LastRunAt, nil
}

// RunNow manually triggers the scraper execution (useful for testing/admin)
func (ncs *NonCompliantCronService) RunNow() error {
	ncs.logger.Info("Manual non-compliant scraper execution triggered")
	return ncs.executeScraper()
}