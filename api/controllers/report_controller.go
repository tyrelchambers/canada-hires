package controllers

import (
	"canada-hires/dto"
	"canada-hires/helpers"
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

type ReportController interface {
	// Public routes
	CreateReport(w http.ResponseWriter, r *http.Request)
	GetReportByID(w http.ResponseWriter, r *http.Request)
	GetReports(w http.ResponseWriter, r *http.Request)
	GetReportsGrouped(w http.ResponseWriter, r *http.Request)

	// Protected routes (auth required)
	GetUserReports(w http.ResponseWriter, r *http.Request)
	UpdateReport(w http.ResponseWriter, r *http.Request)
	DeleteReport(w http.ResponseWriter, r *http.Request)
}

type reportController struct {
	service services.ReportService
}

func NewReportController(service services.ReportService) ReportController {
	return &reportController{service: service}
}

func (c *reportController) CreateReport(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get user from context (required)
	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create service request
	clientIP := helpers.GetClientIP(r)
	serviceReq := &services.CreateReportRequest{
		UserID:          user.ID,
		BusinessName:    req.BusinessName,
		BusinessAddress: req.BusinessAddress,
		ReportSource:    req.ReportSource,
		ConfidenceLevel: req.ConfidenceLevel,
		AdditionalNotes: req.AdditionalNotes,
		IPAddress:       &clientIP,
	}

	report, err := c.service.CreateReport(serviceReq)
	if err != nil {
		log.Error("Failed to create report", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to create report: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.ToReportResponse(report)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) GetReportByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	report, err := c.service.GetReportByID(id)
	if err != nil {
		log.Error("Failed to get report", "error", err, "report_id", id)
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	response := dto.ToReportResponse(report)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) GetReports(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPaginationParams(r)

	// Get search/filter parameters from query string
	query := r.URL.Query().Get("query")
	city := r.URL.Query().Get("city")
	province := r.URL.Query().Get("province")
	year := r.URL.Query().Get("year")
	address := r.URL.Query().Get("address")

	// Get URL parameters for specific filters
	businessName := chi.URLParam(r, "businessName")

	// Handle business name from URL parameter
	if businessName != "" {
		reports, err := c.service.GetBusinessReports(businessName, limit, offset)
		if err != nil {
			log.Error("Failed to get business reports", "error", err, "business_name", businessName)
			http.Error(w, "Failed to get business reports", http.StatusInternalServerError)
			return
		}
		response := dto.ToReportListResponse(reports, limit, offset)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}


	// Handle address filter
	if address != "" {
		reports, err := c.service.GetAddressReports(address)
		if err != nil {
			log.Error("Failed to get address reports", "error", err, "address", address)
			http.Error(w, "Failed to get address reports", http.StatusInternalServerError)
			return
		}
		response := dto.ToReportListResponse(reports, len(reports), 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// If any filters are provided, use the filtered search
	if query != "" || city != "" || province != "" || year != "" {
		filters := services.ReportFilters{
			Query:    query,
			City:     city,
			Province: province,
			Year:     year,
		}

		reports, err := c.service.GetReportsWithFilters(filters, limit, offset)
		if err != nil {
			log.Error("Failed to get filtered reports", "error", err, "filters", filters)
			http.Error(w, "Failed to get reports", http.StatusInternalServerError)
			return
		}

		response := dto.ToReportListResponse(reports, limit, offset)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Otherwise, use the regular GetAllReports
	reports, err := c.service.GetAllReports(limit, offset)
	if err != nil {
		log.Error("Failed to get reports", "error", err)
		http.Error(w, "Failed to get reports", http.StatusInternalServerError)
		return
	}

	response := dto.ToReportListResponse(reports, limit, offset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) GetUserReports(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, offset := getPaginationParams(r)

	reports, err := c.service.GetUserReports(user.ID, limit, offset)
	if err != nil {
		log.Error("Failed to get user reports", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to get user reports", http.StatusInternalServerError)
		return
	}

	response := dto.ToReportListResponse(reports, limit, offset)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) UpdateReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.UpdateReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get existing report to check ownership
	existingReport, err := c.service.GetReportByID(id)
	if err != nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	// Check ownership (only owner or admin can update)
	if existingReport.UserID != user.ID && !user.IsAdmin() {
		http.Error(w, "Forbidden: You can only update your own reports", http.StatusForbidden)
		return
	}

	// Update the report
	existingReport.BusinessName = req.BusinessName
	existingReport.BusinessAddress = req.BusinessAddress
	existingReport.ReportSource = req.ReportSource
	existingReport.ConfidenceLevel = req.ConfidenceLevel
	existingReport.AdditionalNotes = req.AdditionalNotes

	err = c.service.UpdateReport(existingReport)
	if err != nil {
		log.Error("Failed to update report", "error", err, "report_id", id, "user_id", user.ID)
		http.Error(w, "Failed to update report: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := dto.ToReportResponse(existingReport)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *reportController) DeleteReport(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "Report ID is required", http.StatusBadRequest)
		return
	}

	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := c.service.DeleteReport(id, user.ID, user.IsAdmin())
	if err != nil {
		log.Error("Failed to delete report", "error", err, "report_id", id, "user_id", user.ID)
		if strings.Contains(err.Error(), "unauthorized") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Failed to delete report", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (c *reportController) GetReportsGrouped(w http.ResponseWriter, r *http.Request) {
	limit, offset := getPaginationParams(r)

	// Get search/filter parameters
	query := r.URL.Query().Get("query")
	city := r.URL.Query().Get("city")
	province := r.URL.Query().Get("province")
	year := r.URL.Query().Get("year")

	// Build filters if any are provided
	var filters *services.ReportFilters
	if query != "" || city != "" || province != "" || year != "" {
		filters = &services.ReportFilters{
			Query:    query,
			City:     city,
			Province: province,
			Year:     year,
		}
	}

	// Use the single consolidated method
	grouped, err := c.service.GetReportsGroupedByAddress(filters, limit, offset)
	if err != nil {
		log.Error("Failed to get reports grouped by address", "error", err, "filters", filters)
		http.Error(w, "Failed to get grouped reports", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   grouped,
		"limit":  limit,
		"offset": offset,
		"count":  len(grouped),
	})
}

// Helper functions
func getPaginationParams(r *http.Request) (limit, offset int) {
	limit = 50 // default
	offset = 0 // default

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
			if limit > 100 {
				limit = 100 // max limit
			}
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	return limit, offset
}
