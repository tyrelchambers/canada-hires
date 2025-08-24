package controllers

import (
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type NonCompliantController struct {
	service services.NonCompliantService
	logger  *log.Logger
}

func NewNonCompliantController(service services.NonCompliantService, logger *log.Logger) *NonCompliantController {
	return &NonCompliantController{
		service: service,
		logger:  logger,
	}
}

// GetNonCompliantEmployers handles GET /api/non-compliant/employers
func (c *NonCompliantController) GetNonCompliantEmployers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 25 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = parsedLimit
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get employers
	employers, err := c.service.GetNonCompliantEmployers(limit, offset)
	if err != nil {
		c.logger.Error("Failed to get non-compliant employers", "error", err)
		http.Error(w, "Failed to retrieve non-compliant employers", http.StatusInternalServerError)
		return
	}

	// Get total count for pagination
	totalCount, err := c.service.GetNonCompliantEmployersCount()
	if err != nil {
		c.logger.Error("Failed to get non-compliant employers count", "error", err)
		http.Error(w, "Failed to retrieve employers count", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"employers":    employers,
		"total_count":  totalCount,
		"limit":        limit,
		"offset":       offset,
		"has_more":     offset+limit < totalCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNonCompliantReasons handles GET /api/non-compliant/reasons
func (c *NonCompliantController) GetNonCompliantReasons(w http.ResponseWriter, r *http.Request) {
	reasons, err := c.service.GetNonCompliantReasons()
	if err != nil {
		c.logger.Error("Failed to get non-compliant reasons", "error", err)
		http.Error(w, "Failed to retrieve non-compliant reasons", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"reasons": reasons,
		"count":   len(reasons),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TriggerNonCompliantScraper handles POST /api/admin/non-compliant/scrape
func (c *NonCompliantController) TriggerNonCompliantScraper(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Non-compliant scraper triggered via API")

	// Start the scraping process
	job, err := c.service.ScrapeAndStoreNonCompliantEmployers()
	if err != nil {
		c.logger.Error("Failed to trigger non-compliant scraper", "error", err)
		http.Error(w, "Failed to start non-compliant scraper", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Non-compliant employers scraper started successfully",
		"job":     job,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNonCompliantScrapingStatus handles GET /api/admin/non-compliant/status
func (c *NonCompliantController) GetNonCompliantScrapingStatus(w http.ResponseWriter, r *http.Request) {
	job, err := c.service.GetLatestScrapeInfo()
	if err != nil {
		c.logger.Error("Failed to get non-compliant scraping status", "error", err)
		http.Error(w, "Failed to retrieve scraping status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"latest_job": job,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNonCompliantLocations handles GET /api/non-compliant/locations
func (c *NonCompliantController) GetNonCompliantLocations(w http.ResponseWriter, r *http.Request) {
	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 1000 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 5000 {
			limit = parsedLimit
		}
	}

	// Get locations
	response, err := c.service.GetNonCompliantLocationsByPostalCode(limit)
	if err != nil {
		c.logger.Error("Failed to get non-compliant locations", "error", err)
		http.Error(w, "Failed to retrieve non-compliant locations", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetNonCompliantEmployersByPostalCode handles GET /api/non-compliant/employers/postal-code/{postal_code}
func (c *NonCompliantController) GetNonCompliantEmployersByPostalCode(w http.ResponseWriter, r *http.Request) {
	// Get postal code from URL
	postalCode := chi.URLParam(r, "postal_code")
	if postalCode == "" {
		http.Error(w, "Postal code is required", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 100 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 500 {
			limit = parsedLimit
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Get employers by postal code
	response, err := c.service.GetNonCompliantEmployersByPostalCode(postalCode, limit, offset)
	if err != nil {
		c.logger.Error("Failed to get non-compliant employers by postal code", "error", err, "postal_code", postalCode)
		http.Error(w, "Failed to retrieve employers for postal code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TriggerNonCompliantGeocoding handles POST /api/admin/non-compliant/geocode
func (c *NonCompliantController) TriggerNonCompliantGeocoding(w http.ResponseWriter, r *http.Request) {
	c.logger.Info("Non-compliant geocoding triggered via API")

	// Start the geocoding process
	err := c.service.ExtractAndGeocodeEmployers()
	if err != nil {
		c.logger.Error("Failed to trigger non-compliant geocoding", "error", err)
		http.Error(w, "Failed to start geocoding process", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Non-compliant employers geocoding completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetupNonCompliantRoutes sets up all routes for non-compliant employers
func (c *NonCompliantController) SetupNonCompliantRoutes(r chi.Router) {
	// Public API routes
	r.Route("/api/non-compliant", func(r chi.Router) {
		r.Get("/employers", c.GetNonCompliantEmployers)
		r.Get("/reasons", c.GetNonCompliantReasons)
		r.Get("/locations", c.GetNonCompliantLocations)
		r.Get("/employers/postal-code/{postal_code}", c.GetNonCompliantEmployersByPostalCode)
	})

	// Admin routes (these should have authentication middleware in production)
	r.Route("/api/admin/non-compliant", func(r chi.Router) {
		r.Post("/scrape", c.TriggerNonCompliantScraper)
		r.Get("/status", c.GetNonCompliantScrapingStatus)
		r.Post("/geocode", c.TriggerNonCompliantGeocoding)
	})
}