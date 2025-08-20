package services

import (
	"canada-hires/repos"
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type CronService interface {
	Start(ctx context.Context)
	Stop()
	TriggerManualUpdate() error
}

type cronService struct {
	lmiaService LMIAService
	repo        repos.LMIARepository
	ticker      *time.Ticker
	stopChan    chan bool
	running     bool
	mu          sync.Mutex
}

func NewCronService(lmiaService LMIAService, repo repos.LMIARepository) CronService {
	return &cronService{
		lmiaService: lmiaService,
		repo:        repo,
		stopChan:    make(chan bool),
	}
}

func (c *cronService) Start(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		log.Info("Cron service is already running")
		return
	}

	log.Info("Starting LMIA cron service")
	c.running = true

	// Run quarterly (every 3 months)
	// For development/testing, you can change this to a shorter interval
	quarterlyInterval := 90 * 24 * time.Hour // 90 days
	c.ticker = time.NewTicker(quarterlyInterval)

	// Check if we need to run initial update
	go func() {
		// Check when the last successful update was
		lastUpdate, err := c.lmiaService.GetLatestUpdateStatus()
		if err != nil {
			log.Warn("Could not check last update status", "error", err)
		}
		
		// Only run initial update if:
		// 1. No previous update exists, OR
		// 2. Last update was more than 7 days ago, OR
		// 3. Last update failed
		shouldRunUpdate := false
		
		if lastUpdate == nil {
			log.Info("No previous LMIA update found, running initial update")
			shouldRunUpdate = true
		} else if lastUpdate.Status != "completed" {
			log.Info("Last LMIA update failed, retrying", "last_status", lastUpdate.Status)
			shouldRunUpdate = true
		} else {
			// Check if it's been more than 7 days since last successful update
			if lastUpdate.CompletedAt != nil {
				daysSinceUpdate := time.Since(*lastUpdate.CompletedAt).Hours() / 24
				if daysSinceUpdate > 7 {
					log.Info("Last LMIA update was more than 7 days ago, running update", "days_ago", int(daysSinceUpdate))
					shouldRunUpdate = true
				} else {
					log.Info("Recent LMIA update found, skipping initial update", "days_ago", int(daysSinceUpdate))
				}
			} else {
				// CompletedAt is nil, treat as needing update
				log.Info("Last LMIA update has no completion time, running update")
				shouldRunUpdate = true
			}
		}
		
		if shouldRunUpdate {
			log.Info("Running initial LMIA data update")
			err := c.lmiaService.RunFullUpdate()
			if err != nil {
				log.Error("Initial LMIA data update failed", "error", err)
			} else {
				log.Info("Initial LMIA data update completed successfully")
			}
		}
	}()

	// Start the ticker in a goroutine
	go func() {
		for {
			select {
			case <-c.ticker.C:
				log.Info("Running scheduled LMIA data update")
				err := c.lmiaService.RunFullUpdate()
				if err != nil {
					log.Error("Scheduled LMIA data update failed", "error", err)
				} else {
					log.Info("Scheduled LMIA data update completed successfully")
				}
			case <-c.stopChan:
				log.Info("Stopping LMIA cron service")
				return
			case <-ctx.Done():
				log.Info("Context cancelled, stopping LMIA cron service")
				return
			}
		}
	}()
}

func (c *cronService) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return
	}

	log.Info("Stopping LMIA cron service")
	c.running = false

	if c.ticker != nil {
		c.ticker.Stop()
	}

	close(c.stopChan)
}

func (c *cronService) TriggerManualUpdate() error {
	log.Info("Triggering manual LMIA data update")
	
	// Run in background to avoid blocking the HTTP request
	go func() {
		err := c.lmiaService.RunFullUpdate()
		if err != nil {
			log.Error("Manual LMIA data update failed", "error", err)
		} else {
			log.Info("Manual LMIA data update completed successfully")
		}
	}()

	return nil
}