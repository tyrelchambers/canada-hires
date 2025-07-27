package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeNow returns current time - helper for consistent time handling
func TimeNow() time.Time {
	return time.Now()
}

type JobPosting struct {
	ID           string     `json:"id" db:"id"`
	JobBankID    *string    `json:"job_bank_id" db:"job_bank_id"`       // Unique ID from Job Bank (nullable)
	Title        string     `json:"title" db:"title"`                   // Job title
	Employer     string     `json:"employer" db:"employer"`             // Company/employer name
	Location     string     `json:"location" db:"location"`             // City, Province
	Province     *string    `json:"province" db:"province"`             // Parsed province
	City         *string    `json:"city" db:"city"`                     // Parsed city
	SalaryMin    *float64   `json:"salary_min" db:"salary_min"`         // Minimum salary
	SalaryMax    *float64   `json:"salary_max" db:"salary_max"`         // Maximum salary
	SalaryType   *string    `json:"salary_type" db:"salary_type"`       // hourly, weekly, monthly, yearly
	SalaryRaw    *string    `json:"salary_raw" db:"salary_raw"`         // Original salary string from scraper
	PostingDate  *time.Time `json:"posting_date" db:"posting_date"`     // When job was posted
	URL          string     `json:"url" db:"url"`                       // Link to job posting
	IsTFW        bool       `json:"is_tfw" db:"is_tfw"`                 // Whether this is a TFW position
	HasLMIA      bool       `json:"has_lmia" db:"has_lmia"`             // Whether job has LMIA flag
	Description  *string    `json:"description" db:"description"`       // Job description if available
	ScrapingRunID string    `json:"scraping_run_id" db:"scraping_run_id"` // Reference to scraping session
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// ScraperJobData represents the data structure from your scraper
type ScraperJobData struct {
	JobTitle   string  `json:"jobTitle"`
	Business   string  `json:"business"`
	Salary     string  `json:"salary"`
	Location   string  `json:"location"`
	JobUrl     string  `json:"jobUrl"`
	Date       string  `json:"date"`
	JobBankID  *string `json:"jobBankId,omitempty"`
}

// NewJobPostingFromScraperData creates a JobPosting from scraper data
func NewJobPostingFromScraperData(scraperData ScraperJobData, scrapingRunID string) *JobPosting {
	job := &JobPosting{
		Title:         truncateString(scraperData.JobTitle, 500),     // Keep title reasonable length
		Employer:      truncateString(scraperData.Business, 500),    // Match DB constraint
		Location:      truncateString(scraperData.Location, 200),    // Match DB constraint
		URL:           scraperData.JobUrl,                           // URLs can be long
		SalaryRaw:     &scraperData.Salary,
		ScrapingRunID: scrapingRunID,
		IsTFW:         true, // All jobs from TFW scraper are TFW positions
		HasLMIA:       true, // All jobs from TFW scraper have LMIA
		JobBankID:     scraperData.JobBankID,                        // Job Bank ID from scraper
	}

	// Parse posting date
	if postingDate, err := parseScraperDate(scraperData.Date); err == nil {
		job.PostingDate = &postingDate
	}

	// Parse salary information
	job.parseSalary()

	// Parse location into city and province
	job.parseLocation()

	return job
}

// truncateString safely truncates a string to a maximum length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// parseSalary extracts salary min, max, and type from the raw salary string
func (jp *JobPosting) parseSalary() {
	if jp.SalaryRaw == nil {
		return
	}

	salaryStr := strings.TrimSpace(*jp.SalaryRaw)
	
	// Remove currency symbols and common prefixes
	salaryStr = strings.ReplaceAll(salaryStr, "$", "")
	salaryStr = strings.ReplaceAll(salaryStr, ",", "")
	
	// Determine salary type based on keywords
	salaryType := "hourly" // default
	if strings.Contains(strings.ToLower(salaryStr), "yearly") || strings.Contains(strings.ToLower(salaryStr), "annual") {
		salaryType = "yearly"
	} else if strings.Contains(strings.ToLower(salaryStr), "monthly") {
		salaryType = "monthly"
	} else if strings.Contains(strings.ToLower(salaryStr), "weekly") {
		salaryType = "weekly"
	}
	jp.SalaryType = &salaryType

	// Extract numbers using regex
	re := regexp.MustCompile(`(\d+\.?\d*)\s*to\s*(\d+\.?\d*)`)
	matches := re.FindStringSubmatch(salaryStr)
	
	if len(matches) == 3 {
		// Range salary: "20.00 to 25.00"
		if min, err := strconv.ParseFloat(matches[1], 64); err == nil {
			jp.SalaryMin = &min
		}
		if max, err := strconv.ParseFloat(matches[2], 64); err == nil {
			jp.SalaryMax = &max
		}
	} else {
		// Single salary value
		singleRe := regexp.MustCompile(`(\d+\.?\d*)`)
		if match := singleRe.FindString(salaryStr); match != "" {
			if salary, err := strconv.ParseFloat(match, 64); err == nil {
				jp.SalaryMin = &salary
				jp.SalaryMax = &salary
			}
		}
	}
}

// parseLocation extracts city and province from location string
func (jp *JobPosting) parseLocation() {
	if jp.Location == "" {
		return
	}

	// Split location by common separators
	parts := strings.Split(jp.Location, ",")
	if len(parts) >= 2 {
		city := strings.TrimSpace(parts[0])
		province := strings.TrimSpace(parts[len(parts)-1])
		
		// Convert province names to codes for consistency
		province = normalizeProvince(province)
		
		if city != "" {
			city = truncateString(city, 150) // Match new DB constraint
			jp.City = &city
		}
		if province != "" {
			jp.Province = &province
		}
	}
}

// normalizeProvince converts full province names to standard codes or keeps them as-is
func normalizeProvince(province string) string {
	provinceMap := map[string]string{
		"alberta":                     "AB",
		"british columbia":            "BC",
		"manitoba":                    "MB",
		"new brunswick":               "NB",
		"newfoundland and labrador":   "NL",
		"northwest territories":       "NT",
		"nova scotia":                 "NS",
		"nunavut":                     "NU",
		"ontario":                     "ON",
		"prince edward island":        "PE",
		"quebec":                      "QC",
		"saskatchewan":                "SK",
		"yukon":                       "YT",
	}
	
	// Try to match the full name first
	lowerProvince := strings.ToLower(strings.TrimSpace(province))
	if code, exists := provinceMap[lowerProvince]; exists {
		return code
	}
	
	// If it's already a code or unrecognized, return as-is (truncated if too long)
	if len(province) > 50 {
		return province[:50]
	}
	return province
}

// parseScraperDate parses the date format from scraper
func parseScraperDate(dateStr string) (time.Time, error) {
	// Try common date formats
	formats := []string{
		"January 2, 2006",
		"Jan 2, 2006", 
		"2006-01-02",
		"January 02, 2006",
		"Jan 02, 2006",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, nil
}

type JobScrapingRun struct {
	ID              string     `json:"id" db:"id"`
	Status          string     `json:"status" db:"status"`               // running, completed, failed
	StartedAt       time.Time  `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage    *string    `json:"error_message" db:"error_message"`
	TotalPages      int        `json:"total_pages" db:"total_pages"`     // Total pages scraped
	JobsScraped     int        `json:"jobs_scraped" db:"jobs_scraped"`   // Number of jobs found
	JobsStored      int        `json:"jobs_stored" db:"jobs_stored"`     // Number successfully stored
	LastPageScraped int        `json:"last_page_scraped" db:"last_page_scraped"` // For resuming
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}