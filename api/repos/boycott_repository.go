package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BoycottRepository interface {
	Toggle(userID, businessName, businessAddress string) (*models.Boycott, error)
	GetByUserID(userID string, limit, offset int) ([]*models.Boycott, error)
	GetTopBoycotted(limit int) ([]*models.BoycottStats, error)
	GetBoycottCount(businessName, businessAddress string) (int, error)
	IsBoycottedByUser(userID, businessName, businessAddress string) (bool, error)
	Delete(userID, businessName, businessAddress string) error
}

type boycottRepository struct {
	db *sqlx.DB
}

func NewBoycottRepository(db *sqlx.DB) BoycottRepository {
	return &boycottRepository{db: db}
}

func (r *boycottRepository) Toggle(userID, businessName, businessAddress string) (*models.Boycott, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if boycott already exists
	var existing models.Boycott
	var checkQuery string
	var checkArgs []interface{}

	if businessAddress == "" {
		checkQuery = `SELECT * FROM boycotts WHERE user_id = $1 AND business_name = $2 AND (business_address IS NULL OR business_address = '')`
		checkArgs = []interface{}{userID, businessName}
	} else {
		checkQuery = `SELECT * FROM boycotts WHERE user_id = $1 AND business_name = $2 AND business_address = $3`
		checkArgs = []interface{}{userID, businessName, businessAddress}
	}

	err = tx.Get(&existing, checkQuery, checkArgs...)

	if err == nil {
		// Boycott exists, delete it (toggle off)
		deleteQuery := `DELETE FROM boycotts WHERE id = $1`
		_, err = tx.Exec(deleteQuery, existing.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete boycott: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return nil, nil // Return nil to indicate boycott was removed
	}

	// Boycott doesn't exist, create it (toggle on)
	boycott := &models.Boycott{
		ID:              uuid.New().String(),
		UserID:          userID,
		BusinessName:    businessName,
		BusinessAddress: &businessAddress,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	insertQuery := `
		INSERT INTO boycotts (id, user_id, business_name, business_address, created_at, updated_at)
		VALUES (:id, :user_id, :business_name, :business_address, :created_at, :updated_at)
	`

	_, err = tx.NamedExec(insertQuery, boycott)
	if err != nil {
		return nil, fmt.Errorf("failed to insert boycott: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return boycott, nil
}

func (r *boycottRepository) GetByUserID(userID string, limit, offset int) ([]*models.Boycott, error) {
	var boycotts []*models.Boycott
	query := `SELECT * FROM boycotts WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	err := r.db.Select(&boycotts, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get boycotts by user ID: %w", err)
	}

	return boycotts, nil
}

func (r *boycottRepository) GetTopBoycotted(limit int) ([]*models.BoycottStats, error) {
	var stats []*models.BoycottStats
	query := `
		SELECT
			business_name,
			COALESCE(business_address, '') as business_address,
			COUNT(*) as boycott_count
		FROM boycotts
		GROUP BY business_name, business_address
		ORDER BY boycott_count DESC
		LIMIT $1
	`

	err := r.db.Select(&stats, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get top boycotted businesses: %w", err)
	}

	return stats, nil
}

func (r *boycottRepository) GetBoycottCount(businessName, businessAddress string) (int, error) {
	var count int
	var query string
	var args []interface{}

	if businessAddress == "" {
		query = `SELECT COUNT(*) FROM boycotts WHERE business_name = $1 AND (business_address IS NULL OR business_address = '')`
		args = []interface{}{businessName}
	} else {
		query = `SELECT COUNT(*) FROM boycotts WHERE business_name = $1 AND business_address = $2`
		args = []interface{}{businessName, businessAddress}
	}

	err := r.db.Get(&count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get boycott count: %w", err)
	}

	return count, nil
}

func (r *boycottRepository) IsBoycottedByUser(userID, businessName, businessAddress string) (bool, error) {
	var count int
	var query string
	var args []interface{}

	if businessAddress == "" {
		query = `SELECT COUNT(*) FROM boycotts WHERE user_id = $1 AND business_name = $2 AND (business_address IS NULL OR business_address = '')`
		args = []interface{}{userID, businessName}
	} else {
		query = `SELECT COUNT(*) FROM boycotts WHERE user_id = $1 AND business_name = $2 AND business_address = $3`
		args = []interface{}{userID, businessName, businessAddress}
	}

	fmt.Println(count)

	err := r.db.Get(&count, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check boycott status: %w", err)
	}

	return count > 0, nil
}

func (r *boycottRepository) Delete(userID, businessName, businessAddress string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var query string
	var args []interface{}

	if businessAddress == "" {
		query = `DELETE FROM boycotts WHERE user_id = $1 AND business_name = $2 AND (business_address IS NULL OR business_address = '')`
		args = []interface{}{userID, businessName}
	} else {
		query = `DELETE FROM boycotts WHERE user_id = $1 AND business_name = $2 AND business_address = $3`
		args = []interface{}{userID, businessName, businessAddress}
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete boycott: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
