package models

import "time"

type Report struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	BusinessName    string    `json:"business_name" db:"business_name"`
	BusinessAddress string    `json:"business_address" db:"business_address"`
	ReportSource    string    `json:"report_source" db:"report_source"`
	ConfidenceLevel *int      `json:"confidence_level" db:"confidence_level"` // Deprecated: use TFWRatio
	TFWRatio        *string   `json:"tfw_ratio" db:"tfw_ratio"`
	AdditionalNotes *string   `json:"additional_notes" db:"additional_notes"`
	IPAddress       *string   `json:"ip_address" db:"ip_address"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

type ReportsByAddress struct {
	BusinessName    string    `json:"business_name" db:"business_name"`
	BusinessAddress string    `json:"business_address" db:"business_address"`
	ReportCount     int       `json:"report_count" db:"report_count"`
	ConfidenceLevel float64   `json:"confidence_level" db:"confidence_level"` // Deprecated: use TFWRatioDistribution
	TFWRatioFew     int       `json:"tfw_ratio_few" db:"tfw_ratio_few"`
	TFWRatioMany    int       `json:"tfw_ratio_many" db:"tfw_ratio_many"`
	TFWRatioMost    int       `json:"tfw_ratio_most" db:"tfw_ratio_most"`
	TFWRatioAll     int       `json:"tfw_ratio_all" db:"tfw_ratio_all"`
	LatestReport    time.Time `json:"latest_report" db:"latest_report"`
}
