package repos

import (
	"canada-hires/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SessionRepository interface {
	Create(session *models.Session) error
	GetByID(id string) (*models.Session, error)
	DeleteByID(id string) error
	DeleteByUserID(userID string) error
	DeleteExpired() error
	UpdateLastUsed(id string) error
}

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *models.Session) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO sessions (id, user_id, expires_at, ip_address, user_agent, created_at, updated_at)
		VALUES (:id, :user_id, :expires_at, :ip_address, :user_agent, :created_at, :updated_at)
	`

	session.ID = uuid.New().String()
	session.CreatedAt = time.Now().UTC()
	session.UpdatedAt = time.Now().UTC()

	_, err = tx.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetByID(id string) (*models.Session, error) {
	var session models.Session
	query := `SELECT * FROM sessions WHERE id = $1`

	err := r.db.Get(&session, query, id)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) DeleteByID(id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM sessions WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteByUserID(userID string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err = tx.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions for user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *sessionRepository) DeleteExpired() error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err = tx.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *sessionRepository) UpdateLastUsed(id string) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE sessions SET updated_at = NOW() WHERE id = $1`
	_, err = tx.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
