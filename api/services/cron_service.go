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

	// Run initial update in background
	go func() {
		log.Info("Running initial LMIA data update")
		err := c.lmiaService.RunFullUpdate()
		if err != nil {
			log.Error("Initial LMIA data update failed", "error", err)
		} else {
			log.Info("Initial LMIA data update completed successfully")
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