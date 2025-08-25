package repos

import (
	"canada-hires/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LMIARepository interface {
	// LMIA Resources
	CreateResource(resource *models.LMIAResource) error
	GetResourceByResourceID(resourceID string) (*models.LMIAResource, error)
	GetResourcesByLanguage(language string) ([]*models.LMIAResource, error)
	UpdateResourceDownloaded(id string) error
	UpdateResourceProcessed(id string) error
	GetUnprocessedResources() ([]*models.LMIAResource, error)

	// LMIA Employers
	CreateEmployer(employer *models.LMIAEmployer) error
	CreateEmployersBatch(employers []*models.LMIAEmployer) error
	GetEmployersByResourceID(resourceID string) ([]*models.LMIAEmployer, error)
	SearchEmployersByName(name string, limit int) ([]*models.LMIAEmployer, error)
	SearchEmployersByNameAndPeriod(name string, year int, quarter string, limit int) ([]*models.LMIAEmployer, error)
	GetEmployersByLocation(city, province string, limit int) ([]*models.LMIAEmployer, error)
	GetEmployersByYear(year int, limit int) ([]*models.LMIAEmployer, error)
	GetEmployersByYearAndQuarter(year int, quarter string, limit int) ([]*models.LMIAEmployer, error)
	GetEmployersWithGeolocation(year int, quarter string, limit int) ([]*models.LMIAEmployerGeoLocation, error)
	AllEmployersCount() (int, error)
	GetYearRange() (minYear, maxYear int, err error)
	GetDistinctEmployersCount() (int, error)
	GetGeographicSummary(year int) ([]*models.LMIAGeographicSummary, error)

	// Postal Code Methods
	GetPostalCodeLocations(year int, quarter string, limit int) ([]*models.PostalCodeLocation, error)
	GetEmployersByPostalCode(postalCode string, year int, quarter string, limit int) ([]*models.LMIAEmployer, error)
	GetEmployersNeedingPostalCodeExtraction(limit int) ([]*models.LMIAEmployer, error)
	UpdateEmployerPostalCode(id string, postalCode string) error
	GetUngeocodedPostalCodes() (map[string]string, error)

	// Cron Jobs
	CreateCronJob(job *models.CronJob) error
	UpdateCronJobStatus(id string, status string, errorMessage *string) error
	UpdateCronJobCompleted(id string, resourcesProcessed, recordsProcessed int) error
	GetLatestCronJob(jobName string) (*models.CronJob, error)
}

type lmiaRepository struct {
	db *sqlx.DB
}

func NewLMIARepository(db *sqlx.DB) LMIARepository {
	return &lmiaRepository{db: db}
}

