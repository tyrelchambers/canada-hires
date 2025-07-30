package controllers

import (
	"canada-hires/helpers"
	"canada-hires/models"
	"canada-hires/repos"
	"encoding/json"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type SubredditController struct {
	subredditRepo repos.SubredditRepository
}

func NewSubredditController(subredditRepo repos.SubredditRepository) *SubredditController {
	return &SubredditController{
		subredditRepo: subredditRepo,
	}
}

// GetSubreddits returns all subreddits
func (sc *SubredditController) GetSubreddits(w http.ResponseWriter, r *http.Request) {
	subreddits, err := sc.subredditRepo.GetAll()
	if err != nil {
		log.Error("Failed to get subreddits", "error", err)
		http.Error(w, "Failed to get subreddits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"subreddits": subreddits,
	})
}

// GetActiveSubreddits returns only active subreddits
func (sc *SubredditController) GetActiveSubreddits(w http.ResponseWriter, r *http.Request) {
	subreddits, err := sc.subredditRepo.GetActive()
	if err != nil {
		log.Error("Failed to get active subreddits", "error", err)
		http.Error(w, "Failed to get active subreddits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"subreddits": subreddits,
	})
}

// CreateSubreddit creates a new subreddit
func (sc *SubredditController) CreateSubreddit(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateSubredditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode create subreddit request", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Subreddit name is required", http.StatusBadRequest)
		return
	}

	// Check if subreddit already exists
	existing, err := sc.subredditRepo.GetByName(req.Name)
	if err != nil {
		log.Error("Failed to check existing subreddit", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existing != nil {
		http.Error(w, "Subreddit already exists", http.StatusConflict)
		return
	}

	subreddit := &models.Subreddit{
		ID:       uuid.New().String(),
		Name:     req.Name,
		IsActive: true, // Default to active
	}

	if req.IsActive != nil {
		subreddit.IsActive = *req.IsActive
	}

	err = sc.subredditRepo.Create(subreddit)
	if err != nil {
		log.Error("Failed to create subreddit", "error", err)
		http.Error(w, "Failed to create subreddit", http.StatusInternalServerError)
		return
	}

	log.Info("Subreddit created",
		"subreddit_id", subreddit.ID,
		"name", subreddit.Name,
		"created_by", user.Email,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subreddit)
}

// UpdateSubreddit updates an existing subreddit
func (sc *SubredditController) UpdateSubreddit(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subredditID := chi.URLParam(r, "subreddit_id")
	if subredditID == "" {
		http.Error(w, "Subreddit ID is required", http.StatusBadRequest)
		return
	}

	var req models.UpdateSubredditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("Failed to decode update subreddit request", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	// Check if subreddit exists
	existing, err := sc.subredditRepo.GetByID(subredditID)
	if err != nil {
		log.Error("Failed to get subreddit", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Subreddit not found", http.StatusNotFound)
		return
	}

	err = sc.subredditRepo.Update(subredditID, &req)
	if err != nil {
		log.Error("Failed to update subreddit", "error", err)
		http.Error(w, "Failed to update subreddit", http.StatusInternalServerError)
		return
	}

	// Get updated subreddit
	updated, err := sc.subredditRepo.GetByID(subredditID)
	if err != nil {
		log.Error("Failed to get updated subreddit", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info("Subreddit updated",
		"subreddit_id", subredditID,
		"name", existing.Name,
		"updated_by", user.Email,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// DeleteSubreddit deletes a subreddit
func (sc *SubredditController) DeleteSubreddit(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r.Context())
	if user == nil || !user.IsAdmin() {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subredditID := chi.URLParam(r, "subreddit_id")
	if subredditID == "" {
		http.Error(w, "Subreddit ID is required", http.StatusBadRequest)
		return
	}

	// Check if subreddit exists
	existing, err := sc.subredditRepo.GetByID(subredditID)
	if err != nil {
		log.Error("Failed to get subreddit", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.Error(w, "Subreddit not found", http.StatusNotFound)
		return
	}

	err = sc.subredditRepo.Delete(subredditID)
	if err != nil {
		log.Error("Failed to delete subreddit", "error", err)
		http.Error(w, "Failed to delete subreddit", http.StatusInternalServerError)
		return
	}

	log.Info("Subreddit deleted",
		"subreddit_id", subredditID,
		"name", existing.Name,
		"deleted_by", user.Email,
	)

	w.WriteHeader(http.StatusNoContent)
}
