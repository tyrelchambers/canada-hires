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
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	BusinessName    string    `json:"business_name"`
	BusinessAddress string    `json:"business_address"`
	ReportSource    string    `json:"report_source"`
	ConfidenceLevel *int      `json:"confidence_level"`
	AdditionalNotes *string   `json:"additional_notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
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