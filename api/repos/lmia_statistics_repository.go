package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LMIAStatisticsRepository interface {
	// Create and Update
	CreateStatistics(stats *models.LMIAStatistics) error
	UpdateStatistics(stats *models.LMIAStatistics) error
	UpsertStatistics(stats *models.LMIAStatistics) error

	// Get by filters
	GetStatisticsByDateRange(startDate, endDate time.Time, periodType models.PeriodType) ([]*models.LMIAStatistics, error)
	GetStatisticsByDate(date time.Time, periodType models.PeriodType) (*models.LMIAStatistics, error)
	GetLatestStatistics(periodType models.PeriodType, limit int) ([]*models.LMIAStatistics, error)

	// Aggregation helpers - get raw data from job_postings for statistics calculation
	GetJobStatisticsForDate(date time.Time) (*models.JobStatisticsData, error)
	GetJobStatisticsForMonth(year int, month int) (*models.JobStatisticsData, error)
	GetAllDatesWithJobs() ([]time.Time, error)

	// Utility
	DeleteStatistics(id string) error
	GetStatisticsByID(id string) (*models.LMIAStatistics, error)
}

type lmiaStatisticsRepository struct {
	db *sqlx.DB
}

func NewLMIAStatisticsRepository(db *sqlx.DB) LMIAStatisticsRepository {
	return &lmiaStatisticsRepository{db: db}
}

