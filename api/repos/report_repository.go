package repos

import (
	"canada-hires/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReportFilters struct {
	Query    string
	City     string
	Province string
	Year     string
}

type ReportRepository interface {
	Create(report *models.Report) error
	GetByID(id string) (*models.Report, error)
	GetByUserID(userID string, limit, offset int) ([]*models.Report, error)
	GetByBusinessName(businessName string, limit, offset int) ([]*models.Report, error)
	GetAll(limit, offset int) ([]*models.Report, error)
	GetWithFilters(filters ReportFilters, limit, offset int) ([]*models.Report, error)
	GetByAddress(address string) ([]*models.Report, error)
	GetReportsGroupedByAddress(filters *ReportFilters, limit, offset int) ([]*models.ReportsByAddress, error)
	Update(report *models.Report) error
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
			confidence_level, tfw_ratio, additional_notes, ip_address, created_at, updated_at)
		VALUES (:id, :user_id, :business_name, :business_address, :report_source,
			:confidence_level, :tfw_ratio, :additional_notes, :ip_address, :created_at, :updated_at)
	`

	report.ID = uuid.New().String()
	report.CreatedAt = time.Now().UTC()
	report.UpdatedAt = time.Now().UTC()

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

func (r *reportRepository) GetWithFilters(filters ReportFilters, limit, offset int) ([]*models.Report, error) {
	var reports []*models.Report
	var args []interface{}
	var conditions []string
	argCount := 0

	// Base query
	query := `SELECT * FROM reports WHERE 1=1`

	// Add business name/query filter
	if filters.Query != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("business_name ILIKE $%d", argCount))
		args = append(args, "%"+filters.Query+"%")
	}

	// Add city filter (search in business_address)
	if filters.City != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("business_address ILIKE $%d", argCount))
		args = append(args, "%"+filters.City+"%")
	}

	// Add province filter (search in business_address)
	if filters.Province != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("business_address ILIKE $%d", argCount))
		args = append(args, "%"+filters.Province+"%")
	}

	// Add year filter
	if filters.Year != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM created_at) = $%d", argCount))
		args = append(args, filters.Year)
	}

	// Combine conditions
	for _, condition := range conditions {
		query += " AND " + condition
	}

	// Add ordering and pagination
	argCount++
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCount)
	args = append(args, limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	err := r.db.Select(&reports, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports with filters: %w", err)
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
			tfw_ratio = :tfw_ratio,
			additional_notes = :additional_notes,
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

func (r *reportRepository) GetByAddress(address string) ([]*models.Report, error) {
	var reports []*models.Report
	query := `SELECT * FROM reports WHERE business_address = $1 ORDER BY created_at DESC`

	err := r.db.Select(&reports, query, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by address: %w", err)
	}

	return reports, nil
}

func (r *reportRepository) GetReportsGroupedByAddress(filters *ReportFilters, limit, offset int) ([]*models.ReportsByAddress, error) {
	var grouped []*models.ReportsByAddress

	// Build the WHERE conditions
	conditions := []string{"1=1"} // Start with a neutral condition
	args := []interface{}{}
	argCount := 0

	// Apply filters only if filters is not nil
	if filters != nil {
		// Add search query filter (searches business name)
		if filters.Query != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("LOWER(business_name) LIKE LOWER($%d)", argCount))
			args = append(args, "%"+filters.Query+"%")
		}

		// Add city filter (searches business address)
		if filters.City != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("LOWER(business_address) LIKE LOWER($%d)", argCount))
			args = append(args, "%"+filters.City+"%")
		}

		// Add province filter (searches business address)
		if filters.Province != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("LOWER(business_address) LIKE LOWER($%d)", argCount))
			args = append(args, "%"+filters.Province+"%")
		}

		// Add year filter
		if filters.Year != "" {
			argCount++
			conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM created_at) = $%d", argCount))
			args = append(args, filters.Year)
		}
	}

	// Build the final query
	whereClause := strings.Join(conditions, " AND ")
	query := fmt.Sprintf(`
		SELECT
			business_name,
			business_address,
			COUNT(*) as report_count,
			AVG(COALESCE(confidence_level, 5)) as confidence_level,
			SUM(CASE WHEN tfw_ratio = 'few' THEN 1 ELSE 0 END) as tfw_ratio_few,
			SUM(CASE WHEN tfw_ratio = 'many' THEN 1 ELSE 0 END) as tfw_ratio_many,
			SUM(CASE WHEN tfw_ratio = 'most' THEN 1 ELSE 0 END) as tfw_ratio_most,
			SUM(CASE WHEN tfw_ratio = 'all' THEN 1 ELSE 0 END) as tfw_ratio_all,
			MAX(created_at) as latest_report
		FROM reports
		WHERE %s
		GROUP BY business_address, business_name
		ORDER BY report_count DESC, latest_report DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount+1, argCount+2)

	// Add limit and offset to args
	args = append(args, limit, offset)

	err := r.db.Select(&grouped, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get filtered reports grouped by address: %w", err)
	}

	return grouped, nil
}
