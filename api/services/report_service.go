package services

import (
	"canada-hires/models"
	"canada-hires/repos"
	"fmt"
	"slices"
	"strings"
)

type CreateReportRequest struct {
	UserID          string  `json:"user_id"`
	BusinessName    string  `json:"business_name"`
	BusinessAddress string  `json:"business_address"`
	ReportSource    string  `json:"report_source"`
	ConfidenceLevel *int    `json:"confidence_level"`
	AdditionalNotes *string `json:"additional_notes"`
	IPAddress       *string `json:"ip_address"`
}

type ReportService interface {
	CreateReport(req *CreateReportRequest) (*models.Report, error)
	GetReportByID(id string) (*models.Report, error)
	GetAllReports(limit, offset int) ([]*models.Report, error)
	GetUserReports(userID string, limit, offset int) ([]*models.Report, error)
	GetBusinessReports(businessName string, limit, offset int) ([]*models.Report, error)
	GetReportsByStatus(status models.ReportStatus, limit, offset int) ([]*models.Report, error)
	UpdateReport(report *models.Report) error
	ApproveReport(reportID, moderatorID string, notes *string) error
	RejectReport(reportID, moderatorID string, notes *string) error
	FlagReport(reportID, moderatorID string, notes *string) error
	DeleteReport(reportID, userID string, isAdmin bool) error
}

type reportService struct {
	repo repos.ReportRepository
}

func NewReportService(repo repos.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) CreateReport(req *CreateReportRequest) (*models.Report, error) {
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	report := &models.Report{
		UserID:          req.UserID,
		BusinessName:    strings.TrimSpace(req.BusinessName),
		BusinessAddress: strings.TrimSpace(req.BusinessAddress),
		ReportSource:    req.ReportSource,
		ConfidenceLevel: req.ConfidenceLevel,
		AdditionalNotes: req.AdditionalNotes,
		IPAddress:       req.IPAddress,
		Status:          models.ReportPending,
	}

	if err := s.repo.Create(report); err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return report, nil
}

func (s *reportService) validateCreateRequest(req *CreateReportRequest) error {
	if req.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if strings.TrimSpace(req.BusinessName) == "" {
		return fmt.Errorf("business name is required")
	}
	if strings.TrimSpace(req.BusinessAddress) == "" {
		return fmt.Errorf("business address is required")
	}
	if req.ReportSource == "" {
		return fmt.Errorf("report source is required")
	}
	if !isValidReportSource(req.ReportSource) {
		return fmt.Errorf("invalid report source: must be 'employment', 'observation', or 'public_record'")
	}
	if req.ConfidenceLevel != nil && (*req.ConfidenceLevel < 1 || *req.ConfidenceLevel > 10) {
		return fmt.Errorf("confidence level must be between 1 and 10")
	}
	return nil
}


func isValidReportSource(source string) bool {
	validSources := []string{"employment", "observation", "public_record"}
	return slices.Contains(validSources, source)
}

func (s *reportService) GetReportByID(id string) (*models.Report, error) {
	if id == "" {
		return nil, fmt.Errorf("report ID is required")
	}

	report, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	return report, nil
}

func (s *reportService) GetAllReports(limit, offset int) ([]*models.Report, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	reports, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}

func (s *reportService) GetUserReports(userID string, limit, offset int) ([]*models.Report, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	reports, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reports: %w", err)
	}

	return reports, nil
}

func (s *reportService) GetBusinessReports(businessName string, limit, offset int) ([]*models.Report, error) {
	if strings.TrimSpace(businessName) == "" {
		return nil, fmt.Errorf("business name is required")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	reports, err := s.repo.GetByBusinessName(businessName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get business reports: %w", err)
	}

	return reports, nil
}

func (s *reportService) GetReportsByStatus(status models.ReportStatus, limit, offset int) ([]*models.Report, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	reports, err := s.repo.GetByStatus(status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by status: %w", err)
	}

	return reports, nil
}

func (s *reportService) UpdateReport(report *models.Report) error {
	if report.ID == "" {
		return fmt.Errorf("report ID is required")
	}

	if err := s.validateReport(report); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	report.BusinessName = strings.TrimSpace(report.BusinessName)
	report.BusinessAddress = strings.TrimSpace(report.BusinessAddress)

	if err := s.repo.Update(report); err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	return nil
}

func (s *reportService) validateReport(report *models.Report) error {
	if strings.TrimSpace(report.BusinessName) == "" {
		return fmt.Errorf("business name is required")
	}
	if strings.TrimSpace(report.BusinessAddress) == "" {
		return fmt.Errorf("business address is required")
	}
	if report.ReportSource == "" {
		return fmt.Errorf("report source is required")
	}
	if !isValidReportSource(report.ReportSource) {
		return fmt.Errorf("invalid report source: must be 'employment', 'observation', or 'public_record'")
	}
	if report.ConfidenceLevel != nil && (*report.ConfidenceLevel < 1 || *report.ConfidenceLevel > 10) {
		return fmt.Errorf("confidence level must be between 1 and 10")
	}
	return nil
}

func (s *reportService) ApproveReport(reportID, moderatorID string, notes *string) error {
	if reportID == "" {
		return fmt.Errorf("report ID is required")
	}
	if moderatorID == "" {
		return fmt.Errorf("moderator ID is required")
	}

	err := s.repo.UpdateStatus(reportID, models.ReportApproved, &moderatorID, notes)
	if err != nil {
		return fmt.Errorf("failed to approve report: %w", err)
	}

	return nil
}

func (s *reportService) RejectReport(reportID, moderatorID string, notes *string) error {
	if reportID == "" {
		return fmt.Errorf("report ID is required")
	}
	if moderatorID == "" {
		return fmt.Errorf("moderator ID is required")
	}

	err := s.repo.UpdateStatus(reportID, models.ReportRejected, &moderatorID, notes)
	if err != nil {
		return fmt.Errorf("failed to reject report: %w", err)
	}

	return nil
}

func (s *reportService) FlagReport(reportID, moderatorID string, notes *string) error {
	if reportID == "" {
		return fmt.Errorf("report ID is required")
	}
	if moderatorID == "" {
		return fmt.Errorf("moderator ID is required")
	}

	err := s.repo.UpdateStatus(reportID, models.ReportFlagged, &moderatorID, notes)
	if err != nil {
		return fmt.Errorf("failed to flag report: %w", err)
	}

	return nil
}

func (s *reportService) DeleteReport(reportID, userID string, isAdmin bool) error {
	if reportID == "" {
		return fmt.Errorf("report ID is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Get the report to check ownership
	report, err := s.repo.GetByID(reportID)
	if err != nil {
		return fmt.Errorf("failed to get report: %w", err)
	}

	// Only allow report owner or admin to delete
	if !isAdmin && report.UserID != userID {
		return fmt.Errorf("unauthorized: can only delete your own reports")
	}

	err = s.repo.Delete(reportID)
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	return nil
}
