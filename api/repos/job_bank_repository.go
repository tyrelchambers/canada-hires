package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type JobBankRepository interface {
	// Job Scraping Runs
	CreateScrapingRun(run *models.JobScrapingRun) error
	UpdateScrapingRunStatus(id string, status string, errorMessage *string) error
	UpdateScrapingRunProgress(id string, totalPages, jobsScraped, jobsStored, lastPageScraped int) error
	UpdateScrapingRunCompleted(id string, totalPages, jobsScraped, jobsStored int) error
	GetLatestScrapingRun() (*models.JobScrapingRun, error)
	GetScrapingRunByID(id string) (*models.JobScrapingRun, error)

	// Job Postings
	CreateJobPosting(posting *models.JobPosting) error
	CreateJobPostingsBatch(postings []*models.JobPosting) error
	CreateJobPostingsFromScraperData(scraperData []models.ScraperJobData, scrapingRunID string) ([]*models.JobPosting, error)
	GetJobPostingByJobBankID(jobBankID string) (*models.JobPosting, error)
	GetJobPostingByID(id string) (*models.JobPosting, error)
	GetJobPostingByURL(url string) (*models.JobPosting, error)
	UpdateJobPostingRedditStatus(id string, redditPosted bool) error
	UpdateJobRedditApprovalStatus(id string, status string, approvedBy string, approvedAt *time.Time, rejectionReason *string) error
	GetJobPostingsNotPostedToReddit(limit int) ([]*models.JobPosting, error)
	SearchJobPostingsByEmployer(employer string, limit int) ([]*models.JobPosting, error)
	GetJobPostingsByLocation(city, province string, limit int) ([]*models.JobPosting, error)
	GetJobPostingsByScrapingRun(scrapingRunID string) ([]*models.JobPosting, error)
	GetRecentJobPostings(limit int) ([]*models.JobPosting, error)
	SearchJobPostingsAdvanced(filters map[string]interface{}) ([]*models.JobPosting, int, error)
	GetJobPostingsCount() (int, error)
	GetDistinctEmployersCount() (int, error)
	GetEmployerJobCounts(limit int) ([]map[string]interface{}, error)
	DeleteJobPostingsNotInScrapeRun(scrapingRunID string, currentJobBankIDs []string) (int, error)
}

type jobBankRepository struct {
	db *sqlx.DB
}

func NewJobBankRepository(db *sqlx.DB) JobBankRepository {
	return &jobBankRepository{db: db}
}

// Job Scraping Runs methods
func (r *jobBankRepository) CreateScrapingRun(run *models.JobScrapingRun) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO job_scraping_runs (id, status, started_at, created_at)
		VALUES (:id, :status, :started_at, :created_at)
	`

	run.ID = uuid.New().String()
	run.CreatedAt = time.Now()

	_, err = tx.NamedExec(query, run)
	if err != nil {
		return fmt.Errorf("failed to insert job scraping run: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *jobBankRepository) UpdateScrapingRunStatus(id string, status string, errorMessage *string) error {
	query := `UPDATE job_scraping_runs SET status = $2, error_message = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, errorMessage)
	return err
}

func (r *jobBankRepository) UpdateScrapingRunProgress(id string, totalPages, jobsScraped, jobsStored, lastPageScraped int) error {
	query := `
		UPDATE job_scraping_runs
		SET total_pages = $2, jobs_scraped = $3, jobs_stored = $4, last_page_scraped = $5
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, totalPages, jobsScraped, jobsStored, lastPageScraped)
	return err
}

func (r *jobBankRepository) UpdateScrapingRunCompleted(id string, totalPages, jobsScraped, jobsStored int) error {
	query := `
		UPDATE job_scraping_runs
		SET status = 'completed', completed_at = NOW(), total_pages = $2, jobs_scraped = $3, jobs_stored = $4
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, totalPages, jobsScraped, jobsStored)
	return err
}

func (r *jobBankRepository) GetLatestScrapingRun() (*models.JobScrapingRun, error) {
	var run models.JobScrapingRun
	query := `SELECT * FROM job_scraping_runs ORDER BY started_at DESC LIMIT 1`

	err := r.db.Get(&run, query)
	if err != nil {
		return nil, err
	}

	return &run, nil
}

func (r *jobBankRepository) GetScrapingRunByID(id string) (*models.JobScrapingRun, error) {
	var run models.JobScrapingRun
	query := `SELECT * FROM job_scraping_runs WHERE id = $1`

	err := r.db.Get(&run, query, id)
	if err != nil {
		return nil, err
	}

	return &run, nil
}

