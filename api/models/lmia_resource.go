package models

import (
	"time"
)

type LMIAResource struct {
	ID           string     `json:"id" db:"id"`
	ResourceID   string     `json:"resource_id" db:"resource_id"`
	Name         string     `json:"name" db:"name"`
	Quarter      string     `json:"quarter" db:"quarter"`
	Year         int        `json:"year" db:"year"`
	URL          string     `json:"url" db:"url"`
	Format       string     `json:"format" db:"format"`
	Language     string     `json:"language" db:"language"`
	SizeBytes    *int64     `json:"size_bytes" db:"size_bytes"`
	LastModified *time.Time `json:"last_modified" db:"last_modified"`
	DatePublished *time.Time `json:"date_published" db:"date_published"`
	DownloadedAt *time.Time `json:"downloaded_at" db:"downloaded_at"`
	ProcessedAt  *time.Time `json:"processed_at" db:"processed_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type LMIAEmployer struct {
	ID                 string    `json:"id" db:"id"`
	ResourceID         string    `json:"resource_id" db:"resource_id"`
	Quarter            string    `json:"quarter" db:"quarter"`                         // Quarter from resource
	Year               int       `json:"year" db:"year"`                               // Year from filename/resource
	
	// These are the ONLY 8 columns from the actual LMIA CSV files
	ProvinceTerritory  *string   `json:"province_territory" db:"province_territory"`   // "Province/Territory"
	ProgramStream      *string   `json:"program_stream" db:"program_stream"`           // "Program Stream" 
	Employer           string    `json:"employer" db:"employer"`                       // "Employer"
	Address            *string   `json:"address" db:"address"`                         // "Address"
	Occupation         *string   `json:"occupation" db:"occupation"`                   // "Occupation"
	IncorporateStatus  *string   `json:"incorporate_status" db:"incorporate_status"`   // "Incorporate Status"
	ApprovedLMIAs      *int      `json:"approved_lmias" db:"approved_lmias"`           // "Approved LMIAs"
	ApprovedPositions  *int      `json:"approved_positions" db:"approved_positions"`   // "Approved Positions"
	
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CronJob struct {
	ID                string     `json:"id" db:"id"`
	JobName          string     `json:"job_name" db:"job_name"`
	Status           string     `json:"status" db:"status"`
	StartedAt        time.Time  `json:"started_at" db:"started_at"`
	CompletedAt      *time.Time `json:"completed_at" db:"completed_at"`
	ErrorMessage     *string    `json:"error_message" db:"error_message"`
	ResourcesProcessed int      `json:"resources_processed" db:"resources_processed"`
	RecordsProcessed int        `json:"records_processed" db:"records_processed"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
}

type LMIAGeographicSummary struct {
	Province         string `json:"province" db:"province_territory"`
	TotalEmployers   int    `json:"total_employers" db:"total_employers"`
	TotalLMIAs       int    `json:"total_lmias" db:"total_lmias"`
	TotalPositions   int    `json:"total_positions" db:"total_positions"`
	Year            int     `json:"year" db:"year"`
}