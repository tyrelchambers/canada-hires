package models

import "time"

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportApproved ReportStatus = "approved"
	ReportRejected ReportStatus = "rejected"
	ReportFlagged  ReportStatus = "flagged"
)

type Report struct {
	ID              string       `json:"id" db:"id"`
	UserID          string       `json:"user_id" db:"user_id"`
	BusinessName    string       `json:"business_name" db:"business_name"`
	BusinessAddress string       `json:"business_address" db:"business_address"`
	ReportSource    string       `json:"report_source" db:"report_source"`
	ConfidenceLevel *int         `json:"confidence_level" db:"confidence_level"`
	AdditionalNotes *string      `json:"additional_notes" db:"additional_notes"`
	Status          ReportStatus `json:"status" db:"status"`
	ModeratedBy     *string      `json:"moderated_by" db:"moderated_by"`
	ModerationNotes *string      `json:"moderation_notes" db:"moderation_notes"`
	IPAddress       *string      `json:"ip_address" db:"ip_address"`
	CreatedAt       time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" db:"updated_at"`
}