// Job Postings methods
func (r *jobBankRepository) CreateJobPosting(posting *models.JobPosting) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO job_postings (id, job_bank_id, title, employer, location, province, city,
								 salary_min, salary_max, salary_type, posting_date, url, is_tfw,
								 has_lmia, reddit_posted, description, scraping_run_id, created_at, updated_at)
		VALUES (:id, :job_bank_id, :title, :employer, :location, :province, :city,
				:salary_min, :salary_max, :salary_type, :posting_date, :url, :is_tfw,
				:has_lmia, :reddit_posted, :description, :scraping_run_id, :created_at, :updated_at)
		ON CONFLICT (job_bank_id) DO UPDATE SET
			title = EXCLUDED.title,
			employer = EXCLUDED.employer,
			location = EXCLUDED.location,
			province = EXCLUDED.province,
			city = EXCLUDED.city,
			salary_min = EXCLUDED.salary_min,
			salary_max = EXCLUDED.salary_max,
			salary_type = EXCLUDED.salary_type,
			posting_date = EXCLUDED.posting_date,
			url = EXCLUDED.url,
			has_lmia = EXCLUDED.has_lmia,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`

	posting.ID = uuid.New().String()
	posting.CreatedAt = time.Now()
	posting.UpdatedAt = time.Now()

	_, err = tx.NamedExec(query, posting)
	if err != nil {
		return fmt.Errorf("failed to insert job posting: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *jobBankRepository) CreateJobPostingsBatch(postings []*models.JobPosting) error {
	if len(postings) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO job_postings (id, job_bank_id, title, employer, location, province, city,
								 salary_min, salary_max, salary_type, posting_date, url, is_tfw,
								 has_lmia, reddit_posted, description, scraping_run_id, created_at, updated_at)
		VALUES (:id, :job_bank_id, :title, :employer, :location, :province, :city,
				:salary_min, :salary_max, :salary_type, :posting_date, :url, :is_tfw,
				:has_lmia, :reddit_posted, :description, :scraping_run_id, :created_at, :updated_at)
		ON CONFLICT (job_bank_id) DO UPDATE SET
			title = EXCLUDED.title,
			employer = EXCLUDED.employer,
			location = EXCLUDED.location,
			province = EXCLUDED.province,
			city = EXCLUDED.city,
			salary_min = EXCLUDED.salary_min,
			salary_max = EXCLUDED.salary_max,
			salary_type = EXCLUDED.salary_type,
			posting_date = EXCLUDED.posting_date,
			url = EXCLUDED.url,
			has_lmia = EXCLUDED.has_lmia,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
	`

	for _, posting := range postings {
		posting.ID = uuid.New().String()
		posting.CreatedAt = time.Now()
		posting.UpdatedAt = time.Now()
	}

	_, err = tx.NamedExec(query, postings)
	if err != nil {
		return fmt.Errorf("failed to insert job postings batch: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *jobBankRepository) GetJobPostingByJobBankID(jobBankID string) (*models.JobPosting, error) {
	var posting models.JobPosting
	query := `SELECT * FROM job_postings WHERE job_bank_id = $1`

	err := r.db.Get(&posting, query, jobBankID)
	if err != nil {
		return nil, err
	}

	return &posting, nil
}

func (r *jobBankRepository) GetJobPostingByID(id string) (*models.JobPosting, error) {
	var posting models.JobPosting
	query := `SELECT * FROM job_postings WHERE id = $1`

	err := r.db.Get(&posting, query, id)
	if err != nil {
		return nil, err
	}

	return &posting, nil
}

func (r *jobBankRepository) SearchJobPostingsByEmployer(employer string, limit int) ([]*models.JobPosting, error) {
	var postings []*models.JobPosting
	query := `
		SELECT * FROM job_postings
		WHERE employer ILIKE $1
		ORDER BY posting_date DESC, created_at DESC
	`

	searchTerm := "%" + employer + "%"

	if limit > 0 {
		query += " LIMIT $2"
		err := r.db.Select(&postings, query, searchTerm, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&postings, query, searchTerm)
		if err != nil {
			return nil, err
		}
	}

	return postings, nil
}

func (r *jobBankRepository) GetJobPostingsByLocation(city, province string, limit int) ([]*models.JobPosting, error) {
	var postings []*models.JobPosting
	query := `
		SELECT * FROM job_postings
		WHERE ($1 = '' OR city ILIKE $1 OR location ILIKE $1) AND ($2 = '' OR province ILIKE $2 OR location ILIKE $2)
		ORDER BY posting_date DESC, created_at DESC
	`

	citySearch := ""
	provinceSearch := ""
	if city != "" {
		citySearch = "%" + city + "%"
	}
	if province != "" {
		provinceSearch = "%" + province + "%"
	}

	if limit > 0 {
		query += " LIMIT $3"
		err := r.db.Select(&postings, query, citySearch, provinceSearch, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&postings, query, citySearch, provinceSearch)
		if err != nil {
			return nil, err
		}
	}

	return postings, nil
}

func (r *jobBankRepository) GetJobPostingsByScrapingRun(scrapingRunID string) ([]*models.JobPosting, error) {
	var postings []*models.JobPosting
	query := `SELECT * FROM job_postings WHERE scraping_run_id = $1 ORDER BY created_at DESC`

	err := r.db.Select(&postings, query, scrapingRunID)
	if err != nil {
		return nil, err
	}

	return postings, nil
}

func (r *jobBankRepository) GetRecentJobPostings(limit int) ([]*models.JobPosting, error) {
	var postings []*models.JobPosting
	query := `
		SELECT * FROM job_postings
		ORDER BY posting_date DESC, created_at DESC
	`

	if limit > 0 {
		query += " LIMIT $1"
		err := r.db.Select(&postings, query, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&postings, query)
		if err != nil {
			return nil, err
		}
	}

	return postings, nil
}

func (r *jobBankRepository) GetJobPostingsCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM job_postings`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *jobBankRepository) GetDistinctEmployersCount() (int, error) {
	var count int
	query := `SELECT COUNT(DISTINCT employer) FROM job_postings`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *jobBankRepository) GetEmployerJobCounts(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `
		SELECT employer, COUNT(*) as job_count, MIN(posting_date) as earliest_posting, MAX(posting_date) as latest_posting
		FROM job_postings
		GROUP BY employer
		ORDER BY job_count DESC
	`

	if limit > 0 {
		query += " LIMIT $1"
		rows, err := r.db.Query(query, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var employer string
			var jobCount int
			var earliestPosting, latestPosting *time.Time

			if err := rows.Scan(&employer, &jobCount, &earliestPosting, &latestPosting); err != nil {
				return nil, err
			}

			results = append(results, map[string]interface{}{
				"employer":         employer,
				"job_count":        jobCount,
				"earliest_posting": earliestPosting,
				"latest_posting":   latestPosting,
			})
		}
	} else {
		rows, err := r.db.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var employer string
			var jobCount int
			var earliestPosting, latestPosting *time.Time

			if err := rows.Scan(&employer, &jobCount, &earliestPosting, &latestPosting); err != nil {
				return nil, err
			}

			results = append(results, map[string]interface{}{
				"employer":         employer,
				"job_count":        jobCount,
				"earliest_posting": earliestPosting,
				"latest_posting":   latestPosting,
			})
		}
	}

	return results, nil
}

// CreateJobPostingsFromScraperData processes scraper data and creates job postings
func (r *jobBankRepository) CreateJobPostingsFromScraperData(scraperData []models.ScraperJobData, scrapingRunID string) ([]*models.JobPosting, error) {
	if len(scraperData) == 0 {
		return nil, nil
	}

	// Convert scraper data to job postings
	var postings []*models.JobPosting
	for _, data := range scraperData {
		posting := models.NewJobPostingFromScraperData(data, scrapingRunID)
		postings = append(postings, posting)
	}

	// Process in batches to avoid PostgreSQL parameter limit
	// Each job posting has ~19 parameters, so we'll use batches of 1000 to stay well under 65535 limit
	batchSize := 1000
	totalJobs := len(postings)
	var allNewJobIds []string
	
	for i := 0; i < totalJobs; i += batchSize {
		end := i + batchSize
		if end > totalJobs {
			end = totalJobs
		}
		
		batch := postings[i:end]
		newJobIds, err := r.createJobPostingsBatchUpdated(batch)
		if err != nil {
			return nil, fmt.Errorf("failed to insert batch %d-%d: %w", i, end, err)
		}
		
		allNewJobIds = append(allNewJobIds, newJobIds...)
		
		// Log progress for large datasets
		if totalJobs > batchSize {
			fmt.Printf("Processed batch %d-%d of %d jobs (%d new)\n", i+1, end, totalJobs, len(newJobIds))
		}
	}

	// Now fetch only the newly inserted jobs to return
	if len(allNewJobIds) == 0 {
		return []*models.JobPosting{}, nil
	}

	var newJobPostings []*models.JobPosting
	query, args, err := sqlx.In("SELECT * FROM job_postings WHERE id IN (?)", allNewJobIds)
	if err != nil {
		return nil, fmt.Errorf("failed to build query for new jobs: %w", err)
	}
	query = r.db.Rebind(query)
	
	if err := r.db.Select(&newJobPostings, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch new job postings: %w", err)
	}

	return newJobPostings, nil
}

// GetJobPostingByURL retrieves a job posting by its URL
func (r *jobBankRepository) GetJobPostingByURL(url string) (*models.JobPosting, error) {
	var posting models.JobPosting
	query := `SELECT * FROM job_postings WHERE url = $1`

	err := r.db.Get(&posting, query, url)
	if err != nil {
		return nil, err
	}

	return &posting, nil
}

// UpdateJobPostingRedditStatus updates the reddit_posted status for a job posting
func (r *jobBankRepository) UpdateJobPostingRedditStatus(id string, redditPosted bool) error {
	query := `UPDATE job_postings SET reddit_posted = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, redditPosted)
	return err
}

