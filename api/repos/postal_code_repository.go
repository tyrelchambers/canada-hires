package repos

import (
	"canada-hires/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostalCodeRepository interface {
	GetByPostalCode(postalCode string) (*models.PostalCodeCoordinates, error)
	GetAllPostalCodes() (map[string]models.PostalCodeCoordinates, error)
	Create(postalCode *models.PostalCodeCoordinates) error
	Update(postalCode *models.PostalCodeCoordinates) error
	Upsert(postalCode *models.PostalCodeCoordinates) error
}

type postalCodeRepository struct {
	db *sqlx.DB
}

func NewPostalCodeRepository(db *sqlx.DB) PostalCodeRepository {
	return &postalCodeRepository{db: db}
}

func (r *postalCodeRepository) GetByPostalCode(postalCode string) (*models.PostalCodeCoordinates, error) {
	var result models.PostalCodeCoordinates
	
	query := `
		SELECT postal_code, latitude, longitude
		FROM postal_codes 
		WHERE postal_code = $1
	`
	
	err := r.db.Get(&result, query, postalCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get postal code coordinates: %w", err)
	}
	
	return &result, nil
}

func (r *postalCodeRepository) Create(postalCode *models.PostalCodeCoordinates) error {
	// Only save postal codes that have valid coordinates
	if !postalCode.Latitude.Valid || !postalCode.Longitude.Valid {
		return fmt.Errorf("postal code %s has invalid coordinates, not saving", postalCode.PostalCode)
	}
	
	query := `
		INSERT INTO postal_codes (postal_code, latitude, longitude, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`
	
	_, err := r.db.Exec(query, postalCode.PostalCode, postalCode.Latitude, postalCode.Longitude)
	if err != nil {
		return fmt.Errorf("failed to create postal code: %w", err)
	}
	
	return nil
}

func (r *postalCodeRepository) Update(postalCode *models.PostalCodeCoordinates) error {
	query := `
		UPDATE postal_codes 
		SET latitude = $2, longitude = $3, updated_at = NOW()
		WHERE postal_code = $1
	`
	
	_, err := r.db.Exec(query, postalCode.PostalCode, postalCode.Latitude, postalCode.Longitude)
	if err != nil {
		return fmt.Errorf("failed to update postal code: %w", err)
	}
	
	return nil
}

func (r *postalCodeRepository) Upsert(postalCode *models.PostalCodeCoordinates) error {
	// Only save postal codes that have valid coordinates
	if !postalCode.Latitude.Valid || !postalCode.Longitude.Valid {
		return fmt.Errorf("postal code %s has invalid coordinates, not saving", postalCode.PostalCode)
	}
	
	query := `
		INSERT INTO postal_codes (postal_code, latitude, longitude, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (postal_code) 
		DO UPDATE SET 
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			updated_at = NOW()
	`
	
	_, err := r.db.Exec(query, postalCode.PostalCode, postalCode.Latitude, postalCode.Longitude)
	if err != nil {
		return fmt.Errorf("failed to upsert postal code: %w", err)
	}
	
	return nil
}

func (r *postalCodeRepository) GetAllPostalCodes() (map[string]models.PostalCodeCoordinates, error) {
	var rows []models.PostalCodeCoordinates
	
	query := `
		SELECT postal_code, latitude, longitude
		FROM postal_codes 
		WHERE latitude IS NOT NULL AND longitude IS NOT NULL
	`
	
	err := r.db.Select(&rows, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all postal codes: %w", err)
	}
	
	// Convert to map for fast lookup
	result := make(map[string]models.PostalCodeCoordinates)
	for _, row := range rows {
		result[row.PostalCode] = row
	}
	
	return result, nil
}

