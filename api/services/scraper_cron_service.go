package services

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/robfig/cron/v3"
)

type ScraperCronService struct {
	cron           *cron.Cron
	logger         *log.Logger
	scraperService ScraperService
	lastRun        time.Time
}

func NewScraperCronService(logger *log.Logger, scraperService ScraperService) *ScraperCronService {
	c := cron.New(cron.WithLocation(time.UTC))

	return &ScraperCronService{
		cron:           c,
		logger:         logger,
		scraperService: scraperService,
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

	if err := scs.executeScraper(); err != nil {
		scs.logger.Error("Scraper execution failed", "error", err)
		return
	}

	scs.lastRun = time.Now()
	scs.logger.Info("Scraper execution completed successfully", "timestamp", scs.lastRun)
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
	// Check if we missed yesterday's execution
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// If it's after midnight and we haven't run today, check if we missed yesterday
	if now.Hour() > 0 && scs.shouldRunCatchup(yesterday) {
		scs.logger.Info("Detected missed execution, running catch-up scraper")
		go scs.runScraper() // Run asynchronously to not block startup
	}

	return nil
}

func (scs *ScraperCronService) shouldRunCatchup(targetDate time.Time) bool {
	// Check if we have a record of running on the target date
	// For simplicity, we'll check if lastRun is from the target date
	// In production, you might want to store this in the database

	if scs.lastRun.IsZero() {
		// No previous run recorded, assume we should catch up
		return true
	}

	// Check if lastRun was on or after the target date
	targetStart := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, time.UTC)
	targetEnd := targetStart.Add(24 * time.Hour)

	return scs.lastRun.Before(targetStart) || scs.lastRun.After(targetEnd)
}

// GetLastRunTime returns the timestamp of the last successful scraper execution
func (scs *ScraperCronService) GetLastRunTime() time.Time {
	return scs.lastRun
}

// RunNow manually triggers the scraper execution (useful for testing/admin)
func (scs *ScraperCronService) RunNow() error {
	scs.logger.Info("Manual scraper execution triggered")
	return scs.executeScraper()
}