// UpdateJobRedditApprovalStatus updates the Reddit approval status and related fields for a job posting
func (r *jobBankRepository) UpdateJobRedditApprovalStatus(id string, status string, approvedBy string, approvedAt *time.Time, rejectionReason *string) error {
	query := `
		UPDATE job_postings 
		SET reddit_approval_status = $2, 
		    reddit_approved_by = $3, 
		    reddit_approved_at = $4, 
		    reddit_rejection_reason = $5,
		    updated_at = NOW() 
		WHERE id = $1`
	
	_, err := r.db.Exec(query, id, status, approvedBy, approvedAt, rejectionReason)
	return err
}

// GetJobPostingsNotPostedToReddit retrieves job postings that haven't been posted to Reddit
func (r *jobBankRepository) GetJobPostingsNotPostedToReddit(limit int) ([]*models.JobPosting, error) {
	var postings []*models.JobPosting
	query := `
		SELECT * FROM job_postings
		WHERE reddit_posted = FALSE
		ORDER BY posting_date DESC, created_at DESC
	`

	if limit > 0 {
		query += " LIMIT $1"
		err := r.db.Select(&postings, query, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&postings, query)
		if err != nil {
			return nil, err
		}
	}

	return postings, nil
}

// createJobPostingsBatchUpdated creates job postings in batch with updated schema
// Returns a list of newly inserted job IDs
func (r *jobBankRepository) createJobPostingsBatchUpdated(postings []*models.JobPosting) ([]string, error) {
	if len(postings) == 0 {
		return nil, nil
	}

	// Check which job_bank_ids and URLs already exist in the database
	var urls []string
	var jobBankIds []string
	for _, posting := range postings {
		urls = append(urls, posting.URL)
		if posting.JobBankID != nil && *posting.JobBankID != "" {
			jobBankIds = append(jobBankIds, *posting.JobBankID)
		}
	}
	
	existingUrls := make(map[string]bool)
	existingJobBankIds := make(map[string]bool)
	
	// Check existing URLs
	if len(urls) > 0 {
		query, args, err := sqlx.In("SELECT url FROM job_postings WHERE url IN (?)", urls)
		if err != nil {
			return nil, fmt.Errorf("failed to build query for existing URLs: %w", err)
		}
		query = r.db.Rebind(query)
		
		var existingUrlList []string
		if err := r.db.Select(&existingUrlList, query, args...); err != nil {
			return nil, fmt.Errorf("failed to check existing URLs: %w", err)
		}
		
		for _, url := range existingUrlList {
			existingUrls[url] = true
		}
	}
	
	// Check existing job_bank_ids
	if len(jobBankIds) > 0 {
		query, args, err := sqlx.In("SELECT job_bank_id FROM job_postings WHERE job_bank_id IN (?)", jobBankIds)
		if err != nil {
			return nil, fmt.Errorf("failed to build query for existing job_bank_ids: %w", err)
		}
		query = r.db.Rebind(query)
		
		var existingJobBankIdList []string
		if err := r.db.Select(&existingJobBankIdList, query, args...); err != nil {
			return nil, fmt.Errorf("failed to check existing job_bank_ids: %w", err)
		}
		
		for _, jobBankId := range existingJobBankIdList {
			existingJobBankIds[jobBankId] = true
		}
	}

	// Deduplicate by job_bank_id first (if available), then by URL
	duplicateMap := make(map[string]*models.JobPosting)
	var deduplicatedPostings []*models.JobPosting
	var newJobIds []string
	
	for _, posting := range postings {
		// Use job_bank_id as primary key if available, otherwise use URL
		var key string
		if posting.JobBankID != nil && *posting.JobBankID != "" {
			key = "job_bank_id:" + *posting.JobBankID
		} else {
			key = "url:" + posting.URL
		}
		
		if existing, exists := duplicateMap[key]; exists {
			// Keep the one with more recent updated_at
			if posting.UpdatedAt.After(existing.UpdatedAt) {
				duplicateMap[key] = posting
			}
		} else {
			duplicateMap[key] = posting
		}
	}
	
	for _, posting := range duplicateMap {
		// Assign ID and timestamps before checking if it's new
		posting.ID = uuid.New().String()
		posting.CreatedAt = time.Now()
		posting.UpdatedAt = time.Now()
		
		// Track which jobs are truly new (not in database)
		isExisting := existingUrls[posting.URL]
		if posting.JobBankID != nil && *posting.JobBankID != "" {
			isExisting = isExisting || existingJobBankIds[*posting.JobBankID]
		}
		
		if !isExisting {
			newJobIds = append(newJobIds, posting.ID)
		}
		
		deduplicatedPostings = append(deduplicatedPostings, posting)
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO job_postings (id, job_bank_id, title, employer, location, province, city,
								 salary_min, salary_max, salary_type, salary_raw, posting_date, url, is_tfw,
								 has_lmia, description, scraping_run_id, created_at, updated_at)
		VALUES (:id, :job_bank_id, :title, :employer, :location, :province, :city,
				:salary_min, :salary_max, :salary_type, :salary_raw, :posting_date, :url, :is_tfw,
				:has_lmia, :description, :scraping_run_id, :created_at, :updated_at)
		ON CONFLICT (url) DO UPDATE SET
			title = EXCLUDED.title,
			employer = EXCLUDED.employer,
			location = EXCLUDED.location,
			province = EXCLUDED.province,
			city = EXCLUDED.city,
			salary_min = EXCLUDED.salary_min,
			salary_max = EXCLUDED.salary_max,
			salary_type = EXCLUDED.salary_type,
			salary_raw = EXCLUDED.salary_raw,
			posting_date = EXCLUDED.posting_date,
			has_lmia = EXCLUDED.has_lmia,
			description = EXCLUDED.description,
			updated_at = EXCLUDED.updated_at
		WHERE job_postings.updated_at < EXCLUDED.updated_at
	`

	_, err = tx.NamedExec(query, deduplicatedPostings)
	if err != nil {
		return nil, fmt.Errorf("failed to insert job postings batch: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newJobIds, nil
}

// SearchJobPostingsAdvanced performs advanced search with multiple filters and pagination
func (r *jobBankRepository) SearchJobPostingsAdvanced(filters map[string]interface{}) ([]*models.JobPosting, int, error) {
	var postings []*models.JobPosting
	var totalCount int
	
	// Build base query
	whereClause := " WHERE 1=1"
	var args []interface{}
	argIndex := 1
	
	// Add search filter (searches across title, employer, location)
	if search, ok := filters["search"].(string); ok && search != "" {
		whereClause += fmt.Sprintf(" AND (title ILIKE $%d OR employer ILIKE $%d OR location ILIKE $%d)", argIndex, argIndex, argIndex)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argIndex++
	}
	
	// Add employer filter
	if employer, ok := filters["employer"].(string); ok && employer != "" {
		whereClause += fmt.Sprintf(" AND employer ILIKE $%d", argIndex)
		args = append(args, "%"+employer+"%")
		argIndex++
	}
	
	// Add title filter
	if title, ok := filters["title"].(string); ok && title != "" {
		whereClause += fmt.Sprintf(" AND title ILIKE $%d", argIndex)
		args = append(args, "%"+title+"%")
		argIndex++
	}
	
	// Add city filter
	if city, ok := filters["city"].(string); ok && city != "" {
		whereClause += fmt.Sprintf(" AND (city ILIKE $%d OR location ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+city+"%")
		argIndex++
	}
	
	// Add province filter
	if province, ok := filters["province"].(string); ok && province != "" {
		whereClause += fmt.Sprintf(" AND (province ILIKE $%d OR location ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+province+"%")
		argIndex++
	}
	
	// Add salary filter
	if salaryMin, ok := filters["salary_min"].(*float64); ok && salaryMin != nil {
		whereClause += fmt.Sprintf(" AND salary_min >= $%d", argIndex)
		args = append(args, *salaryMin)
		argIndex++
	}
	
	// Add days filter (jobs posted within X days)
	// Only apply the filter if days > 0, otherwise show all jobs
	if days, ok := filters["days"].(int); ok && days > 0 {
		whereClause += fmt.Sprintf(" AND posting_date >= NOW() - INTERVAL '%d days'", days)
	}
	
	// Add Reddit approval status filter
	if approvalStatus, ok := filters["reddit_approval_status"].(string); ok && approvalStatus != "" {
		whereClause += fmt.Sprintf(" AND reddit_approval_status = $%d", argIndex)
		args = append(args, approvalStatus)
		argIndex++
	}
	
	// Get total count
	countQuery := "SELECT COUNT(*) FROM job_postings" + whereClause
	err := r.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}
	
	// Add sorting
	sortBy := "posting_date"
	sortOrder := "DESC"
	if sort, ok := filters["sort_by"].(string); ok && sort != "" {
		// Validate sort field to prevent SQL injection
		validSorts := map[string]bool{
			"posting_date": true,
			"created_at":   true,
			"title":        true,
			"employer":     true,
			"salary_min":   true,
			"salary_max":   true,
		}
		if validSorts[sort] {
			sortBy = sort
		}
	}
	if order, ok := filters["sort_order"].(string); ok && (order == "ASC" || order == "DESC" || order == "asc" || order == "desc") {
		sortOrder = order
	}
	
	selectQuery := fmt.Sprintf(`
		SELECT id, job_bank_id, title, employer, location, province, city,
			   salary_min, salary_max, salary_type, salary_raw, posting_date, url, 
			   is_tfw, has_lmia, reddit_posted, reddit_approval_status, reddit_approved_by, 
			   reddit_approved_at, reddit_rejection_reason, description, scraping_run_id, 
			   created_at, updated_at
		FROM job_postings%s ORDER BY %s %s`, whereClause, sortBy, sortOrder)
	
	// Add pagination
	limit := 25
	offset := 0
	if l, ok := filters["limit"].(int); ok && l > 0 {
		limit = l
	}
	if o, ok := filters["offset"].(int); ok && o >= 0 {
		offset = o
	}
	
	selectQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)
	
	// Execute the query
	err = r.db.Select(&postings, selectQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute search query: %w", err)
	}
	
	return postings, totalCount, nil
}

// DeleteJobPostingsNotInScrapeRun deletes job postings that are not in the current scrape run
// This removes jobs that existed in previous scrapes but are no longer available on the job bank site
func (r *jobBankRepository) DeleteJobPostingsNotInScrapeRun(scrapingRunID string, currentJobBankIDs []string) (int, error) {
	if len(currentJobBankIDs) == 0 {
		// If no job bank IDs provided, don't delete anything to be safe
		return 0, nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete jobs that:
	// 1. Were created from previous scraping runs (not the current one)
	// 2. Have job_bank_id values that are not in the current scrape
	// 3. Are TFW/LMIA jobs (to avoid deleting manually added jobs)
	var deletedCount int64
	
	if len(currentJobBankIDs) > 0 {
		// Build the NOT IN clause dynamically based on the number of IDs
		query, args, err := sqlx.In(`
			DELETE FROM job_postings 
			WHERE scraping_run_id != ? 
			AND job_bank_id IS NOT NULL 
			AND job_bank_id NOT IN (?)
			AND is_tfw = true
			AND has_lmia = true
		`, scrapingRunID, currentJobBankIDs)
		if err != nil {
			return 0, fmt.Errorf("failed to build delete query: %w", err)
		}
		
		// Rebind the query for the specific database driver
		query = r.db.Rebind(query)
		
		result, err := tx.Exec(query, args...)
		if err != nil {
			return 0, fmt.Errorf("failed to delete orphaned job postings: %w", err)
		}
		
		deletedCount, err = result.RowsAffected()
		if err != nil {
			return 0, fmt.Errorf("failed to get deleted rows count: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return int(deletedCount), nil
}