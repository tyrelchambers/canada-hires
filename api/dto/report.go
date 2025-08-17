package dto

import (
	"canada-hires/models"
	"time"
)

type CreateReportRequest struct {
	BusinessName    string  `json:"business_name" validate:"required,min=1,max=255"`
	BusinessAddress string  `json:"business_address" validate:"required,min=1"`
	ReportSource    string  `json:"report_source" validate:"required,oneof=employment observation public_record"`
	ConfidenceLevel *int    `json:"confidence_level" validate:"omitempty,gte=1,lte=10"`
	AdditionalNotes *string `json:"additional_notes" validate:"omitempty,max=1000"`
}

type UpdateReportRequest struct {
	BusinessName    string  `json:"business_name" validate:"required,min=1,max=255"`
	BusinessAddress string  `json:"business_address" validate:"required,min=1"`
	ReportSource    string  `json:"report_source" validate:"required,oneof=employment observation public_record"`
	ConfidenceLevel *int    `json:"confidence_level" validate:"omitempty,gte=1,lte=10"`
	AdditionalNotes *string `json:"additional_notes" validate:"omitempty,max=1000"`
}

type ReportResponse struct {
	ID              string                 `json:"id"`
	UserID          string                 `json:"user_id"`
	BusinessName    string                 `json:"business_name"`
	BusinessAddress string                 `json:"business_address"`
	ReportSource    string                 `json:"report_source"`
	ConfidenceLevel *int                   `json:"confidence_level"`
	AdditionalNotes *string                `json:"additional_notes"`
	Status          models.ReportStatus    `json:"status"`
	ModeratedBy     *string                `json:"moderated_by,omitempty"`
	ModerationNotes *string                `json:"moderation_notes,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

type ReportListResponse struct {
	Reports    []*ReportResponse `json:"reports"`
	Pagination PaginationInfo    `json:"pagination"`
}

type PaginationInfo struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total,omitempty"` // Optional - can be expensive to calculate
}

type ModerationRequest struct {
	Notes *string `json:"notes" validate:"omitempty,max=500"`
}

type ReportStatsResponse struct {
	TotalReports    int `json:"total_reports"`
	PendingReports  int `json:"pending_reports"`
	ApprovedReports int `json:"approved_reports"`
	RejectedReports int `json:"rejected_reports"`
	FlaggedReports  int `json:"flagged_reports"`
}

// Helper function to convert model to DTO
func ToReportResponse(report *models.Report) *ReportResponse {
	return &ReportResponse{
		ID:              report.ID,
		UserID:          report.UserID,
		BusinessName:    report.BusinessName,
		BusinessAddress: report.BusinessAddress,
		ReportSource:    report.ReportSource,
		ConfidenceLevel: report.ConfidenceLevel,
		AdditionalNotes: report.AdditionalNotes,
		Status:          report.Status,
		ModeratedBy:     report.ModeratedBy,
		ModerationNotes: report.ModerationNotes,
		CreatedAt:       report.CreatedAt,
		UpdatedAt:       report.UpdatedAt,
	}
}

// Helper function to convert multiple models to DTOs
func ToReportListResponse(reports []*models.Report, limit, offset int) *ReportListResponse {
	responses := make([]*ReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = ToReportResponse(report)
	}

	return &ReportListResponse{
		Reports: responses,
		Pagination: PaginationInfo{
			Limit:  limit,
			Offset: offset,
		},
	}
}