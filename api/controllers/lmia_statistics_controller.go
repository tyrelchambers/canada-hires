package controllers

import (
	"canada-hires/models"
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

type LMIAStatisticsController interface {
	// Public endpoints
	GetDailyTrends(w http.ResponseWriter, r *http.Request)
	GetMonthlyTrends(w http.ResponseWriter, r *http.Request)
	GetTrendsSummary(w http.ResponseWriter, r *http.Request)

	// Admin endpoints (for manual operations)
	BackfillHistoricalStatistics(w http.ResponseWriter, r *http.Request)
	GenerateStatisticsForDateRange(w http.ResponseWriter, r *http.Request)
	RunDailyAggregation(w http.ResponseWriter, r *http.Request)
}

type lmiaStatisticsController struct {
	service services.LMIAStatisticsService
}

func NewLMIAStatisticsController(service services.LMIAStatisticsService) LMIAStatisticsController {
	return &lmiaStatisticsController{service: service}
}

// GetDailyTrends returns daily job posting trends
func (c *lmiaStatisticsController) GetDailyTrends(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")

	var stats []*models.LMIAStatistics
	var err error

	if startDateStr != "" && endDateStr != "" {
		// Parse dates
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		stats, err = c.service.GetStatisticsByDateRange(startDate, endDate, models.PeriodTypeDaily)
	} else {
		// Get latest statistics
		limit := 30 // default last 30 days
		if limitStr != "" {
			if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 {
				limit = l
				if limit > 365 {
					limit = 365 // max 1 year
				}
			}
		}

		stats, err = c.service.GetLatestStatistics(models.PeriodTypeDaily, limit)
	}

	if err != nil {
		log.Error("Failed to get daily trends", "error", err)
		http.Error(w, "Failed to get daily trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  stats,
		"count": len(stats),
	})
}

// GetMonthlyTrends returns monthly job posting trends
func (c *lmiaStatisticsController) GetMonthlyTrends(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")
	limitStr := r.URL.Query().Get("limit")

	var stats []*models.LMIAStatistics
	var err error

	if startDateStr != "" && endDateStr != "" {
		// Parse dates (expect first day of month format)
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(w, "Invalid start_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(w, "Invalid end_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		stats, err = c.service.GetStatisticsByDateRange(startDate, endDate, models.PeriodTypeMonthly)
	} else {
		// Get latest statistics
		limit := 12 // default last 12 months
		if limitStr != "" {
			if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 {
				limit = l
				if limit > 60 {
					limit = 60 // max 5 years
				}
			}
		}

		stats, err = c.service.GetLatestStatistics(models.PeriodTypeMonthly, limit)
	}

	if err != nil {
		log.Error("Failed to get monthly trends", "error", err)
		http.Error(w, "Failed to get monthly trends", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  stats,
		"count": len(stats),
	})
}

// GetTrendsSummary returns a summary of current trends
func (c *lmiaStatisticsController) GetTrendsSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := c.service.GetTrendsSummary()
	if err != nil {
		log.Error("Failed to get trends summary", "error", err)
		http.Error(w, "Failed to get trends summary", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// BackfillHistoricalStatistics backfills all historical statistics (admin only)
func (c *lmiaStatisticsController) BackfillHistoricalStatistics(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting backfill of historical statistics")

	err := c.service.BackfillAllHistoricalStatistics()
	if err != nil {
		log.Error("Failed to backfill historical statistics", "error", err)
		http.Error(w, "Failed to backfill historical statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Historical statistics backfill completed successfully",
	})
}

// GenerateStatisticsForDateRange generates statistics for a specific date range (admin only)
func (c *lmiaStatisticsController) GenerateStatisticsForDateRange(w http.ResponseWriter, r *http.Request) {
	var request struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		http.Error(w, "Invalid start_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		http.Error(w, "Invalid end_date format (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	log.Info("Generating statistics for date range", "start_date", request.StartDate, "end_date", request.EndDate)

	err = c.service.GenerateStatisticsForDateRange(startDate, endDate)
	if err != nil {
		log.Error("Failed to generate statistics for date range", "error", err)
		http.Error(w, "Failed to generate statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Statistics generated successfully for date range",
	})
}

// RunDailyAggregation manually runs the daily aggregation job (admin only)
func (c *lmiaStatisticsController) RunDailyAggregation(w http.ResponseWriter, r *http.Request) {
	log.Info("Manually running daily aggregation job")

	err := c.service.RunDailyAggregation()
	if err != nil {
		log.Error("Failed to run daily aggregation", "error", err)
		http.Error(w, "Failed to run daily aggregation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Daily aggregation completed successfully",
	})
}
