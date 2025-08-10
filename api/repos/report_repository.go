package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReportRepository interface {
	Create(report *models.Report) error
	GetByID(id string) (*models.Report, error)
	GetByUserID(userID string, limit, offset int) ([]*models.Report, error)
	GetByBusinessName(businessName string, limit, offset int) ([]*models.Report, error)
	GetByStatus(status models.ReportStatus, limit, offset int) ([]*models.Report, error)
	GetAll(limit, offset int) ([]*models.Report, error)
	Update(report *models.Report) error
	UpdateStatus(id string, status models.ReportStatus, moderatorID *string, notes *string) error
	Delete(id string) error
}

type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *models.Report) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO reports (id, user_id, business_name, business_address, report_source, 
			confidence_level, additional_notes, status, moderated_by, moderation_notes, 
			ip_address, created_at, updated_at)
		VALUES (:id, :user_id, :business_name, :business_address, :report_source, 
			:confidence_level, :additional_notes, :status, :moderated_by, :moderation_notes, 
			:ip_address, :created_at, :updated_at)
	`

	report.ID = uuid.New().String()
	report.CreatedAt = time.Now().UTC()
	report.UpdatedAt = time.Now().UTC()
	
	if report.Status == "" {
		report.Status = models.ReportPending
	}

	_, err = tx.NamedExec(query, report)
	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *reportRepository) GetByID(id string) (*models.Report, error) {
	var report models.Report
	query := `SELECT * FROM reports WHERE id = $1`

	err := r.db.Get(&report, query, id)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *reportRepository) GetAll(limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	query := `SELECT * FROM reports ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	err := r.db.Select(&reports, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}

func (r *reportRepository) GetByUserID(userID string, limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	query := `SELECT * FROM reports WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	err := r.db.Select(&reports, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by user ID: %w", err)
	}

	return reports, nil
}

func (r *reportRepository) GetByBusinessName(businessName string, limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	query := `SELECT * FROM reports WHERE business_name ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	err := r.db.Select(&reports, query, "%"+businessName+"%", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by business name: %w", err)
	}

	return reports, nil
}

func (r *reportRepository) GetByStatus(status models.ReportStatus, limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	query := `SELECT * FROM reports WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	err := r.db.Select(&reports, query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by status: %w", err)
	}

	return reports, nil
}

func (r *reportRepository) Update(report *models.Report) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE reports SET 
			business_name = :business_name,
			business_address = :business_address,
			report_source = :report_source,
			confidence_level = :confidence_level,
			additional_notes = :additional_notes,
			status = :status,
			updated_at = NOW()
		WHERE id = :id
	`

	report.UpdatedAt = time.Now().UTC()

	_, err = tx.NamedExec(query, report)
	if err != nil {
		return fmt.Errorf("failed to update report: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *reportRepository) UpdateStatus(id string, status models.ReportStatus, moderatorID *string, notes *string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE reports SET 
			status = $1,
			moderated_by = $2,
			moderation_notes = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	_, err = tx.Exec(query, status, moderatorID, notes, id)
	if err != nil {
		return fmt.Errorf("failed to update report status: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *reportRepository) Delete(id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM reports WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete report: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
