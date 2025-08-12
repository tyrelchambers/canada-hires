package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type LMIAStatisticsService interface {
	// Generate statistics
	GenerateDailyStatistics(date time.Time) (*models.LMIAStatistics, error)
	GenerateMonthlyStatistics(year int, month int) (*models.LMIAStatistics, error)
	
	// Batch operations
	GenerateStatisticsForDateRange(startDate, endDate time.Time) error
	BackfillAllHistoricalStatistics() error
	
	// Get statistics
	GetStatisticsByDateRange(startDate, endDate time.Time, periodType models.PeriodType) ([]*models.LMIAStatistics, error)
	GetLatestStatistics(periodType models.PeriodType, limit int) ([]*models.LMIAStatistics, error)
	GetTrendsSummary() (*TrendsSummary, error)
	
	// Daily aggregation job
	RunDailyAggregation() error
}

type TrendsSummary struct {
	TotalJobsToday      int                   `json:"total_jobs_today"`
	TotalJobsThisMonth  int                   `json:"total_jobs_this_month"`
	TotalJobsLastMonth  int                   `json:"total_jobs_last_month"`
	PercentageChange    float64               `json:"percentage_change"`
	TopProvincesToday   []models.RegionData   `json:"top_provinces_today"`
	TopCitiesToday      []models.RegionData   `json:"top_cities_today"`
	RecentTrends        []*models.LMIAStatistics `json:"recent_trends"`
}

type lmiaStatisticsService struct {
	repo repos.LMIAStatisticsRepository
}

func NewLMIAStatisticsService(repo repos.LMIAStatisticsRepository) LMIAStatisticsService {
	return &lmiaStatisticsService{repo: repo}
}

// GenerateDailyStatistics generates and stores daily statistics for a given date
func (s *lmiaStatisticsService) GenerateDailyStatistics(date time.Time) (*models.LMIAStatistics, error) {
	// Get aggregated job data for the date
	jobData, err := s.repo.GetJobStatisticsForDate(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get job statistics for date %v: %w", date, err)
	}

	// Create statistics model
	stats := models.NewLMIAStatistics(date, models.PeriodTypeDaily, *jobData)

	// Upsert the statistics (insert or update if exists)
	err = s.repo.UpsertStatistics(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert daily statistics: %w", err)
	}

	log.Info("Generated daily statistics", "date", date.Format("2006-01-02"), "total_jobs", stats.TotalJobs, "unique_employers", stats.UniqueEmployers)
	return stats, nil
}

// GenerateMonthlyStatistics generates and stores monthly statistics for a given month
func (s *lmiaStatisticsService) GenerateMonthlyStatistics(year int, month int) (*models.LMIAStatistics, error) {
	// Get aggregated job data for the month
	jobData, err := s.repo.GetJobStatisticsForMonth(year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get job statistics for month %d/%d: %w", month, year, err)
	}

	// Use the first day of the month as the date
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	
	// Create statistics model
	stats := models.NewLMIAStatistics(date, models.PeriodTypeMonthly, *jobData)

	// Upsert the statistics
	err = s.repo.UpsertStatistics(stats)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert monthly statistics: %w", err)
	}

	log.Info("Generated monthly statistics", "year", year, "month", month, "total_jobs", stats.TotalJobs, "unique_employers", stats.UniqueEmployers)
	return stats, nil
}

