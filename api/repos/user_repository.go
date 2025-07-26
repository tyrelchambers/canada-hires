package repos

import (
	"canada-hires/models"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	UpdateLastActive(id string) error
	AddIPAddress(userID string, ipAddress string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO users (id, email, verification_tier, email_domain, ip_addresses, created_at, last_active, updated_at)
		VALUES (:id, :email, :verification_tier, :email_domain, :ip_addresses, :created_at, :last_active, :updated_at)
	`

	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.LastActive = time.Now()
	user.UpdatedAt = time.Now()
	user.VerificationTier = models.VerificationBasic

	// Extract email domain
	if emailParts := strings.Split(user.Email, "@"); len(emailParts) == 2 {
		user.EmailDomain = &emailParts[1]
	}

	_, err = tx.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateLastActive(id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET last_active = NOW(), updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update last active: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *userRepository) AddIPAddress(userID string, ipAddress string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE users
		SET ip_addresses = CASE
			WHEN ip_addresses @> $2::jsonb THEN ip_addresses
			ELSE ip_addresses || $2::jsonb
		END,
		updated_at = NOW()
		WHERE id = $1
	`

	ipJSON := fmt.Sprintf(`["%s"]`, ipAddress)
	_, err = tx.Exec(query, userID, ipJSON)
	if err != nil {
		return fmt.Errorf("failed to add IP address: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
