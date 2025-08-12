package controllers

import (
	"canada-hires/helpers"
	"canada-hires/models"
	"canada-hires/repos"
	"canada-hires/services"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type JobController struct {
	jobBankRepo      repos.JobBankRepository
	jobService       services.JobService
	redditService    services.RedditService
	scraperCronService *services.ScraperCronService
}

func NewJobController(jobBankRepo repos.JobBankRepository, jobService services.JobService, redditService services.RedditService, scraperCronService *services.ScraperCronService) *JobController {
	return &JobController{
		jobBankRepo:        jobBankRepo,
		jobService:         jobService,
		redditService:      redditService,
		scraperCronService: scraperCronService,
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
		"message":         "Jobs successfully stored",
		"jobs_processed":  jobsCount,
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
		"message":         "Scraping run completed successfully",
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
		"jobs":     jobs,
		"total":    totalCount,
		"limit":    limit,
		"offset":   offset,
		"has_more": totalCount > offset+limit,
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

// ADMIN ENDPOINTS FOR REDDIT APPROVAL WORKFLOW

// GetPendingJobsForReddit retrieves jobs pending Reddit approval
func (jc *JobController) GetPendingJobsForReddit(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	filters := map[string]interface{}{
		"reddit_approval_status": "pending",
		"sort_by":                "posting_date",
		"sort_order":             "desc",
		"limit":                  limit,
		"offset":                 offset,
	}

	jobs, totalCount, err := jc.jobBankRepo.SearchJobPostingsAdvanced(filters)
	if err != nil {
		log.Error("Failed to retrieve pending jobs", "error", err)
		http.Error(w, "Failed to retrieve pending jobs", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"jobs":     jobs,
		"total":    totalCount,
		"limit":    limit,
		"offset":   offset,
		"has_more": totalCount > offset+limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPostedJobsForReddit retrieves jobs that have been approved and posted to Reddit
func (jc *JobController) GetPostedJobsForReddit(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	offset := 0
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	filters := map[string]interface{}{
		"reddit_approval_status": "approved",
		"sort_by":                "reddit_approved_at",
		"sort_order":             "desc",
		"limit":                  limit,
		"offset":                 offset,
	}

	jobs, totalCount, err := jc.jobBankRepo.SearchJobPostingsAdvanced(filters)
	if err != nil {
		log.Error("Failed to retrieve posted jobs", "error", err)
		http.Error(w, "Failed to retrieve posted jobs", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"jobs":     jobs,
		"total":    totalCount,
		"limit":    limit,
		"offset":   offset,
		"has_more": totalCount > offset+limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ApproveJobForReddit approves a job for Reddit posting
func (jc *JobController) ApproveJobForReddit(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	var body struct {
		ApprovedBy string `json:"approved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to decode approval data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if body.ApprovedBy == "" {
		http.Error(w, "approved_by is required", http.StatusBadRequest)
		return
	}

	// Update job approval status
	now := time.Now()
	err := jc.jobBankRepo.UpdateJobRedditApprovalStatus(jobID, "approved", body.ApprovedBy, &now, nil)
	if err != nil {
		log.Error("Failed to approve job for Reddit", "error", err, "job_id", jobID)
		http.Error(w, "Failed to approve job", http.StatusInternalServerError)
		return
	}

	// Get the updated job
	job, err := jc.jobBankRepo.GetJobPostingByID(jobID)
	if err != nil {
		log.Error("Failed to retrieve approved job", "error", err, "job_id", jobID)
		http.Error(w, "Job approved but failed to retrieve details", http.StatusInternalServerError)
		return
	}

	// Post to Reddit asynchronously (skip in development mode)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := jc.redditService.PostJob(ctx, job); err != nil {
			log.Error("Failed to post approved job to Reddit",
				"error", err,
				"job_id", job.ID,
				"job_title", job.Title,
			)
		} else {
			log.Info("Successfully posted approved job to Reddit",
				"job_id", job.ID,
				"job_title", job.Title,
				"approved_by", body.ApprovedBy,
			)
		}
	}()

	response := map[string]interface{}{
		"message": "Job approved for Reddit posting",
		"job_id":  jobID,
		"status":  "approved",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RejectJobForReddit rejects a job for Reddit posting
func (jc *JobController) RejectJobForReddit(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "job_id")
	if jobID == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	var body struct {
		RejectedBy string `json:"rejected_by"`
		Reason     string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to decode rejection data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if body.RejectedBy == "" {
		http.Error(w, "rejected_by is required", http.StatusBadRequest)
		return
	}

	// Update job approval status
	now := time.Now()
	var reason *string
	if body.Reason != "" {
		reason = &body.Reason
	}

	err := jc.jobBankRepo.UpdateJobRedditApprovalStatus(jobID, "rejected", body.RejectedBy, &now, reason)
	if err != nil {
		log.Error("Failed to reject job for Reddit", "error", err, "job_id", jobID)
		http.Error(w, "Failed to reject job", http.StatusInternalServerError)
		return
	}

	log.Info("Job rejected for Reddit posting",
		"job_id", jobID,
		"rejected_by", body.RejectedBy,
		"reason", body.Reason,
	)

	response := map[string]interface{}{
		"message": "Job rejected for Reddit posting",
		"job_id":  jobID,
		"status":  "rejected",
		"reason":  body.Reason,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BulkApproveJobsForReddit approves multiple jobs for Reddit posting
func (jc *JobController) BulkApproveJobsForReddit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		JobIDs     []string `json:"job_ids"`
		ApprovedBy string   `json:"approved_by"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to decode bulk approval data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if len(body.JobIDs) == 0 {
		http.Error(w, "job_ids is required", http.StatusBadRequest)
		return
	}

	if body.ApprovedBy == "" {
		http.Error(w, "approved_by is required", http.StatusBadRequest)
		return
	}

	approvedCount := 0
	failedIDs := []string{}

	for _, jobID := range body.JobIDs {
		now := time.Now()
		err := jc.jobBankRepo.UpdateJobRedditApprovalStatus(jobID, "approved", body.ApprovedBy, &now, nil)
		if err != nil {
			log.Error("Failed to approve job in bulk operation", "error", err, "job_id", jobID)
			failedIDs = append(failedIDs, jobID)
			continue
		}

		// Get the job and post to Reddit asynchronously (skip in development mode)
		if jc.redditService != nil && !helpers.IsDev() {
			go func(id string) {
				job, err := jc.jobBankRepo.GetJobPostingByID(id)
				if err != nil {
					log.Error("Failed to retrieve job for Reddit posting", "error", err, "job_id", id)
					return
				}

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := jc.redditService.PostJob(ctx, job); err != nil {
					log.Error("Failed to post bulk approved job to Reddit",
						"error", err,
						"job_id", job.ID,
						"job_title", job.Title,
					)
				} else {
					log.Info("Successfully posted bulk approved job to Reddit",
						"job_id", job.ID,
						"job_title", job.Title,
						"approved_by", body.ApprovedBy,
					)
				}
			}(jobID)
		} else if helpers.IsDev() {
			log.Info("Skipping Reddit post in development mode",
				"job_id", jobID,
				"approved_by", body.ApprovedBy,
			)
		}

		approvedCount++
	}

	log.Info("Bulk approval completed",
		"total_requested", len(body.JobIDs),
		"approved", approvedCount,
		"failed", len(failedIDs),
		"approved_by", body.ApprovedBy,
	)

	response := map[string]interface{}{
		"message":        "Bulk approval completed",
		"approved_count": approvedCount,
		"failed_count":   len(failedIDs),
		"failed_ids":     failedIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// BulkRejectJobsForReddit rejects multiple jobs for Reddit posting
func (jc *JobController) BulkRejectJobsForReddit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		JobIDs     []string `json:"job_ids"`
		RejectedBy string   `json:"rejected_by"`
		Reason     string   `json:"reason,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to decode bulk rejection data", "error", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if len(body.JobIDs) == 0 {
		http.Error(w, "job_ids is required", http.StatusBadRequest)
		return
	}

	if body.RejectedBy == "" {
		http.Error(w, "rejected_by is required", http.StatusBadRequest)
		return
	}

	rejectedCount := 0
	failedIDs := []string{}

	for _, jobID := range body.JobIDs {
		now := time.Now()
		var rejectionReason *string
		if body.Reason != "" {
			rejectionReason = &body.Reason
		}

		err := jc.jobBankRepo.UpdateJobRedditApprovalStatus(jobID, "rejected", body.RejectedBy, &now, rejectionReason)
		if err != nil {
			log.Error("Failed to reject job in bulk operation", "error", err, "job_id", jobID)
			failedIDs = append(failedIDs, jobID)
			continue
		}

		// Log the rejection
		log.Info("Job rejected for Reddit posting",
			"job_id", jobID,
			"rejected_by", body.RejectedBy,
			"reason", rejectionReason,
		)

		rejectedCount++
	}

	log.Info("Bulk rejection completed",
		"total_requested", len(body.JobIDs),
		"rejected", rejectedCount,
		"failed", len(failedIDs),
		"rejected_by", body.RejectedBy,
	)

	response := map[string]interface{}{
		"message":        "Bulk rejection completed",
		"rejected_count": rejectedCount,
		"failed_count":   len(failedIDs),
		"failed_ids":     failedIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


// TriggerScraper manually triggers the scraper and statistics aggregation
func (jc *JobController) TriggerScraper(w http.ResponseWriter, r *http.Request) {
	log.Info("Manual scraper trigger requested")

	// Trigger the scraper execution
	err := jc.scraperCronService.RunNow()
	if err != nil {
		log.Error("Failed to trigger scraper", "error", err)
		http.Error(w, "Failed to trigger scraper: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Scraper job triggered successfully",
		"status":  "started",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TriggerStatisticsAggregation manually triggers the statistics aggregation
func (jc *JobController) TriggerStatisticsAggregation(w http.ResponseWriter, r *http.Request) {
	log.Info("Manual statistics aggregation trigger requested")

	// Trigger the statistics aggregation
	err := jc.scraperCronService.RunStatisticsAggregationNow()
	if err != nil {
		log.Error("Failed to trigger statistics aggregation", "error", err)
		http.Error(w, "Failed to trigger statistics aggregation: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Statistics aggregation triggered successfully",
		"status":  "started",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
