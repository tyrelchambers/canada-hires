package controllers

import (
	"canada-hires/helpers"
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
)

type BoycottController interface {
	// Protected routes (auth required)
	ToggleBoycott(w http.ResponseWriter, r *http.Request)
	GetUserBoycotts(w http.ResponseWriter, r *http.Request)

	// Public routes
	GetTopBoycotted(w http.ResponseWriter, r *http.Request)
	GetBoycottStats(w http.ResponseWriter, r *http.Request)
}

type boycottController struct {
	service services.BoycottService
}

func NewBoycottController(service services.BoycottService) BoycottController {
	return &boycottController{service: service}
}

type ToggleBoycottRequest struct {
	BusinessName    string `json:"business_name"`
	BusinessAddress string `json:"business_address"`
}

type BoycottResponse struct {
	ID              string `json:"id"`
	BusinessName    string `json:"business_name"`
	BusinessAddress string `json:"business_address"`
	CreatedAt       string `json:"created_at"`
	IsBoycotting    bool   `json:"is_boycotting"`
}

type BoycottStatsResponse struct {
	BusinessName    string `json:"business_name"`
	BusinessAddress string `json:"business_address"`
	BoycottCount    int    `json:"boycott_count"`
}

func (c *boycottController) ToggleBoycott(w http.ResponseWriter, r *http.Request) {
	var req ToggleBoycottRequest
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
	serviceReq := &services.ToggleBoycottRequest{
		UserID:          user.ID,
		BusinessName:    req.BusinessName,
		BusinessAddress: req.BusinessAddress,
	}

	boycott, err := c.service.ToggleBoycott(serviceReq)
	if err != nil {
		log.Error("Failed to toggle boycott", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to toggle boycott: "+err.Error(), http.StatusBadRequest)
		return
	}

	var response BoycottResponse
	if boycott != nil {
		// Boycott was created
		response = BoycottResponse{
			ID:              boycott.ID,
			BusinessName:    boycott.BusinessName,
			BusinessAddress: *boycott.BusinessAddress,
			CreatedAt:       boycott.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsBoycotting:    true,
		}
	} else {
		// Boycott was removed
		response = BoycottResponse{
			BusinessName:    req.BusinessName,
			BusinessAddress: req.BusinessAddress,
			IsBoycotting:    false,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *boycottController) GetUserBoycotts(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	limit, offset := getPaginationParams(r)

	boycotts, err := c.service.GetUserBoycotts(user.ID, limit, offset)
	if err != nil {
		log.Error("Failed to get user boycotts", "error", err, "user_id", user.ID)
		http.Error(w, "Failed to get user boycotts", http.StatusInternalServerError)
		return
	}

	response := make([]BoycottResponse, len(boycotts))
	for i, boycott := range boycotts {
		businessAddress := ""
		if boycott.BusinessAddress != nil {
			businessAddress = *boycott.BusinessAddress
		}
		response[i] = BoycottResponse{
			ID:              boycott.ID,
			BusinessName:    boycott.BusinessName,
			BusinessAddress: businessAddress,
			CreatedAt:       boycott.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsBoycotting:    true,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":   response,
		"limit":  limit,
		"offset": offset,
		"count":  len(response),
	})
}

func (c *boycottController) GetTopBoycotted(w http.ResponseWriter, r *http.Request) {
	limit := 3 // Default to top 3
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 10 {
			limit = l
		}
	}

	stats, err := c.service.GetTopBoycottedBusinesses(limit)
	if err != nil {
		log.Error("Failed to get top boycotted businesses", "error", err)
		http.Error(w, "Failed to get top boycotted businesses", http.StatusInternalServerError)
		return
	}

	response := make([]BoycottStatsResponse, len(stats))
	for i, stat := range stats {
		response[i] = BoycottStatsResponse{
			BusinessName:    stat.BusinessName,
			BusinessAddress: stat.BusinessAddress,
			BoycottCount:    stat.BoycottCount,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *boycottController) GetBoycottStats(w http.ResponseWriter, r *http.Request) {
	businessName := r.URL.Query().Get("business_name")
	businessAddress := r.URL.Query().Get("business_address")

	if businessName == "" {
		http.Error(w, "business_name parameter is required", http.StatusBadRequest)
		return
	}

	count, err := c.service.GetBoycottCount(businessName, businessAddress)
	if err != nil {
		log.Error("Failed to get boycott count", "error", err, "business_name", businessName)
		http.Error(w, "Failed to get boycott count", http.StatusInternalServerError)
		return
	}

	// If user is authenticated, also check if they're boycotting
	user := helpers.GetUserFromContext(r.Context())
	isBoycottedByUser := false
	if user != nil {
		isBoycottedByUser, err = c.service.IsBoycottedByUser(user.ID, businessName, businessAddress)
		if err != nil {
			log.Warn("Failed to check user boycott status", "error", err, "user_id", user.ID)
			// Don't fail the request, just log the warning
		}
	}

	response := map[string]interface{}{
		"business_name":        businessName,
		"business_address":     businessAddress,
		"boycott_count":        count,
		"is_boycotted_by_user": isBoycottedByUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
