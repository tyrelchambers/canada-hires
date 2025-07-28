package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type ScraperJobRepository interface {
	GetScraperJobByType(jobType string) (*models.ScraperJob, error)
	UpdateLastRunTime(jobType string, lastRunTime time.Time) error
	UpdateStatus(jobType string, status string) error
	UpdateNextScheduledRun(jobType string, nextRun time.Time) error
	CreateScraperJob(scraperJob *models.ScraperJob) error
}

type scraperJobRepository struct {
	db *sqlx.DB
}

func NewScraperJobRepository(db *sqlx.DB) ScraperJobRepository {
	return &scraperJobRepository{db: db}
}

func (r *scraperJobRepository) GetScraperJobByType(jobType string) (*models.ScraperJob, error) {
	var scraperJob models.ScraperJob
	
	query := `
		SELECT id, job_type, last_run_at, next_scheduled_run, status, created_at, updated_at
		FROM scraper_jobs 
		WHERE job_type = $1
	`
	
	err := r.db.Get(&scraperJob, query, jobType)
	if err != nil {
		return nil, fmt.Errorf("failed to get scraper job by type %s: %w", jobType, err)
	}
	
	return &scraperJob, nil
}

func (r *scraperJobRepository) UpdateLastRunTime(jobType string, lastRunTime time.Time) error {
	query := `
		UPDATE scraper_jobs 
		SET last_run_at = $1, updated_at = NOW()
		WHERE job_type = $2
	`
	
	result, err := r.db.Exec(query, lastRunTime, jobType)
	if err != nil {
		return fmt.Errorf("failed to update last run time for job type %s: %w", jobType, err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no scraper job found with type %s", jobType)
	}
	
	return nil
}

func (r *scraperJobRepository) UpdateStatus(jobType string, status string) error {
	query := `
		UPDATE scraper_jobs 
		SET status = $1, updated_at = NOW()
		WHERE job_type = $2
	`
	
	result, err := r.db.Exec(query, status, jobType)
	if err != nil {
		return fmt.Errorf("failed to update status for job type %s: %w", jobType, err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no scraper job found with type %s", jobType)
	}
	
	return nil
}

func (r *scraperJobRepository) UpdateNextScheduledRun(jobType string, nextRun time.Time) error {
	query := `
		UPDATE scraper_jobs 
		SET next_scheduled_run = $1, updated_at = NOW()
		WHERE job_type = $2
	`
	
	result, err := r.db.Exec(query, nextRun, jobType)
	if err != nil {
		return fmt.Errorf("failed to update next scheduled run for job type %s: %w", jobType, err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("no scraper job found with type %s", jobType)
	}
	
	return nil
}

func (r *scraperJobRepository) CreateScraperJob(scraperJob *models.ScraperJob) error {
	query := `
		INSERT INTO scraper_jobs (job_type, last_run_at, next_scheduled_run, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRow(
		query,
		scraperJob.JobType,
		scraperJob.LastRunAt,
		scraperJob.NextScheduledRun,
		scraperJob.Status,
	).Scan(&scraperJob.ID, &scraperJob.CreatedAt, &scraperJob.UpdatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create scraper job: %w", err)
	}
	
	return nil
}