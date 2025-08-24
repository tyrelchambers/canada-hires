package repos

import (
	"canada-hires/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type NonCompliantRepository interface {
	// Employers
	CreateEmployer(employer *models.NonCompliantEmployer) error
	CreateEmployersBatch(employers []models.NonCompliantEmployer) error
	GetEmployerByID(id string) (*models.NonCompliantEmployer, error)
	GetEmployersWithReasons(limit, offset int) ([]models.NonCompliantEmployerWithReasons, error)
	GetEmployersCount() (int, error)
	UpdateEmployer(employer *models.NonCompliantEmployer) error
	DeleteEmployer(id string) error

	// Reasons
	CreateReason(reason *models.NonCompliantReason) error
	GetReasonByCode(code string) (*models.NonCompliantReason, error)
	GetAllReasons() ([]models.NonCompliantReason, error)
	UpsertReason(code, description string) (*models.NonCompliantReason, error)

	// Bulk operations
	CreateEmployersWithReasons(data []models.ScraperNonCompliantData) error
	UpsertEmployersWithReasons(data []models.ScraperNonCompliantData) error
	ClearAllNonCompliantData() error

	// Stats
	GetLatestScrapedDate() (*time.Time, error)
	GetTotalEmployersCount() (int, error)

	// Geolocation methods
	GetEmployersWithoutPostalCodes() ([]models.NonCompliantEmployer, error)
	GetEmployersWithoutExtractablePostalCodes() ([]models.NonCompliantEmployer, error)
	UpdateEmployerPostalCode(employerID, postalCode string) error
	UpdateEmployerAddress(employerID, address string) error
	GetLocationsByPostalCode(limit int) ([]models.NonCompliantPostalCodeLocation, error)
	GetEmployersByPostalCode(postalCode string, limit, offset int) ([]models.NonCompliantEmployerWithReasons, error)
	GetEmployersByCoordinates(lat, lng float64, limit, offset int) ([]models.NonCompliantEmployerWithReasons, error)
	GetAllEmployers() ([]models.NonCompliantEmployer, error)
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
			date_of_final_decision, penalty_amount, penalty_currency, status, reason_codes, scraped_at
		) VALUES (
			COALESCE(NULLIF($1, ''), gen_random_uuid()), $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		employer.ID, employer.BusinessOperatingName, employer.BusinessLegalName,
		employer.Address, employer.DateOfFinalDecision, employer.PenaltyAmount,
		employer.PenaltyCurrency, employer.Status, employer.ReasonCodes, employer.ScrapedAt,
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
		       reason_codes, postal_code, scraped_at, created_at, updated_at
		FROM non_compliant_employers
		WHERE id = $1`

	err := r.db.Get(&employer, query, id)
	if err != nil {
		return nil, err
	}

	// Load associated reasons from reason codes
	if len(employer.ReasonCodes) > 0 {
		var reasons []models.NonCompliantReason
		for _, code := range employer.ReasonCodes {
			if reason, err := r.GetReasonByCode(code); err == nil {
				reasons = append(reasons, *reason)
			}
		}
		employer.Reasons = reasons
	}

	return &employer, nil
}

func (r *nonCompliantRepository) GetEmployersWithReasons(limit, offset int) ([]models.NonCompliantEmployerWithReasons, error) {
	query := `
		SELECT
			e.id, e.business_operating_name, e.business_legal_name, e.address,
			e.date_of_final_decision, e.penalty_amount, e.penalty_currency, e.status,
			e.reason_codes, e.postal_code, e.scraped_at, e.created_at, e.updated_at
		FROM non_compliant_employers e
		ORDER BY e.date_of_final_decision DESC, e.business_operating_name
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employers []models.NonCompliantEmployerWithReasons
	for rows.Next() {
		var employer models.NonCompliantEmployerWithReasons
		var reasonCodesArray interface{}

		err := rows.Scan(
			&employer.ID, &employer.BusinessOperatingName, &employer.BusinessLegalName,
			&employer.Address, &employer.DateOfFinalDecision, &employer.PenaltyAmount,
			&employer.PenaltyCurrency, &employer.Status, &reasonCodesArray, &employer.PostalCode,
			&employer.ScrapedAt, &employer.CreatedAt, &employer.UpdatedAt,
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

		// Get full reason objects with descriptions for the reason codes
		if len(employer.ReasonCodes) > 0 {
			var reasons []models.NonCompliantReason
			for _, code := range employer.ReasonCodes {
				if reason, err := r.GetReasonByCode(code); err == nil {
					reasons = append(reasons, *reason)
				}
			}
			employer.Reasons = reasons
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

		// Convert reason codes slice to PostgreSQL array format
		reasonCodes := item.ReasonCodes

		employerQuery := `
			INSERT INTO non_compliant_employers (
				business_operating_name, business_legal_name, address,
				date_of_final_decision, penalty_amount, penalty_currency, status, reason_codes, scraped_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			RETURNING id`

		var employerID string
		err = tx.QueryRow(employerQuery,
			item.BusinessOperatingName, &item.BusinessLegalName, &item.Address,
			finalDecisionDate, &item.PenaltyAmount, item.PenaltyCurrency,
			&item.Status, pq.Array(reasonCodes), time.Now(),
		).Scan(&employerID)
		if err != nil {
			return fmt.Errorf("failed to insert employer: %w", err)
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

		// Convert reason codes slice to PostgreSQL array format
		reasonCodes := item.ReasonCodes

		employerQuery := `
			INSERT INTO non_compliant_employers (
				business_operating_name, business_legal_name, address,
				date_of_final_decision, penalty_amount, penalty_currency, status, reason_codes, scraped_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (business_operating_name, COALESCE(date_of_final_decision, '1900-01-01'::date))
			DO UPDATE SET
				business_legal_name = EXCLUDED.business_legal_name,
				penalty_amount = EXCLUDED.penalty_amount,
				penalty_currency = EXCLUDED.penalty_currency,
				status = EXCLUDED.status,
				reason_codes = EXCLUDED.reason_codes,
				scraped_at = EXCLUDED.scraped_at
			RETURNING id`

		var employerID string
		err = tx.QueryRow(employerQuery,
			item.BusinessOperatingName, &item.BusinessLegalName, &item.Address,
			finalDecisionDate, &item.PenaltyAmount, item.PenaltyCurrency,
			&item.Status, pq.Array(reasonCodes), time.Now(),
		).Scan(&employerID)
		if err != nil {
			return fmt.Errorf("failed to upsert employer: %w", err)
		}
	}

	return tx.Commit()
}

func (r *nonCompliantRepository) ClearAllNonCompliantData() error {
	_, err := r.db.Exec("DELETE FROM non_compliant_employers")
	return err
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

// GetLocationsByPostalCode returns aggregated location data grouped by postal code and includes direct coordinates
func (r *nonCompliantRepository) GetLocationsByPostalCode(limit int) ([]models.NonCompliantPostalCodeLocation, error) {
	var locations []models.NonCompliantPostalCodeLocation
	query := `
		WITH postal_code_locations AS (
			-- Get locations based on postal codes
			SELECT
				epc.extracted_postal_code as postal_code,
				pc.latitude,
				pc.longitude,
				COUNT(*) as employer_count,
				COALESCE(SUM(epc.penalty_amount), 0) as total_penalty_amount,
				COUNT(*) as violation_count,
				MAX(epc.date_of_final_decision) as most_recent_violation
			FROM (
				SELECT
					e.*,
					-- Extract Canadian postal codes from address using regex
					UPPER(SUBSTRING(e.address FROM '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d')) as extracted_postal_code
				FROM non_compliant_employers e
				WHERE e.address IS NOT NULL AND e.address != ''
			) epc
			JOIN postal_codes pc ON epc.extracted_postal_code = pc.postal_code
			WHERE epc.extracted_postal_code IS NOT NULL
			  AND pc.latitude IS NOT NULL
			  AND pc.longitude IS NOT NULL
			GROUP BY epc.extracted_postal_code, pc.latitude, pc.longitude
		),
		address_locations AS (
			-- Get locations based on address geocoding cache (only for employers without extractable postal codes)
			SELECT
				'COORD_' || ROUND(agc.latitude::numeric, 6) || '_' || ROUND(agc.longitude::numeric, 6) as postal_code,
				agc.latitude,
				agc.longitude,
				COUNT(*) as employer_count,
				COALESCE(SUM(e.penalty_amount), 0) as total_penalty_amount,
				COUNT(*) as violation_count,
				MAX(e.date_of_final_decision) as most_recent_violation
			FROM non_compliant_employers e
			JOIN address_geocoding_cache agc ON agc.normalized_address = LOWER(TRIM(REGEXP_REPLACE(REGEXP_REPLACE(e.address, '[,.]', '', 'g'), '\s+', ' ', 'g')))
			WHERE e.address IS NOT NULL
			  AND e.address != ''
			  AND UPPER(SUBSTRING(e.address FROM '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d')) IS NULL
			GROUP BY agc.latitude, agc.longitude
		)
		-- Union both location types
		SELECT * FROM postal_code_locations
		UNION ALL
		SELECT * FROM address_locations
		ORDER BY employer_count DESC, total_penalty_amount DESC
		LIMIT $1`

	err := r.db.Select(&locations, query, limit)
	return locations, err
}

// GetEmployersByPostalCode returns all employers for a specific postal code
func (r *nonCompliantRepository) GetEmployersByPostalCode(postalCode string, limit, offset int) ([]models.NonCompliantEmployerWithReasons, error) {
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
			epc.reason_codes, epc.postal_code, epc.scraped_at, epc.created_at, epc.updated_at
		FROM extracted_postal_codes epc
		WHERE epc.extracted_postal_code = $1
		ORDER BY epc.date_of_final_decision DESC, epc.business_operating_name
		LIMIT $2 OFFSET $3`

	var employers []models.NonCompliantEmployerWithReasons
	err := r.db.Select(&employers, query, postalCode, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range employers {
		// Get full reason objects with descriptions for the reason codes
		if len(employers[i].ReasonCodes) > 0 {
			var reasons []models.NonCompliantReason
			for _, code := range employers[i].ReasonCodes {
				if reason, err := r.GetReasonByCode(code); err == nil {
					reasons = append(reasons, *reason)
				}
			}
			employers[i].Reasons = reasons
		}
	}

	return employers, err
}

// GetEmployersWithoutExtractablePostalCodes returns employers that don't have extractable
// postal codes for address geocoding
func (r *nonCompliantRepository) GetEmployersWithoutExtractablePostalCodes() ([]models.NonCompliantEmployer, error) {
	var employers []models.NonCompliantEmployer
	query := `
  		SELECT id, business_operating_name, business_legal_name, address,
  		       date_of_final_decision, penalty_amount, penalty_currency, status,
  		       postal_code, scraped_at, created_at, updated_at
  		FROM non_compliant_employers
  		WHERE address IS NOT NULL
  		  AND address != ''
  		  AND UPPER(SUBSTRING(address FROM '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d')) IS NULL`

	err := r.db.Select(&employers, query)
	return employers, err
}

// GetEmployersByCoordinates returns all employers at specific lat/lng coordinates
func (r *nonCompliantRepository) GetEmployersByCoordinates(lat, lng float64, limit, offset int) ([]models.NonCompliantEmployerWithReasons, error) {
	query := `
		SELECT
			e.*
		FROM non_compliant_employers e
		JOIN address_geocoding_cache agc ON agc.normalized_address = LOWER(TRIM(REGEXP_REPLACE(REGEXP_REPLACE(e.address, '[,.]', '', 'g'), '\s+', ' ', 'g')))
		WHERE e.address IS NOT NULL
		  AND e.address != ''
		  AND ABS(agc.latitude - $1) < 0.001
		  AND ABS(agc.longitude - $2) < 0.001
		ORDER BY e.date_of_final_decision DESC, e.business_operating_name
		LIMIT $3 OFFSET $4`

	var employers []models.NonCompliantEmployerWithReasons

	err := r.db.Select(&employers, query, lat, lng, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range employers {
		for _, rc := range employers[i].ReasonCodes {
			c, err := r.GetReasonByCode(rc)
			if err != nil {
				return nil, err
			}

			employers[i].Reasons = append(employers[i].Reasons, *c)
		}
	}

	return employers, err
}

// UpdateEmployerAddress updates the address for an employer
func (r *nonCompliantRepository) UpdateEmployerAddress(employerID, address string) error {
	query := `
		UPDATE non_compliant_employers
		SET address = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`

	_, err := r.db.Exec(query, employerID, address)
	return err
}

// GetAllEmployers returns all employers for address cleaning
func (r *nonCompliantRepository) GetAllEmployers() ([]models.NonCompliantEmployer, error) {
	var employers []models.NonCompliantEmployer
	query := `
		SELECT id, business_operating_name, business_legal_name, address,
		       date_of_final_decision, penalty_amount, penalty_currency, status,
		       postal_code, scraped_at, created_at, updated_at
		FROM non_compliant_employers
		ORDER BY created_at DESC`

	err := r.db.Select(&employers, query)
	return employers, err
}