// GenerateStatisticsForDateRange generates statistics for all dates in a range
func (s *lmiaStatisticsService) GenerateStatisticsForDateRange(startDate, endDate time.Time) error {
	log.Info("Starting statistics generation for date range", "start", startDate.Format("2006-01-02"), "end", endDate.Format("2006-01-02"))

	// Generate daily statistics
	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		_, err := s.GenerateDailyStatistics(currentDate)
		if err != nil {
			log.Error("Failed to generate daily statistics", "date", currentDate.Format("2006-01-02"), "error", err)
			// Continue with next date instead of failing completely
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	// Generate monthly statistics for all months in the range
	currentMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())
	
	for currentMonth.Before(endMonth) || currentMonth.Equal(endMonth) {
		_, err := s.GenerateMonthlyStatistics(currentMonth.Year(), int(currentMonth.Month()))
		if err != nil {
			log.Error("Failed to generate monthly statistics", "year", currentMonth.Year(), "month", int(currentMonth.Month()), "error", err)
			// Continue with next month instead of failing completely
		}
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	log.Info("Completed statistics generation for date range")
	return nil
}

// BackfillAllHistoricalStatistics generates statistics for all historical job data
func (s *lmiaStatisticsService) BackfillAllHistoricalStatistics() error {
	log.Info("Starting backfill of all historical statistics")

	// Get all unique dates with jobs
	dates, err := s.repo.GetAllDatesWithJobs()
	if err != nil {
		return fmt.Errorf("failed to get dates with jobs: %w", err)
	}

	if len(dates) == 0 {
		log.Info("No job posting dates found, nothing to backfill")
		return nil
	}

	startDate := dates[0]
	endDate := dates[len(dates)-1]

	log.Info("Backfilling statistics", "start_date", startDate.Format("2006-01-02"), "end_date", endDate.Format("2006-01-02"), "total_dates", len(dates))

	return s.GenerateStatisticsForDateRange(startDate, endDate)
}

// GetStatisticsByDateRange retrieves statistics for a date range
func (s *lmiaStatisticsService) GetStatisticsByDateRange(startDate, endDate time.Time, periodType models.PeriodType) ([]*models.LMIAStatistics, error) {
	return s.repo.GetStatisticsByDateRange(startDate, endDate, periodType)
}

// GetLatestStatistics retrieves the most recent statistics
func (s *lmiaStatisticsService) GetLatestStatistics(periodType models.PeriodType, limit int) ([]*models.LMIAStatistics, error) {
	return s.repo.GetLatestStatistics(periodType, limit)
}

// GetTrendsSummary provides a summary of current trends
func (s *lmiaStatisticsService) GetTrendsSummary() (*TrendsSummary, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	// Get today's statistics
	todayStats, err := s.repo.GetStatisticsByDate(today, models.PeriodTypeDaily)
	var totalJobsToday int
	var topProvincesToday, topCitiesToday []models.RegionData
	
	if err == nil {
		totalJobsToday = todayStats.TotalJobs
		topProvincesToday, _ = todayStats.GetTopProvinces()
		topCitiesToday, _ = todayStats.GetTopCities()
	}

	// Get this month's statistics
	thisMonthFirst := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	thisMonthStats, err := s.repo.GetStatisticsByDate(thisMonthFirst, models.PeriodTypeMonthly)
	var totalJobsThisMonth int
	if err == nil {
		totalJobsThisMonth = thisMonthStats.TotalJobs
	}

	// Get last month's statistics
	lastMonthFirst := thisMonthFirst.AddDate(0, -1, 0)
	lastMonthStats, err := s.repo.GetStatisticsByDate(lastMonthFirst, models.PeriodTypeMonthly)
	var totalJobsLastMonth int
	if err == nil {
		totalJobsLastMonth = lastMonthStats.TotalJobs
	}

	// Calculate percentage change
	var percentageChange float64
	if totalJobsLastMonth > 0 {
		percentageChange = float64(totalJobsThisMonth-totalJobsLastMonth) / float64(totalJobsLastMonth) * 100
	}

	// Get recent daily trends (last 30 days)
	thirtyDaysAgo := today.AddDate(0, 0, -30)
	recentTrends, err := s.repo.GetStatisticsByDateRange(thirtyDaysAgo, today, models.PeriodTypeDaily)
	if err != nil {
		log.Error("Failed to get recent trends", "error", err)
		recentTrends = []*models.LMIAStatistics{}
	}

	return &TrendsSummary{
		TotalJobsToday:     totalJobsToday,
		TotalJobsThisMonth: totalJobsThisMonth,
		TotalJobsLastMonth: totalJobsLastMonth,
		PercentageChange:   percentageChange,
		TopProvincesToday:  topProvincesToday,
		TopCitiesToday:     topCitiesToday,
		RecentTrends:       recentTrends,
	}, nil
}

// RunDailyAggregation runs the daily aggregation job (typically called by cron)
func (s *lmiaStatisticsService) RunDailyAggregation() error {
	log.Info("Starting daily aggregation job")

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	// Generate statistics for yesterday
	_, err := s.GenerateDailyStatistics(yesterday)
	if err != nil {
		return fmt.Errorf("failed to generate daily statistics for yesterday: %w", err)
	}

	// Also update monthly statistics for the current month
	_, err = s.GenerateMonthlyStatistics(now.Year(), int(now.Month()))
	if err != nil {
		log.Error("Failed to update monthly statistics", "error", err)
		// Don't return error for monthly stats failure, daily is more important
	}

	// If it's the first day of the month, also update last month's stats
	if now.Day() == 1 {
		lastMonth := now.AddDate(0, -1, 0)
		_, err = s.GenerateMonthlyStatistics(lastMonth.Year(), int(lastMonth.Month()))
		if err != nil {
			log.Error("Failed to update last month's statistics", "error", err)
		}
	}

	log.Info("Daily aggregation job completed successfully")
	return nil
}