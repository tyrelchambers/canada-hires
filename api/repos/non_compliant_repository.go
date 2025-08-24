package repos

import (
	"canada-hires/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type NonCompliantRepository interface {
	// Employers
	CreateEmployer(employer *models.NonCompliantEmployer) error
	CreateEmployersBatch(employers []models.NonCompliantEmployer) error
	GetEmployerByID(id string) (*models.NonCompliantEmployer, error)
	GetEmployersWithReasons(limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error)
	GetEmployersCount() (int, error)
	UpdateEmployer(employer *models.NonCompliantEmployer) error
	DeleteEmployer(id string) error

	// Reasons
	CreateReason(reason *models.NonCompliantReason) error
	GetReasonByCode(code string) (*models.NonCompliantReason, error)
	GetAllReasons() ([]models.NonCompliantReason, error)
	UpsertReason(code, description string) (*models.NonCompliantReason, error)

	// Employer-Reason relationships
	AddEmployerReason(employerID string, reasonID int) error
	RemoveEmployerReason(employerID string, reasonID int) error
	GetEmployerReasons(employerID string) ([]models.NonCompliantReason, error)

	// Bulk operations
	CreateEmployersWithReasons(data []models.ScraperNonCompliantData) error
	UpsertEmployersWithReasons(data []models.ScraperNonCompliantData) error
	ClearAllNonCompliantData() error

	// Stats
	GetLatestScrapedDate() (*time.Time, error)
	GetTotalEmployersCount() (int, error)

	// Geolocation methods
	GetEmployersWithoutPostalCodes() ([]models.NonCompliantEmployer, error)
	UpdateEmployerPostalCode(employerID, postalCode string) error
	GetLocationsByPostalCode(limit int) ([]models.NonCompliantPostalCodeLocation, error)
	GetEmployersByPostalCode(postalCode string, limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error)
}

type nonCompliantRepository struct {
	db *sqlx.DB
}

func NewNonCompliantRepository(db *sqlx.DB) NonCompliantRepository {
	return &nonCompliantRepository{db: db}
}