// LMIA Resources methods
func (r *lmiaRepository) CreateResource(resource *models.LMIAResource) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lmia_resources (id, resource_id, name, quarter, year, url, format, language, size_bytes,
								   last_modified, date_published, created_at, updated_at)
		VALUES (:id, :resource_id, :name, :quarter, :year, :url, :format, :language, :size_bytes,
				:last_modified, :date_published, :created_at, :updated_at)
	`

	resource.ID = uuid.New().String()
	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	_, err = tx.NamedExec(query, resource)
	if err != nil {
		return fmt.Errorf("failed to insert LMIA resource: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *lmiaRepository) GetResourceByResourceID(resourceID string) (*models.LMIAResource, error) {
	var resource models.LMIAResource
	query := `SELECT * FROM lmia_resources WHERE resource_id = $1`

	err := r.db.Get(&resource, query, resourceID)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *lmiaRepository) GetResourcesByLanguage(language string) ([]*models.LMIAResource, error) {
	var resources []*models.LMIAResource
	query := `SELECT * FROM lmia_resources WHERE language = $1 ORDER BY year DESC, quarter DESC`

	err := r.db.Select(&resources, query, language)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (r *lmiaRepository) UpdateResourceDownloaded(id string) error {
	query := `UPDATE lmia_resources SET downloaded_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *lmiaRepository) UpdateResourceProcessed(id string) error {
	query := `UPDATE lmia_resources SET processed_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *lmiaRepository) GetUnprocessedResources() ([]*models.LMIAResource, error) {
	var resources []*models.LMIAResource
	query := `
		SELECT DISTINCT r.* 
		FROM lmia_resources r 
		LEFT JOIN (
			SELECT DISTINCT resource_id 
			FROM lmia_employers
		) e ON r.id = e.resource_id 
		WHERE e.resource_id IS NULL 
		ORDER BY r.year DESC, r.quarter DESC`

	err := r.db.Select(&resources, query)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

// LMIA Employers methods
func (r *lmiaRepository) CreateEmployer(employer *models.LMIAEmployer) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
		INSERT INTO lmia_employers (id, resource_id, province_territory, program_stream, employer,
								   address, occupation, incorporate_status, approved_lmias, approved_positions,
								   quarter, year, created_at, updated_at, postal_code)
		VALUES (:id, :resource_id, :province_territory, :program_stream, :employer,
				:address, :occupation, :incorporate_status, :approved_lmias, :approved_positions,
				:quarter, :year, :created_at, :updated_at, :postal_code)
	`

	employer.ID = uuid.New().String()
	employer.CreatedAt = time.Now()
	employer.UpdatedAt = time.Now()

	_, err = tx.NamedExec(query, employer)
	if err != nil {
		return fmt.Errorf("failed to insert LMIA employer: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *lmiaRepository) CreateEmployersBatch(employers []*models.LMIAEmployer) error {
	if len(employers) == 0 {
		return nil
	}

	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO lmia_employers (id, resource_id, province_territory, program_stream, employer,
								   address, occupation, incorporate_status, approved_lmias, approved_positions,
								   quarter, year, created_at, updated_at, postal_code)
		VALUES (:id, :resource_id, :province_territory, :program_stream, :employer,
				:address, :occupation, :incorporate_status, :approved_lmias, :approved_positions,
				:quarter, :year, :created_at, :updated_at, :postal_code)
	`

	for _, employer := range employers {
		employer.ID = uuid.New().String()
		employer.CreatedAt = time.Now()
		employer.UpdatedAt = time.Now()
	}

	_, err = tx.NamedExec(query, employers)
	if err != nil {
		return fmt.Errorf("failed to insert LMIA employers batch: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *lmiaRepository) GetEmployersByResourceID(resourceID string) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer
	query := `
		SELECT e.*, r.quarter, r.year
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE e.resource_id = $1
		ORDER BY e.employer
	`

	err := r.db.Select(&employers, query, resourceID)
	if err != nil {
		return nil, err
	}

	return employers, nil
}

func (r *lmiaRepository) SearchEmployersByName(name string, limit int) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer
	query := `
		SELECT e.*, r.quarter, r.year
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE e.employer ILIKE $1
		ORDER BY r.year DESC, r.quarter DESC, e.employer
	`

	searchTerm := "%" + name + "%"

	// If limit is 0 or negative, return all records
	if limit > 0 {
		query += " LIMIT $2"
		err := r.db.Select(&employers, query, searchTerm, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&employers, query, searchTerm)
		if err != nil {
			return nil, err
		}
	}

	return employers, nil
}

func (r *lmiaRepository) GetEmployersByLocation(city, province string, limit int) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer
	query := `
		SELECT e.*, r.quarter, r.year
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE ($1 = '' OR e.address ILIKE $1) AND ($2 = '' OR e.province_territory ILIKE $2)
		ORDER BY r.year DESC, r.quarter DESC, e.employer
	`

	citySearch := ""
	provinceSearch := ""
	if city != "" {
		citySearch = "%" + city + "%"
	}
	if province != "" {
		provinceSearch = "%" + province + "%"
	}

	// If limit is 0 or negative, return all records
	if limit > 0 {
		query += " LIMIT $3"
		err := r.db.Select(&employers, query, citySearch, provinceSearch, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&employers, query, citySearch, provinceSearch)
		if err != nil {
			return nil, err
		}
	}

	return employers, nil
}

// Cron Jobs methods
func (r *lmiaRepository) CreateCronJob(job *models.CronJob) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO cron_jobs (id, job_name, status, started_at, created_at)
		VALUES (:id, :job_name, :status, :started_at, :created_at)
	`

	job.ID = uuid.New().String()
	job.CreatedAt = time.Now()

	_, err = tx.NamedExec(query, job)
	if err != nil {
		return fmt.Errorf("failed to insert cron job: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *lmiaRepository) UpdateCronJobStatus(id string, status string, errorMessage *string) error {
	query := `UPDATE cron_jobs SET status = $2, error_message = $3 WHERE id = $1`
	_, err := r.db.Exec(query, id, status, errorMessage)
	return err
}

func (r *lmiaRepository) UpdateCronJobCompleted(id string, resourcesProcessed, recordsProcessed int) error {
	query := `
		UPDATE cron_jobs
		SET status = 'completed', completed_at = NOW(), resources_processed = $2, records_processed = $3
		WHERE id = $1
	`
	_, err := r.db.Exec(query, id, resourcesProcessed, recordsProcessed)
	return err
}

func (r *lmiaRepository) GetLatestCronJob(jobName string) (*models.CronJob, error) {
	var job models.CronJob
	query := `SELECT * FROM cron_jobs WHERE job_name = $1 ORDER BY started_at DESC LIMIT 1`

	err := r.db.Get(&job, query, jobName)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (r *lmiaRepository) GetEmployersByYear(year int, limit int) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer

	query := `SELECT e.*, r.quarter, r.year
			FROM lmia_employers e
			JOIN lmia_resources r ON e.resource_id = r.id
			WHERE r.year = $1
			ORDER BY r.year DESC, r.quarter DESC, e.employer`

	// If limit is 0 or negative, return all records
	if limit > 0 {
		query += " LIMIT $2"
		err := r.db.Select(&employers, query, year, limit)
		if err != nil {
			return nil, err
		}
	} else {
		err := r.db.Select(&employers, query, year)
		if err != nil {
			return nil, err
		}
	}

	return employers, nil
}

func (r *lmiaRepository) AllEmployersCount() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM lmia_employers`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *lmiaRepository) GetYearRange() (minYear, maxYear int, err error) {
	query := `SELECT MIN(year) as min_year, MAX(year) as max_year FROM lmia_resources WHERE processed_at IS NOT NULL`

	var result struct {
		MinYear *int `db:"min_year"`
		MaxYear *int `db:"max_year"`
	}

	err = r.db.Get(&result, query)
	if err != nil {
		return 0, 0, err
	}

	if result.MinYear == nil || result.MaxYear == nil {
		return 0, 0, nil
	}

	return *result.MinYear, *result.MaxYear, nil
}

func (r *lmiaRepository) GetDistinctEmployersCount() (int, error) {
	var count int
	query := `SELECT COUNT(DISTINCT employer) FROM lmia_employers`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *lmiaRepository) GetGeographicSummary(year int) ([]*models.LMIAGeographicSummary, error) {
	query := `
		SELECT
			COALESCE(province_territory, 'Unknown') as province_territory,
			COUNT(DISTINCT employer) as total_employers,
			SUM(COALESCE(approved_lmias, 0)) as total_lmias,
			SUM(COALESCE(approved_positions, 0)) as total_positions,
			$2 as year
		FROM lmia_employers
		WHERE year = $1
		  AND province_territory IS NOT NULL
		  AND province_territory != ''
		  AND province_territory != 'N/A'
		GROUP BY province_territory
		ORDER BY total_positions DESC
	`

	var summaries []*models.LMIAGeographicSummary
	err := r.db.Select(&summaries, query, year, year)
	if err != nil {
		return nil, fmt.Errorf("failed to get geographic summary: %w", err)
	}

	return summaries, nil
}

func (r *lmiaRepository) SearchEmployersByNameAndPeriod(name string, year int, quarter string, limit int) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer

	query := `
		SELECT e.*, r.quarter, r.year
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE e.employer ILIKE $1 AND r.year = $2
	`
	args := []interface{}{"%" + name + "%", year}
	argIndex := 3

	if quarter != "" {
		query += " AND r.quarter = $" + strconv.Itoa(argIndex)
		args = append(args, quarter)
		argIndex++
	}

	query += " ORDER BY r.year DESC, r.quarter DESC, e.employer"

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, limit)
	}

	err := r.db.Select(&employers, query, args...)
	if err != nil {
		return nil, err
	}

	return employers, nil
}

func (r *lmiaRepository) GetEmployersByYearAndQuarter(year int, quarter string, limit int) ([]*models.LMIAEmployer, error) {
	var employers []*models.LMIAEmployer

	query := `
		SELECT e.*, r.quarter, r.year
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE r.year = $1
	`
	args := []interface{}{year}
	argIndex := 2

	if quarter != "" {
		query += " AND r.quarter = $" + strconv.Itoa(argIndex)
		args = append(args, quarter)
		argIndex++
	}

	query += " ORDER BY r.year DESC, r.quarter DESC, e.employer"

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, limit)
	}

	err := r.db.Select(&employers, query, args...)
	if err != nil {
		return nil, err
	}

	return employers, nil
}

func (r *lmiaRepository) GetEmployersWithGeolocation(year int, quarter string, limit int) ([]*models.LMIAEmployerGeoLocation, error) {
	var employers []*models.LMIAEmployerGeoLocation

	// Use postal code coordinates from postal_codes table
	query := `
		SELECT
			e.id,
			e.employer,
			e.address,
			e.province_territory,
			e.approved_lmias,
			e.approved_positions,
			r.quarter,
			r.year,
			COALESCE(pc.latitude, 0) as latitude,
			COALESCE(pc.longitude, 0) as longitude,
			SUM(COALESCE(e.approved_lmias, 0)) OVER (PARTITION BY e.employer) as total_lmias
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		LEFT JOIN postal_codes pc ON e.postal_code = pc.postal_code
		WHERE r.year = $1
		AND e.postal_code IS NOT NULL
		AND e.postal_code != ''
		AND pc.latitude IS NOT NULL
		AND pc.longitude IS NOT NULL
	`
	args := []interface{}{year}
	argIndex := 2

	if quarter != "" {
		query += " AND r.quarter = $" + strconv.Itoa(argIndex)
		args = append(args, quarter)
		argIndex++
	}

	query += " ORDER BY e.approved_lmias DESC"

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, limit)
	}

	err := r.db.Select(&employers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get employers with geolocation: %w", err)
	}

	return employers, nil
}

// GetPostalCodeLocations returns employers grouped by postal code with coordinates
func (r *lmiaRepository) GetPostalCodeLocations(year int, quarter string, limit int) ([]*models.PostalCodeLocation, error) {
	// Join with postal_codes table to get coordinates
	query := `
		SELECT
			e.postal_code,
			COALESCE(pc.latitude, 0) as latitude,
			COALESCE(pc.longitude, 0) as longitude,
			JSON_AGG(
				JSON_BUILD_OBJECT(
					'employer', e.employer,
					'occupation', COALESCE(e.occupation, 'Not specified'),
					'approved_lmias', COALESCE(e.approved_lmias, 0),
					'approved_positions', COALESCE(e.approved_positions, 0)
				)
			) as businesses,
			SUM(COALESCE(e.approved_lmias, 0)) as total_lmias,
			COUNT(*) as business_count
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		LEFT JOIN postal_codes pc ON e.postal_code = pc.postal_code
		WHERE r.year = $1
		AND e.postal_code IS NOT NULL
		AND e.postal_code != ''
		AND pc.latitude IS NOT NULL
		AND pc.longitude IS NOT NULL
	`
	args := []interface{}{year}
	argIndex := 2

	if quarter != "" {
		query += " AND r.quarter = $" + strconv.Itoa(argIndex)
		args = append(args, quarter)
		argIndex++
	}

	query += `
		GROUP BY e.postal_code, pc.latitude, pc.longitude
		ORDER BY total_lmias DESC
	`

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get postal code locations: %w", err)
	}
	defer rows.Close()

	var locations []*models.PostalCodeLocation
	for rows.Next() {
		var location models.PostalCodeLocation
		var businessesJSON string

		err := rows.Scan(
			&location.PostalCode,
			&location.Latitude,
			&location.Longitude,
			&businessesJSON,
			&location.TotalLMIAs,
			&location.BusinessCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan postal code location: %w", err)
		}

		// Parse the JSON businesses data
		if err := json.Unmarshal([]byte(businessesJSON), &location.Businesses); err != nil {
			return nil, fmt.Errorf("failed to parse businesses JSON: %w", err)
		}

		locations = append(locations, &location)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating postal code locations: %w", err)
	}

	return locations, nil
}

// GetEmployersByPostalCode returns all employers for a specific postal code
func (r *lmiaRepository) GetEmployersByPostalCode(postalCode string, year int, quarter string, limit int) ([]*models.LMIAEmployer, error) {
	query := `
		SELECT e.* 
		FROM lmia_employers e
		JOIN lmia_resources r ON e.resource_id = r.id
		WHERE e.postal_code = $1
		AND r.year = $2
	`
	args := []interface{}{postalCode, year}
	argIndex := 3

	if quarter != "" {
		query += " AND r.quarter = $" + strconv.Itoa(argIndex)
		args = append(args, quarter)
		argIndex++
	}

	query += " ORDER BY e.approved_lmias DESC"

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(argIndex)
		args = append(args, limit)
	}

	var employers []*models.LMIAEmployer
	err := r.db.Select(&employers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get employers by postal code: %w", err)
	}

	return employers, nil
}

// GetEmployersNeedingPostalCodeExtraction returns employers that need postal code extraction
func (r *lmiaRepository) GetEmployersNeedingPostalCodeExtraction(limit int) ([]*models.LMIAEmployer, error) {
	query := `
		SELECT id, employer, address, province_territory
		FROM lmia_employers
		WHERE address IS NOT NULL
		AND address != ''
		AND postal_code IS NULL
		AND address ~ '[A-Za-z]\d[A-Za-z]\s*\d[A-Za-z]\d'
		ORDER BY created_at DESC
	`
	args := []interface{}{}
	if limit > 0 {
		query += " LIMIT $1"
		args = append(args, limit)
	}

	var employers []*models.LMIAEmployer
	err := r.db.Select(&employers, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get employers needing postal code extraction: %w", err)
	}

	return employers, nil
}

// UpdateEmployerPostalCode updates the postal code for an employer
func (r *lmiaRepository) UpdateEmployerPostalCode(id string, postalCode string) error {
	query := `UPDATE lmia_employers SET postal_code = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, id, postalCode)
	if err != nil {
		return fmt.Errorf("failed to update employer postal code: %w", err)
	}
	return nil
}

// GetUngeocodedPostalCodes returns a map of postal codes to their most common province for geocoding validation
func (r *lmiaRepository) GetUngeocodedPostalCodes() (map[string]string, error) {
	// First get count of postal codes that need geocoding (postal codes in lmia_employers but not in postal_codes table)
	var postalCodeCount int
	countQuery := `
		SELECT COUNT(DISTINCT e.postal_code)
		FROM lmia_employers e
		LEFT JOIN postal_codes pc ON REPLACE(e.postal_code, ' ', '') = REPLACE(pc.postal_code, ' ', '')
		WHERE e.postal_code IS NOT NULL
		AND e.postal_code != ''
		AND pc.postal_code IS NULL
	`
	err := r.db.Get(&postalCodeCount, countQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to count ungeocoded postal codes: %w", err)
	}

	// Get postal codes with their most common province/territory that don't exist in postal_codes table
	query := `
		SELECT 
			REPLACE(e.postal_code, ' ', '') as normalized_postal_code,
			COALESCE(e.province_territory, '') as province_territory,
			COUNT(*) as count
		FROM lmia_employers e
		LEFT JOIN postal_codes pc ON REPLACE(e.postal_code, ' ', '') = REPLACE(pc.postal_code, ' ', '')
		WHERE e.postal_code IS NOT NULL
		AND e.postal_code != ''
		AND pc.postal_code IS NULL
		GROUP BY REPLACE(e.postal_code, ' ', ''), COALESCE(e.province_territory, '')
		ORDER BY normalized_postal_code, count DESC
	`

	type PostalCodeProvince struct {
		PostalCode string `db:"normalized_postal_code"`
		Province   string `db:"province_territory"`
		Count      int    `db:"count"`
	}

	var results []PostalCodeProvince
	err = r.db.Select(&results, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get ungeocoded postal codes with provinces: %w", err)
	}

	// Build map with most common province per postal code
	postalCodeProvinces := make(map[string]string)
	for _, result := range results {
		if _, exists := postalCodeProvinces[result.PostalCode]; !exists {
			// First occurrence (highest count due to ORDER BY) becomes the representative province
			postalCodeProvinces[result.PostalCode] = result.Province
		}
	}

	log.Info("Retrieved postal codes for geocoding", 
		"postal_codes_needing_geocoding", postalCodeCount,
		"unique_postal_codes", len(postalCodeProvinces))

	return postalCodeProvinces, nil
}


