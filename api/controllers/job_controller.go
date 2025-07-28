package controllers

import (
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type JobController struct {
	jobBankRepo repos.JobBankRepository
	jobService  services.JobService
}

func NewJobController(jobBankRepo repos.JobBankRepository, jobService services.JobService) *JobController {
	return &JobController{
		jobBankRepo: jobBankRepo,
		jobService:  jobService,
	}
}

// CreateScrapingRun starts a new job scraping session
func (jc *JobController) CreateScrapingRun(w http.ResponseWriter, r *http.Request) {
	scrapingRun := &models.JobScrapingRun{
		ID:              uuid.New().String(),
		Status:          "running",
		StartedAt:       models.TimeNow(),
		TotalPages:      0,
		JobsScraped:     0,
		JobsStored:      0,
		LastPageScraped: 0,
	}

	if err := jc.jobBankRepo.CreateScrapingRun(scrapingRun); err != nil {
		log.Error("Failed to create scraping run", "error", err)
		http.Error(w, "Failed to create scraping run", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scrapingRun)
}

// SubmitScraperJobs accepts job data from your scraper
func (jc *JobController) SubmitScraperJobs(w http.ResponseWriter, r *http.Request) {
	scrapingRunID := chi.URLParam(r, "scraping_run_id")
	if scrapingRunID == "" {
		http.Error(w, "Scraping run ID is required", http.StatusBadRequest)
		return
	}

	var scraperData []models.ScraperJobData
	if err := json.NewDecoder(r.Body).Decode(&scraperData); err != nil {
		log.Error("Failed to decode scraper data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if len(scraperData) == 0 {
		http.Error(w, "No job data provided", http.StatusBadRequest)
		return
	}

	// Store the job postings
	_, err := jc.jobBankRepo.CreateJobPostingsFromScraperData(scraperData, scrapingRunID)
	if err != nil {
		log.Error("Failed to store scraper job data", "error", err, "scraping_run_id", scrapingRunID)
		http.Error(w, "Failed to store job data", http.StatusInternalServerError)
		return
	}

	// Update scraping run progress
	jobsCount := len(scraperData)
	if err := jc.jobBankRepo.UpdateScrapingRunProgress(scrapingRunID, 0, jobsCount, jobsCount, 0); err != nil {
		log.Error("Failed to update scraping run progress", "error", err)
	}

	response := map[string]interface{}{
		"message":        "Jobs successfully stored",
		"jobs_processed": jobsCount,
		"scraping_run_id": scrapingRunID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompleteScrapingRun marks a scraping session as completed
func (jc *JobController) CompleteScrapingRun(w http.ResponseWriter, r *http.Request) {
	scrapingRunID := chi.URLParam(r, "scraping_run_id")
	if scrapingRunID == "" {
		http.Error(w, "Scraping run ID is required", http.StatusBadRequest)
		return
	}

	var body struct {
		TotalPages  int `json:"total_pages"`
		JobsScraped int `json:"jobs_scraped"`
		JobsStored  int `json:"jobs_stored"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to decode completion data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := jc.jobBankRepo.UpdateScrapingRunCompleted(scrapingRunID, body.TotalPages, body.JobsScraped, body.JobsStored); err != nil {
		log.Error("Failed to complete scraping run", "error", err, "scraping_run_id", scrapingRunID)
		http.Error(w, "Failed to complete scraping run", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":        "Scraping run completed successfully",
		"scraping_run_id": scrapingRunID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetJobPostings retrieves job postings with filtering and pagination
func (jc *JobController) GetJobPostings(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	search := r.URL.Query().Get("search")
	employer := r.URL.Query().Get("employer")
	city := r.URL.Query().Get("city")
	province := r.URL.Query().Get("province")
	title := r.URL.Query().Get("title")
	salaryMinStr := r.URL.Query().Get("salary_min")
	sortBy := r.URL.Query().Get("sort_by")
	sortOrder := r.URL.Query().Get("sort_order")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	daysStr := r.URL.Query().Get("days")

	// Parse pagination
	limit := 25 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 10000 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Parse days filter (default to all jobs if not specified)
	days := 0 // 0 means no date filter
	if daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays >= 0 {
			days = parsedDays
		}
	}

	// Parse salary filter
	var salaryMin *float64
	if salaryMinStr != "" {
		if parsedSalary, err := strconv.ParseFloat(salaryMinStr, 64); err == nil {
			salaryMin = &parsedSalary
		}
	}

	// Set default sort
	if sortBy == "" {
		sortBy = "posting_date"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Create filter parameters
	filters := map[string]interface{}{
		"search":     search,
		"employer":   employer,
		"city":       city,
		"province":   province,
		"title":      title,
		"salary_min": salaryMin,
		"days":       days,
		"sort_by":    sortBy,
		"sort_order": sortOrder,
		"limit":      limit,
		"offset":     offset,
	}

	jobs, totalCount, err := jc.jobBankRepo.SearchJobPostingsAdvanced(filters)
	if err != nil {
		log.Error("Failed to retrieve job postings", "error", err)
		http.Error(w, "Failed to retrieve job postings", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"jobs":       jobs,
		"total":      totalCount,
		"limit":      limit,
		"offset":     offset,
		"has_more":   totalCount > offset+limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetJobStats returns statistics about job postings
func (jc *JobController) GetJobStats(w http.ResponseWriter, r *http.Request) {
	totalJobs, err := jc.jobBankRepo.GetJobPostingsCount()
	if err != nil {
		log.Error("Failed to get job count", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	totalEmployers, err := jc.jobBankRepo.GetDistinctEmployersCount()
	if err != nil {
		log.Error("Failed to get employer count", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	topEmployers, err := jc.jobBankRepo.GetEmployerJobCounts(10)
	if err != nil {
		log.Error("Failed to get top employers", "error", err)
		http.Error(w, "Failed to get statistics", http.StatusInternalServerError)
		return
	}

	stats := map[string]interface{}{
		"total_jobs":      totalJobs,
		"total_employers": totalEmployers,
		"top_employers":   topEmployers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetScrapingRuns retrieves recent scraping runs
func (jc *JobController) GetScrapingRuns(w http.ResponseWriter, r *http.Request) {
	// For now, just get the latest run
	run, err := jc.jobBankRepo.GetLatestScrapingRun()
	if err != nil {
		log.Error("Failed to get scraping runs", "error", err)
		http.Error(w, "Failed to get scraping runs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]*models.JobScrapingRun{run})
}