func (r *nonCompliantRepository) CreateEmployer(employer *models.NonCompliantEmployer) error {
	query := `
		INSERT INTO non_compliant_employers (
			id, business_operating_name, business_legal_name, address,
			date_of_final_decision, penalty_amount, penalty_currency, status, scraped_at
		) VALUES (
			COALESCE(NULLIF($1, ''), gen_random_uuid()), $2, $3, $4, $5, $6, $7, $8, $9
		)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		employer.ID, employer.BusinessOperatingName, employer.BusinessLegalName,
		employer.Address, employer.DateOfFinalDecision, employer.PenaltyAmount,
		employer.PenaltyCurrency, employer.Status, employer.ScrapedAt,
	).Scan(&employer.ID, &employer.CreatedAt, &employer.UpdatedAt)
}

func (r *nonCompliantRepository) CreateEmployersBatch(employers []models.NonCompliantEmployer) error {
	if len(employers) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO non_compliant_employers (
			id, business_operating_name, business_legal_name, address,
			date_of_final_decision, penalty_amount, penalty_currency, status, scraped_at
		) VALUES (
			COALESCE(NULLIF($1, ''), gen_random_uuid()), $2, $3, $4, $5, $6, $7, $8, $9
		) ON CONFLICT (business_operating_name, COALESCE(date_of_final_decision, '1900-01-01'::date)) DO UPDATE SET
			business_legal_name = EXCLUDED.business_legal_name,
			address = EXCLUDED.address,
			penalty_amount = EXCLUDED.penalty_amount,
			penalty_currency = EXCLUDED.penalty_currency,
			status = EXCLUDED.status,
			scraped_at = EXCLUDED.scraped_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id`

	for i := range employers {
		err = tx.QueryRow(query,
			employers[i].ID, employers[i].BusinessOperatingName, employers[i].BusinessLegalName,
			employers[i].Address, employers[i].DateOfFinalDecision, employers[i].PenaltyAmount,
			employers[i].PenaltyCurrency, employers[i].Status, employers[i].ScrapedAt,
		).Scan(&employers[i].ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *nonCompliantRepository) GetEmployerByID(id string) (*models.NonCompliantEmployer, error) {
	var employer models.NonCompliantEmployer
	query := `
		SELECT id, business_operating_name, business_legal_name, address,
		       date_of_final_decision, penalty_amount, penalty_currency, status,
		       scraped_at, created_at, updated_at
		FROM non_compliant_employers
		WHERE id = $1`

	err := r.db.Get(&employer, query, id)
	if err != nil {
		return nil, err
	}

	// Load associated reasons
	reasons, err := r.GetEmployerReasons(id)
	if err == nil {
		employer.Reasons = reasons
	}

	return &employer, nil
}

func (r *nonCompliantRepository) GetEmployersWithReasons(limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error) {
	query := `
		SELECT
			e.id, e.business_operating_name, e.business_legal_name, e.address,
			e.date_of_final_decision, e.penalty_amount, e.penalty_currency, e.status,
			e.postal_code, e.scraped_at, e.created_at, e.updated_at,
			COALESCE(
				array_agg(r.reason_code ORDER BY r.reason_code)
				FILTER (WHERE r.reason_code IS NOT NULL),
				ARRAY[]::VARCHAR[]
			) as reason_codes
		FROM non_compliant_employers e
		LEFT JOIN non_compliant_employer_reasons er ON e.id = er.employer_id
		LEFT JOIN non_compliant_reasons r ON er.reason_id = r.id
		GROUP BY e.id, e.business_operating_name, e.business_legal_name, e.address,
		         e.date_of_final_decision, e.penalty_amount, e.penalty_currency, e.status,
		         e.postal_code, e.scraped_at, e.created_at, e.updated_at
		ORDER BY e.date_of_final_decision DESC, e.business_operating_name
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employers []models.NonCompliantEmployerWithReasonCodes
	for rows.Next() {
		var employer models.NonCompliantEmployerWithReasonCodes
		var reasonCodesArray interface{}

		err := rows.Scan(
			&employer.ID, &employer.BusinessOperatingName, &employer.BusinessLegalName,
			&employer.Address, &employer.DateOfFinalDecision, &employer.PenaltyAmount,
			&employer.PenaltyCurrency, &employer.Status, &employer.PostalCode,
			&employer.ScrapedAt, &employer.CreatedAt, &employer.UpdatedAt,
			&reasonCodesArray,
		)
		if err != nil {
			return nil, err
		}

		// Convert PostgreSQL array to string slice
		if reasonCodesArray != nil {
			if codes, ok := reasonCodesArray.([]interface{}); ok {
				employer.ReasonCodes = make([]string, len(codes))
				for i, code := range codes {
					if codeStr, ok := code.(string); ok {
						employer.ReasonCodes[i] = codeStr
					}
				}
			}
		}

		employers = append(employers, employer)
	}

	return employers, rows.Err()
}

func (r *nonCompliantRepository) GetEmployersCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM non_compliant_employers")
	return count, err
}

func (r *nonCompliantRepository) UpdateEmployer(employer *models.NonCompliantEmployer) error {
	query := `
		UPDATE non_compliant_employers SET
			business_operating_name = $2, business_legal_name = $3, address = $4,
			date_of_final_decision = $5, penalty_amount = $6, penalty_currency = $7,
			status = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, employer.ID, employer.BusinessOperatingName,
		employer.BusinessLegalName, employer.Address, employer.DateOfFinalDecision,
		employer.PenaltyAmount, employer.PenaltyCurrency, employer.Status)
	return err
}

func (r *nonCompliantRepository) DeleteEmployer(id string) error {
	_, err := r.db.Exec("DELETE FROM non_compliant_employers WHERE id = $1", id)
	return err
}

func (r *nonCompliantRepository) CreateReason(reason *models.NonCompliantReason) error {
	query := `
		INSERT INTO non_compliant_reasons (reason_code, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, reason.ReasonCode, reason.Description).Scan(
		&reason.ID, &reason.CreatedAt, &reason.UpdatedAt)
}

func (r *nonCompliantRepository) GetReasonByCode(code string) (*models.NonCompliantReason, error) {
	var reason models.NonCompliantReason
	query := "SELECT id, reason_code, description, created_at, updated_at FROM non_compliant_reasons WHERE reason_code = $1"
	err := r.db.Get(&reason, query, code)
	if err != nil {
		return nil, err
	}
	return &reason, nil
}

func (r *nonCompliantRepository) GetAllReasons() ([]models.NonCompliantReason, error) {
	var reasons []models.NonCompliantReason
	query := "SELECT id, reason_code, description, created_at, updated_at FROM non_compliant_reasons ORDER BY reason_code"
	err := r.db.Select(&reasons, query)
	return reasons, err
}

func (r *nonCompliantRepository) UpsertReason(code, description string) (*models.NonCompliantReason, error) {
	var reason models.NonCompliantReason
	query := `
		INSERT INTO non_compliant_reasons (reason_code, description)
		VALUES ($1, $2)
		ON CONFLICT (reason_code) DO UPDATE SET
			description = EXCLUDED.description,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, reason_code, description, created_at, updated_at`

	err := r.db.QueryRow(query, code, description).Scan(
		&reason.ID, &reason.ReasonCode, &reason.Description, &reason.CreatedAt, &reason.UpdatedAt)
	return &reason, err
}

func (r *nonCompliantRepository) AddEmployerReason(employerID string, reasonID int) error {
	query := `
		INSERT INTO non_compliant_employer_reasons (employer_id, reason_id)
		VALUES ($1, $2)
		ON CONFLICT (employer_id, reason_id) DO NOTHING`

	_, err := r.db.Exec(query, employerID, reasonID)
	return err
}

func (r *nonCompliantRepository) RemoveEmployerReason(employerID string, reasonID int) error {
	_, err := r.db.Exec("DELETE FROM non_compliant_employer_reasons WHERE employer_id = $1 AND reason_id = $2", employerID, reasonID)
	return err
}

func (r *nonCompliantRepository) GetEmployerReasons(employerID string) ([]models.NonCompliantReason, error) {
	var reasons []models.NonCompliantReason
	query := `
		SELECT r.id, r.reason_code, r.description, r.created_at, r.updated_at
		FROM non_compliant_reasons r
		JOIN non_compliant_employer_reasons er ON r.id = er.reason_id
		WHERE er.employer_id = $1
		ORDER BY r.reason_code`

	err := r.db.Select(&reasons, query, employerID)
	return reasons, err
}

func (r *nonCompliantRepository) CreateEmployersWithReasons(data []models.ScraperNonCompliantData) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range data {
		// Parse date
		var finalDecisionDate *time.Time
		if item.DateOfFinalDecision != "" {
			if parsedDate, err := time.Parse("2006-01-02", item.DateOfFinalDecision); err == nil {
				finalDecisionDate = &parsedDate
			}
		}

		// Create employer
		employer := models.NonCompliantEmployer{
			BusinessOperatingName: item.BusinessOperatingName,
			BusinessLegalName:     &item.BusinessLegalName,
			Address:               &item.Address,
			DateOfFinalDecision:   finalDecisionDate,
			PenaltyAmount:         &item.PenaltyAmount,
			PenaltyCurrency:       item.PenaltyCurrency,
			Status:                &item.Status,
			ScrapedAt:             time.Now(),
		}

		// For now, just insert without ON CONFLICT since the constraint might not exist yet
		employerQuery := `
			INSERT INTO non_compliant_employers (
				business_operating_name, business_legal_name, address,
				date_of_final_decision, penalty_amount, penalty_currency, status, scraped_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`

		err = tx.QueryRow(employerQuery,
			employer.BusinessOperatingName, employer.BusinessLegalName, employer.Address,
			employer.DateOfFinalDecision, employer.PenaltyAmount, employer.PenaltyCurrency,
			employer.Status, employer.ScrapedAt,
		).Scan(&employer.ID)
		if err != nil {
			return fmt.Errorf("failed to insert employer: %w", err)
		}

		// Handle reason codes
		for _, reasonCode := range item.ReasonCodes {
			if reasonCode == "" {
				continue
			}

			// Upsert reason
			var reasonID int
			reasonQuery := `
				INSERT INTO non_compliant_reasons (reason_code, description)
				VALUES ($1, $2)
				ON CONFLICT (reason_code) DO UPDATE SET
					description = COALESCE(EXCLUDED.description, non_compliant_reasons.description),
					updated_at = CURRENT_TIMESTAMP
				RETURNING id`

			err = tx.QueryRow(reasonQuery, reasonCode, fmt.Sprintf("Reason code %s", reasonCode)).Scan(&reasonID)
			if err != nil {
				return fmt.Errorf("failed to upsert reason code %s: %w", reasonCode, err)
			}

			// Link employer to reason
			_, err = tx.Exec(`
				INSERT INTO non_compliant_employer_reasons (employer_id, reason_id)
				VALUES ($1, $2)
				ON CONFLICT (employer_id, reason_id) DO NOTHING`,
				employer.ID, reasonID)
			if err != nil {
				return fmt.Errorf("failed to link employer to reason: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *nonCompliantRepository) UpsertEmployersWithReasons(data []models.ScraperNonCompliantData) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range data {
		// Parse date
		var finalDecisionDate *time.Time
		if item.DateOfFinalDecision != "" {
			if parsedDate, err := time.Parse("2006-01-02", item.DateOfFinalDecision); err == nil {
				finalDecisionDate = &parsedDate
			}
		}

		// Upsert employer using ON CONFLICT
		employer := models.NonCompliantEmployer{
			BusinessOperatingName: item.BusinessOperatingName,
			BusinessLegalName:     &item.BusinessLegalName,
			Address:               &item.Address,
			DateOfFinalDecision:   finalDecisionDate,
			PenaltyAmount:         &item.PenaltyAmount,
			PenaltyCurrency:       item.PenaltyCurrency,
			Status:                &item.Status,
			ScrapedAt:             time.Now(),
		}

		employerQuery := `
			INSERT INTO non_compliant_employers (
				business_operating_name, business_legal_name, address,
				date_of_final_decision, penalty_amount, penalty_currency, status, scraped_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (business_operating_name, COALESCE(date_of_final_decision, '1900-01-01'::date)) 
			DO UPDATE SET 
				business_legal_name = EXCLUDED.business_legal_name,
				penalty_amount = EXCLUDED.penalty_amount,
				penalty_currency = EXCLUDED.penalty_currency,
				status = EXCLUDED.status,
				scraped_at = EXCLUDED.scraped_at
			RETURNING id`

		err = tx.QueryRow(employerQuery,
			employer.BusinessOperatingName, employer.BusinessLegalName, employer.Address,
			employer.DateOfFinalDecision, employer.PenaltyAmount, employer.PenaltyCurrency,
			employer.Status, employer.ScrapedAt,
		).Scan(&employer.ID)
		if err != nil {
			return fmt.Errorf("failed to upsert employer: %w", err)
		}

		// Clear existing reason associations for this employer
		_, err = tx.Exec("DELETE FROM non_compliant_employer_reasons WHERE employer_id = $1", employer.ID)
		if err != nil {
			return fmt.Errorf("failed to clear existing reasons: %w", err)
		}

		// Handle reason codes
		for _, reasonCode := range item.ReasonCodes {
			if reasonCode == "" {
				continue
			}

			// Get or create reason
			var reason models.NonCompliantReason
			reasonQuery := `
				INSERT INTO non_compliant_reasons (reason_code) 
				VALUES ($1) 
				ON CONFLICT (reason_code) DO NOTHING 
				RETURNING id, reason_code`
			err = tx.QueryRow(reasonQuery, reasonCode).Scan(&reason.ID, &reason.ReasonCode)
			if err != nil {
				// If no rows returned, reason already exists, so fetch it
				if err == sql.ErrNoRows {
					err = tx.QueryRow("SELECT id, reason_code FROM non_compliant_reasons WHERE reason_code = $1", reasonCode).
						Scan(&reason.ID, &reason.ReasonCode)
					if err != nil {
						return fmt.Errorf("failed to fetch existing reason: %w", err)
					}
				} else {
					return fmt.Errorf("failed to upsert reason: %w", err)
				}
			}

			// Create employer-reason association
			_, err = tx.Exec("INSERT INTO non_compliant_employer_reasons (employer_id, reason_id) VALUES ($1, $2)", 
				employer.ID, reason.ID)
			if err != nil {
				return fmt.Errorf("failed to create employer-reason association: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (r *nonCompliantRepository) ClearAllNonCompliantData() error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete in correct order due to foreign key constraints
	_, err = tx.Exec("DELETE FROM non_compliant_employer_reasons")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM non_compliant_employers")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *nonCompliantRepository) GetLatestScrapedDate() (*time.Time, error) {
	var latestDate sql.NullTime
	err := r.db.Get(&latestDate, "SELECT MAX(scraped_at) FROM non_compliant_employers")
	if err != nil {
		return nil, err
	}

	if latestDate.Valid {
		return &latestDate.Time, nil
	}
	return nil, nil
}

func (r *nonCompliantRepository) GetTotalEmployersCount() (int, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM non_compliant_employers")
	return count, err
}

// GetEmployersWithoutPostalCodes returns all employers that don't have postal codes yet
func (r *nonCompliantRepository) GetEmployersWithoutPostalCodes() ([]models.NonCompliantEmployer, error) {
	var employers []models.NonCompliantEmployer
	query := `
		SELECT id, business_operating_name, business_legal_name, address,
		       date_of_final_decision, penalty_amount, penalty_currency, status,
		       postal_code, scraped_at, created_at, updated_at
		FROM non_compliant_employers
		WHERE postal_code IS NULL AND address IS NOT NULL AND address != ''`

	err := r.db.Select(&employers, query)
	return employers, err
}

// UpdateEmployerPostalCode updates the postal code for an employer
func (r *nonCompliantRepository) UpdateEmployerPostalCode(employerID, postalCode string) error {
	query := `
		UPDATE non_compliant_employers 
		SET postal_code = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, employerID, postalCode)
	return err
}

// GetLocationsByPostalCode returns aggregated location data grouped by postal code
func (r *nonCompliantRepository) GetLocationsByPostalCode(limit int) ([]models.NonCompliantPostalCodeLocation, error) {
	var locations []models.NonCompliantPostalCodeLocation
	query := `
		WITH extracted_postal_codes AS (
			SELECT 
				e.*,
				-- Extract Canadian postal codes from address using regex
				UPPER(SUBSTRING(e.address FROM '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d')) as extracted_postal_code
			FROM non_compliant_employers e
			WHERE e.address IS NOT NULL AND e.address != ''
		)
		SELECT 
			epc.extracted_postal_code as postal_code,
			pc.latitude,
			pc.longitude,
			COUNT(*) as employer_count,
			COALESCE(SUM(epc.penalty_amount), 0) as total_penalty_amount,
			COUNT(*) as violation_count,
			MAX(epc.date_of_final_decision) as most_recent_violation
		FROM extracted_postal_codes epc
		JOIN postal_codes pc ON epc.extracted_postal_code = pc.postal_code
		WHERE epc.extracted_postal_code IS NOT NULL 
		  AND pc.latitude IS NOT NULL 
		  AND pc.longitude IS NOT NULL
		GROUP BY epc.extracted_postal_code, pc.latitude, pc.longitude
		ORDER BY employer_count DESC, total_penalty_amount DESC
		LIMIT $1`

	err := r.db.Select(&locations, query, limit)
	return locations, err
}

// GetEmployersByPostalCode returns all employers for a specific postal code
func (r *nonCompliantRepository) GetEmployersByPostalCode(postalCode string, limit, offset int) ([]models.NonCompliantEmployerWithReasonCodes, error) {
	query := `
		WITH extracted_postal_codes AS (
			SELECT 
				e.*,
				-- Extract Canadian postal codes from address using regex
				UPPER(SUBSTRING(e.address FROM '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d')) as extracted_postal_code
			FROM non_compliant_employers e
			WHERE e.address IS NOT NULL AND e.address != ''
		)
		SELECT
			epc.id, epc.business_operating_name, epc.business_legal_name, epc.address,
			epc.date_of_final_decision, epc.penalty_amount, epc.penalty_currency, epc.status,
			epc.postal_code, epc.scraped_at, epc.created_at, epc.updated_at,
			COALESCE(
				array_agg(r.reason_code ORDER BY r.reason_code)
				FILTER (WHERE r.reason_code IS NOT NULL),
				ARRAY[]::VARCHAR[]
			) as reason_codes
		FROM extracted_postal_codes epc
		LEFT JOIN non_compliant_employer_reasons er ON epc.id = er.employer_id
		LEFT JOIN non_compliant_reasons r ON er.reason_id = r.id
		WHERE epc.extracted_postal_code = $1
		GROUP BY epc.id, epc.business_operating_name, epc.business_legal_name, epc.address,
		         epc.date_of_final_decision, epc.penalty_amount, epc.penalty_currency, epc.status,
		         epc.postal_code, epc.scraped_at, epc.created_at, epc.updated_at
		ORDER BY epc.date_of_final_decision DESC, epc.business_operating_name
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, postalCode, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employers []models.NonCompliantEmployerWithReasonCodes
	for rows.Next() {
		var employer models.NonCompliantEmployerWithReasonCodes
		var reasonCodesArray interface{}

		err := rows.Scan(
			&employer.ID, &employer.BusinessOperatingName, &employer.BusinessLegalName,
			&employer.Address, &employer.DateOfFinalDecision, &employer.PenaltyAmount,
			&employer.PenaltyCurrency, &employer.Status, &employer.PostalCode,
			&employer.ScrapedAt, &employer.CreatedAt, &employer.UpdatedAt,
			&reasonCodesArray,
		)
		if err != nil {
			return nil, err
		}

		// Convert PostgreSQL array to string slice
		if reasonCodesArray != nil {
			if codes, ok := reasonCodesArray.([]interface{}); ok {
				employer.ReasonCodes = make([]string, len(codes))
				for i, code := range codes {
					if codeStr, ok := code.(string); ok {
						employer.ReasonCodes[i] = codeStr
					}
				}
			}
		}

		employers = append(employers, employer)
	}

	return employers, rows.Err()
}
