package models

import "time"

type Report struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	BusinessName    string    `json:"business_name" db:"business_name"`
	BusinessAddress string    `json:"business_address" db:"business_address"`
	ReportSource    string    `json:"report_source" db:"report_source"`
	ConfidenceLevel *int      `json:"confidence_level" db:"confidence_level"`
	AdditionalNotes *string   `json:"additional_notes" db:"additional_notes"`
	IPAddress       *string   `json:"ip_address" db:"ip_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ReportsByAddress struct {
	BusinessName    string    `json:"business_name" db:"business_name"`
	BusinessAddress string    `json:"business_address" db:"business_address"`
	ReportCount     int       `json:"report_count" db:"report_count"`
	ConfidenceLevel float64   `json:"confidence_level" db:"confidence_level"`
	LatestReport    time.Time `json:"latest_report" db:"latest_report"`
}
