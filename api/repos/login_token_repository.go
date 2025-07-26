package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type LoginTokenRepository interface {
	Create(token *models.LoginToken) error
	GetByToken(token string) (*models.LoginToken, error)
	MarkAsUsed(tokenID string) error
	DeleteExpiredTokens() error
	DeleteUserTokens(userID string) error
}

type loginTokenRepository struct {
	db *sqlx.DB
}

func NewLoginTokenRepository(db *sqlx.DB) LoginTokenRepository {
	return &loginTokenRepository{db: db}
}

func (r *loginTokenRepository) Create(token *models.LoginToken) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO login_tokens (id, user_id, token, expires_at, ip_address, created_at)
		VALUES (:id, :user_id, :token, :expires_at, :ip_address, :created_at)
	`

	token.ID = uuid.New().String()
	token.CreatedAt = time.Now()

	_, err = tx.NamedExec(query, token)
	if err != nil {
		return fmt.Errorf("failed to insert login token: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *loginTokenRepository) GetByToken(token string) (*models.LoginToken, error) {
	var loginToken models.LoginToken
	query := `SELECT * FROM login_tokens WHERE token = $1 AND expires_at > NOW() AND used_at IS NULL`

	err := r.db.Get(&loginToken, query, token)
	if err != nil {
		return nil, err
	}

	return &loginToken, nil
}

func (r *loginTokenRepository) MarkAsUsed(tokenID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE login_tokens SET used_at = NOW() WHERE id = $1`
	_, err = tx.Exec(query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *loginTokenRepository) DeleteExpiredTokens() error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM login_tokens WHERE expires_at < NOW()`
	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *loginTokenRepository) DeleteUserTokens(userID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM login_tokens WHERE user_id = $1`
	_, err = tx.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user tokens: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
