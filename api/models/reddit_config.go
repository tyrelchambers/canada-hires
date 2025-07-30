package models

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type RedditConfig struct {
	ID                string    `json:"id" db:"id"`
	Subreddit         string    `json:"subreddit" db:"subreddit"`                     // Target subreddit (e.g., "canadajobs")
	PostTitleTemplate string    `json:"post_title_template" db:"post_title_template"` // Template for post title
	PostBodyTemplate  string    `json:"post_body_template" db:"post_body_template"`   // Template for post body
	IsEnabled         bool      `json:"is_enabled" db:"is_enabled"`                   // Enable/disable Reddit posting
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// Default Reddit configuration
func DefaultRedditConfig() *RedditConfig {
	return &RedditConfig{
		Subreddit:         os.Getenv("REDDIT_SUBREDDIT"),
		PostTitleTemplate: "🇨🇦 New LMIA Job: {{.Title}} at {{.Employer}} - {{.Location}}",
		PostBodyTemplate: `Apparently, the {{.Employer}} is unable to find anyone in {{.Location}} or surrounding area to work for them. They are so in need of employees that they are applying for a Labour Market Impact Assessment (LMIA).

If the LMIA application is successful, {{.Employer}} will be able to bring in Temporary Foreign Workers to do the jobs.

**Position:** {{.Title}}

**Employer:** {{.Employer}}

**Location:** {{.Location}}

{{if .SalaryRaw}}**Salary:** {{.SalaryRaw}}{{end}}

{{if .PostingDate}}**Posted:** {{.PostingDate}}{{end}}

**Job Details:** [View on Job Bank]({{.URL}})

Apply, and if you don't hear back, follow the links on the ad to report the business.

---
*This posting was automatically detected from Government of Canada Job Bank TFW listings. Data provided by JobWatch Canada for transparency in hiring practices.*

*See more LMIA listings at [JobWatch Canada](https://jobwatchcanada.com)*`,
		IsEnabled: true,
	}
}

// RedditPostData contains the data needed to create a Reddit post
type RedditPostData struct {
	Title      string
	Body       string
	Subreddit  string
	JobPosting *JobPosting
}

// GeneratePostData creates Reddit post data from a job posting using the config templates
func (rc *RedditConfig) GeneratePostData(job *JobPosting) *RedditPostData {
	if job == nil {
		return nil
	}

	title := rc.processTemplate(rc.PostTitleTemplate, job)
	body := rc.processTemplate(rc.PostBodyTemplate, job)

	return &RedditPostData{
		Title:      title,
		Body:       body,
		Subreddit:  rc.Subreddit,
		JobPosting: job,
	}
}

// buildSalaryRangeText formats salary information from min, max, and type
func buildSalaryRangeText(min *float64, max *float64, salaryType *string) string {
	if min == nil && max == nil {
		return ""
	}
	
	// Default type to "hourly" if not specified
	typeStr := "per hour"
	if salaryType != nil {
		switch *salaryType {
		case "yearly":
			typeStr = "per year"
		case "monthly":
			typeStr = "per month"
		case "weekly":
			typeStr = "per week"
		case "biweekly":
			typeStr = "biweekly"
		case "hourly":
			typeStr = "per hour"
		default:
			typeStr = "per hour"
		}
	}
	
	// Format salary amounts
	if min != nil && max != nil {
		if *min == *max {
			// Same min and max, show single salary
			return fmt.Sprintf("$%.2f %s", *min, typeStr)
		} else {
			// Different min and max, show range
			return fmt.Sprintf("$%.2f - $%.2f %s", *min, *max, typeStr)
		}
	} else if min != nil {
		// Only minimum specified
		return fmt.Sprintf("$%.2f+ %s", *min, typeStr)
	} else if max != nil {
		// Only maximum specified
		return fmt.Sprintf("Up to $%.2f %s", *max, typeStr)
	}
	
	return ""
}

// processTemplate replaces template placeholders with job data
func (rc *RedditConfig) processTemplate(template string, job *JobPosting) string {
	result := template

	// Replace basic fields
	result = strings.ReplaceAll(result, "{{.Title}}", job.Title)
	result = strings.ReplaceAll(result, "{{.Employer}}", job.Employer)
	result = strings.ReplaceAll(result, "{{.Location}}", job.Location)
	result = strings.ReplaceAll(result, "{{.URL}}", job.URL)

	// Replace optional fields with conditional logic
	salaryText := buildSalaryRangeText(job.SalaryMin, job.SalaryMax, job.SalaryType)
	if salaryText != "" {
		result = strings.ReplaceAll(result, "{{if .SalaryRaw}}**Salary:** {{.SalaryRaw}}{{end}}", "**Salary:** "+salaryText)
	} else {
		result = strings.ReplaceAll(result, "{{if .SalaryRaw}}**Salary:** {{.SalaryRaw}}{{end}}", "")
	}

	if job.PostingDate != nil {
		dateStr := job.PostingDate.Format("January 2, 2006")
		result = strings.ReplaceAll(result, "{{if .PostingDate}}**Posted:** {{.PostingDate}}{{end}}", "**Posted:** "+dateStr)
	} else {
		result = strings.ReplaceAll(result, "{{if .PostingDate}}**Posted:** {{.PostingDate}}{{end}}", "")
	}

	// Clean up any double newlines from removed optional fields
	result = strings.ReplaceAll(result, "\n\n\n", "\n\n")

	return result
}

// ValidateConfig checks if the Reddit configuration is valid
func (rc *RedditConfig) ValidateConfig() error {
	if rc.Subreddit == "" {
		return fmt.Errorf("subreddit cannot be empty")
	}
	if rc.PostTitleTemplate == "" {
		return fmt.Errorf("post title template cannot be empty")
	}
	if rc.PostBodyTemplate == "" {
		return fmt.Errorf("post body template cannot be empty")
	}
	return nil
}