// CreateStatistics creates a new statistics record
func (r *lmiaStatisticsRepository) CreateStatistics(stats *models.LMIAStatistics) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lmia_job_statistics (id, date, period_type, total_jobs, unique_employers, 
										avg_salary_min, avg_salary_max, top_provinces, top_cities, 
										created_at, updated_at)
		VALUES (:id, :date, :period_type, :total_jobs, :unique_employers,
				:avg_salary_min, :avg_salary_max, :top_provinces, :top_cities,
				:created_at, :updated_at)
	`

	stats.ID = uuid.New().String()
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	_, err = tx.NamedExec(query, stats)
	if err != nil {
		return fmt.Errorf("failed to insert statistics: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateStatistics updates an existing statistics record
func (r *lmiaStatisticsRepository) UpdateStatistics(stats *models.LMIAStatistics) error {
	query := `
		UPDATE lmia_job_statistics 
		SET total_jobs = :total_jobs, unique_employers = :unique_employers,
			avg_salary_min = :avg_salary_min, avg_salary_max = :avg_salary_max,
			top_provinces = :top_provinces, top_cities = :top_cities,
			updated_at = :updated_at
		WHERE id = :id
	`

	stats.UpdatedAt = time.Now()

	_, err := r.db.NamedExec(query, stats)
	if err != nil {
		return fmt.Errorf("failed to update statistics: %w", err)
	}

	return nil
}

// UpsertStatistics inserts or updates statistics for a given date and period
func (r *lmiaStatisticsRepository) UpsertStatistics(stats *models.LMIAStatistics) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lmia_job_statistics (id, date, period_type, total_jobs, unique_employers, 
										avg_salary_min, avg_salary_max, top_provinces, top_cities, 
										created_at, updated_at)
		VALUES (:id, :date, :period_type, :total_jobs, :unique_employers,
				:avg_salary_min, :avg_salary_max, :top_provinces, :top_cities,
				:created_at, :updated_at)
		ON CONFLICT (date, period_type) DO UPDATE SET
			total_jobs = EXCLUDED.total_jobs,
			unique_employers = EXCLUDED.unique_employers,
			avg_salary_min = EXCLUDED.avg_salary_min,
			avg_salary_max = EXCLUDED.avg_salary_max,
			top_provinces = EXCLUDED.top_provinces,
			top_cities = EXCLUDED.top_cities,
			updated_at = EXCLUDED.updated_at
	`

	if stats.ID == "" {
		stats.ID = uuid.New().String()
	}
	stats.CreatedAt = time.Now()
	stats.UpdatedAt = time.Now()

	_, err = tx.NamedExec(query, stats)
	if err != nil {
		return fmt.Errorf("failed to upsert statistics: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetStatisticsByDateRange retrieves statistics within a date range
func (r *lmiaStatisticsRepository) GetStatisticsByDateRange(startDate, endDate time.Time, periodType models.PeriodType) ([]*models.LMIAStatistics, error) {
	var stats []*models.LMIAStatistics
	query := `
		SELECT * FROM lmia_job_statistics
		WHERE date >= $1 AND date <= $2 AND period_type = $3
		ORDER BY date ASC
	`

	err := r.db.Select(&stats, query, startDate, endDate, periodType)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics by date range: %w", err)
	}

	return stats, nil
}

// GetStatisticsByDate retrieves statistics for a specific date and period
func (r *lmiaStatisticsRepository) GetStatisticsByDate(date time.Time, periodType models.PeriodType) (*models.LMIAStatistics, error) {
	var stats models.LMIAStatistics
	query := `SELECT * FROM lmia_job_statistics WHERE date = $1 AND period_type = $2`

	err := r.db.Get(&stats, query, date, periodType)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetLatestStatistics retrieves the most recent statistics
func (r *lmiaStatisticsRepository) GetLatestStatistics(periodType models.PeriodType, limit int) ([]*models.LMIAStatistics, error) {
	var stats []*models.LMIAStatistics
	query := `
		SELECT * FROM lmia_job_statistics
		WHERE period_type = $1
		ORDER BY date DESC
		LIMIT $2
	`

	err := r.db.Select(&stats, query, periodType, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest statistics: %w", err)
	}

	return stats, nil
}

// GetJobStatisticsForDate aggregates job posting data for a specific date
func (r *lmiaStatisticsRepository) GetJobStatisticsForDate(date time.Time) (*models.JobStatisticsData, error) {
	// Set time to start and end of day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// Get basic aggregated data
	query := `
		SELECT 
			COUNT(*) as total_jobs,
			COUNT(DISTINCT employer) as unique_employers,
			AVG(salary_min) as avg_salary_min,
			AVG(salary_max) as avg_salary_max
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true
	`

	var data models.JobStatisticsData
	row := r.db.QueryRow(query, startOfDay, endOfDay)
	err := row.Scan(&data.TotalJobs, &data.UniqueEmployers, &data.AvgSalaryMin, &data.AvgSalaryMax)
	if err != nil {
		return nil, fmt.Errorf("failed to get job statistics for date: %w", err)
	}

	// Get province counts
	data.ProvincesCounts = make(map[string]int)
	provinceQuery := `
		SELECT province, COUNT(*) as count
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true AND province IS NOT NULL
		GROUP BY province
	`
	rows, err := r.db.Query(provinceQuery, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get province counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var province string
		var count int
		if err := rows.Scan(&province, &count); err != nil {
			return nil, fmt.Errorf("failed to scan province count: %w", err)
		}
		data.ProvincesCounts[province] = count
	}

	// Get city counts
	data.CitiesCounts = make(map[string]int)
	cityQuery := `
		SELECT city, COUNT(*) as count
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true AND city IS NOT NULL
		GROUP BY city
	`
	rows, err = r.db.Query(cityQuery, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get city counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var city string
		var count int
		if err := rows.Scan(&city, &count); err != nil {
			return nil, fmt.Errorf("failed to scan city count: %w", err)
		}
		data.CitiesCounts[city] = count
	}

	return &data, nil
}

// GetJobStatisticsForMonth aggregates job posting data for a specific month
func (r *lmiaStatisticsRepository) GetJobStatisticsForMonth(year int, month int) (*models.JobStatisticsData, error) {
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	// Get basic aggregated data
	query := `
		SELECT 
			COUNT(*) as total_jobs,
			COUNT(DISTINCT employer) as unique_employers,
			AVG(salary_min) as avg_salary_min,
			AVG(salary_max) as avg_salary_max
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true
	`

	var data models.JobStatisticsData
	row := r.db.QueryRow(query, startOfMonth, endOfMonth)
	err := row.Scan(&data.TotalJobs, &data.UniqueEmployers, &data.AvgSalaryMin, &data.AvgSalaryMax)
	if err != nil {
		return nil, fmt.Errorf("failed to get job statistics for month: %w", err)
	}

	// Get province counts
	data.ProvincesCounts = make(map[string]int)
	provinceQuery := `
		SELECT province, COUNT(*) as count
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true AND province IS NOT NULL
		GROUP BY province
	`
	rows, err := r.db.Query(provinceQuery, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get province counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var province string
		var count int
		if err := rows.Scan(&province, &count); err != nil {
			return nil, fmt.Errorf("failed to scan province count: %w", err)
		}
		data.ProvincesCounts[province] = count
	}

	// Get city counts
	data.CitiesCounts = make(map[string]int)
	cityQuery := `
		SELECT city, COUNT(*) as count
		FROM job_postings 
		WHERE posting_date >= $1 AND posting_date <= $2
		AND is_tfw = true AND has_lmia = true AND city IS NOT NULL
		GROUP BY city
	`
	rows, err = r.db.Query(cityQuery, startOfMonth, endOfMonth)
	if err != nil {
		return nil, fmt.Errorf("failed to get city counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var city string
		var count int
		if err := rows.Scan(&city, &count); err != nil {
			return nil, fmt.Errorf("failed to scan city count: %w", err)
		}
		data.CitiesCounts[city] = count
	}

	return &data, nil
}

// GetAllDatesWithJobs returns all unique dates that have job postings
func (r *lmiaStatisticsRepository) GetAllDatesWithJobs() ([]time.Time, error) {
	var dates []time.Time
	query := `
		SELECT DISTINCT DATE(posting_date) as date
		FROM job_postings 
		WHERE is_tfw = true AND has_lmia = true AND posting_date IS NOT NULL
		ORDER BY date ASC
	`

	err := r.db.Select(&dates, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all dates with jobs: %w", err)
	}

	return dates, nil
}

// DeleteStatistics deletes a statistics record
func (r *lmiaStatisticsRepository) DeleteStatistics(id string) error {
	query := `DELETE FROM lmia_job_statistics WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// GetStatisticsByID retrieves a statistics record by ID
func (r *lmiaStatisticsRepository) GetStatisticsByID(id string) (*models.LMIAStatistics, error) {
	var stats models.LMIAStatistics
	query := `SELECT * FROM lmia_job_statistics WHERE id = $1`

	err := r.db.Get(&stats, query, id)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}