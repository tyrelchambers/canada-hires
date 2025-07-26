package controllers

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type LMIAController struct {
	lmiaService services.LMIAService
	cronService services.CronService
	repo        repos.LMIARepository
}

func NewLMIAController(lmiaService services.LMIAService, cronService services.CronService, repo repos.LMIARepository) *LMIAController {
	return &LMIAController{
		lmiaService: lmiaService,
		cronService: cronService,
		repo:        repo,
	}
}

// GetUpdateStatus returns the status of the latest LMIA data update
func (c *LMIAController) GetUpdateStatus(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting LMIA update status")

	job, err := c.lmiaService.GetLatestUpdateStatus()
	if err != nil {
		log.Error("Failed to get update status", "error", err)
		http.Error(w, "Failed to get update status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

// SearchEmployers searches for employers by name
func (c *LMIAController) SearchEmployers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 0 // default to return all records (no limit)
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	yearStr := r.URL.Query().Get("year")
	var year int
	if yearStr == "" {
		now := time.Now()
		year = time.Time.Year(now)
	} else {
		if parsedYear, err := strconv.Atoi(yearStr); err == nil && parsedYear >= 2000 && parsedYear <= time.Now().Year() {
			year = parsedYear
		}
	}

	var employers []*models.LMIAEmployer
	if query != "*" {
		var err error
		employers, err = c.repo.SearchEmployersByName(query, limit)
		if err != nil {
			log.Error("Failed to search employers", "error", err)
			http.Error(w, "Failed to search employers", http.StatusInternalServerError)
			return
		}
	} else {
		var err error
		employers, err = c.repo.GetEmployersByYear(year, limit)
		if err != nil {
			log.Error("Failed to get all employers", "error", err)
			http.Error(w, "Failed to get all employers", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employers": employers,
		"count":     len(employers),
		"query":     query,
		"limit":     limit,
	})
}

// GetEmployersByLocation gets employers by city and/or province
func (c *LMIAController) GetEmployersByLocation(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	province := r.URL.Query().Get("province")

	if city == "" && province == "" {
		http.Error(w, "At least one of 'city' or 'province' parameters is required", http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 0 // default to return all records (no limit)
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	log.Info("Getting employers by location", "city", city, "province", province, "limit", limit)

	employers, err := c.repo.GetEmployersByLocation(city, province, limit)
	if err != nil {
		log.Error("Failed to get employers by location", "error", err)
		http.Error(w, "Failed to get employers by location", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employers": employers,
		"count":     len(employers),
		"city":      city,
		"province":  province,
		"limit":     limit,
	})
}

// GetResources returns available LMIA resources
func (c *LMIAController) GetResources(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting LMIA resources")

	resources, err := c.repo.GetResourcesByLanguage("en")
	if err != nil {
		log.Error("Failed to get resources", "error", err)
		http.Error(w, "Failed to get resources", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"resources": resources,
		"count":     len(resources),
	})
}

// GetEmployersByResource gets employers for a specific resource
func (c *LMIAController) GetEmployersByResource(w http.ResponseWriter, r *http.Request) {
	resourceID := chi.URLParam(r, "resourceID")
	if resourceID == "" {
		http.Error(w, "Resource ID is required", http.StatusBadRequest)
		return
	}

	log.Info("Getting employers by resource", "resource_id", resourceID)

	employers, err := c.repo.GetEmployersByResourceID(resourceID)
	if err != nil {
		log.Error("Failed to get employers by resource", "error", err)
		http.Error(w, "Failed to get employers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"employers":   employers,
		"count":       len(employers),
		"resource_id": resourceID,
	})
}

// GetStats returns basic statistics about LMIA data
func (c *LMIAController) GetStats(w http.ResponseWriter, r *http.Request) {
	log.Info("Getting LMIA statistics")

	// Get all resources
	resources, err := c.repo.GetResourcesByLanguage("en")
	if err != nil {
		log.Error("Failed to get resources for stats", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	// Count processed resources
	processedCount := 0
	for _, resource := range resources {
		if resource.ProcessedAt != nil {
			processedCount++
		}
	}

	// Get latest update status
	latestJob, err := c.lmiaService.GetLatestUpdateStatus()
	if err != nil {
		log.Warn("Could not get latest job status", "error", err)
	}

	totalRecords, err := c.repo.AllEmployersCount()
	if err != nil {
		log.Error("Failed to get total employers count", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	distinctEmployers, err := c.repo.GetDistinctEmployersCount()
	if err != nil {
		log.Error("Failed to get distinct employers count", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	minYear, maxYear, err := c.repo.GetYearRange()
	if err != nil {
		log.Error("Failed to get year range", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	stats := map[string]interface{}{
		"total_resources":     len(resources),
		"processed_resources": processedCount,
		"last_update":         nil,
		"last_update_status":  "unknown",
		"total_records":       totalRecords,
		"distinct_employers":  distinctEmployers,
		"year_range": map[string]interface{}{
			"min_year": minYear,
			"max_year": maxYear,
		},
	}

	if latestJob != nil {
		stats["last_update"] = latestJob.StartedAt
		stats["last_update_status"] = latestJob.Status
		stats["total_records_processed"] = latestJob.RecordsProcessed
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
