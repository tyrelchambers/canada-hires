package models

import "time"

type ScraperJob struct {
	ID                int        `json:"id" db:"id"`
	JobType           string     `json:"job_type" db:"job_type"`
	LastRunAt         *time.Time `json:"last_run_at" db:"last_run_at"`
	NextScheduledRun  *time.Time `json:"next_scheduled_run" db:"next_scheduled_run"`
	Status            string     `json:"status" db:"status"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

// ShouldRun checks if the scraper job should run based on last run time
func (sj *ScraperJob) ShouldRun() bool {
	if sj.LastRunAt == nil {
		// Never run before, should run
		return true
	}
	
	// Check if it's been more than 24 hours since last run
	return time.Since(*sj.LastRunAt) > 24*time.Hour
}

// IsOverdue checks if the scraper is significantly overdue (more than 25 hours)
func (sj *ScraperJob) IsOverdue() bool {
	if sj.LastRunAt == nil {
		return true
	}
	
	return time.Since(*sj.LastRunAt) > 25*time.Hour
}