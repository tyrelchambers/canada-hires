package repos

import (
	"canada-hires/models"
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

type SubredditRepository interface {
	GetAll() ([]*models.Subreddit, error)
	GetActive() ([]*models.Subreddit, error)
	GetByID(id string) (*models.Subreddit, error)
	GetByIDs(ids []string) ([]*models.Subreddit, error)
	GetByName(name string) (*models.Subreddit, error)
	Create(subreddit *models.Subreddit) error
	Update(id string, updates *models.UpdateSubredditRequest) error
	Delete(id string) error
	IncrementPostCount(id string) error
	UpdateLastPostedAt(id string, postedAt time.Time) error
}

type subredditRepository struct {
	db *sqlx.DB
}

func NewSubredditRepository(db *sqlx.DB) SubredditRepository {
	return &subredditRepository{db: db}
}

func (r *subredditRepository) GetAll() ([]*models.Subreddit, error) {
	subreddits := []*models.Subreddit{}
	query := `
		SELECT id, name, is_active, post_count, 
		       last_posted_at, created_at, updated_at
		FROM subreddits 
		ORDER BY created_at ASC
	`
	
	err := r.db.Select(&subreddits, query)
	if err != nil {
		log.Error("Failed to get all subreddits", "error", err)
		return nil, err
	}
	
	return subreddits, nil
}

func (r *subredditRepository) GetActive() ([]*models.Subreddit, error) {
	subreddits := []*models.Subreddit{}
	query := `
		SELECT id, name, is_active, post_count, 
		       last_posted_at, created_at, updated_at
		FROM subreddits 
		WHERE is_active = true
		ORDER BY created_at ASC
	`
	
	err := r.db.Select(&subreddits, query)
	if err != nil {
		log.Error("Failed to get active subreddits", "error", err)
		return nil, err
	}
	
	return subreddits, nil
}

func (r *subredditRepository) GetByID(id string) (*models.Subreddit, error) {
	subreddit := &models.Subreddit{}
	query := `
		SELECT id, name, is_active, post_count, 
		       last_posted_at, created_at, updated_at
		FROM subreddits 
		WHERE id = $1
	`
	
	err := r.db.Get(subreddit, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("Failed to get subreddit by ID", "error", err, "id", id)
		return nil, err
	}
	
	return subreddit, nil
}

func (r *subredditRepository) GetByIDs(ids []string) ([]*models.Subreddit, error) {
	if len(ids) == 0 {
		return []*models.Subreddit{}, nil
	}

	subreddits := []*models.Subreddit{}
	query, args, err := sqlx.In(`
		SELECT id, name, is_active, post_count, 
		       last_posted_at, created_at, updated_at
		FROM subreddits 
		WHERE id IN (?)
		ORDER BY created_at ASC
	`, ids)
	
	if err != nil {
		log.Error("Failed to build IN query for subreddits", "error", err)
		return nil, err
	}
	
	query = r.db.Rebind(query)
	err = r.db.Select(&subreddits, query, args...)
	if err != nil {
		log.Error("Failed to get subreddits by IDs", "error", err, "ids", ids)
		return nil, err
	}
	
	return subreddits, nil
}

func (r *subredditRepository) GetByName(name string) (*models.Subreddit, error) {
	subreddit := &models.Subreddit{}
	query := `
		SELECT id, name, is_active, post_count, 
		       last_posted_at, created_at, updated_at
		FROM subreddits 
		WHERE name = $1
	`
	
	err := r.db.Get(subreddit, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Error("Failed to get subreddit by name", "error", err, "name", name)
		return nil, err
	}
	
	return subreddit, nil
}

func (r *subredditRepository) Create(subreddit *models.Subreddit) error {
	query := `
		INSERT INTO subreddits (id, name, is_active, created_at, updated_at)
		VALUES (:id, :name, :is_active, :created_at, :updated_at)
	`
	
	now := time.Now()
	subreddit.CreatedAt = now
	subreddit.UpdatedAt = now
	
	_, err := r.db.NamedExec(query, subreddit)
	if err != nil {
		log.Error("Failed to create subreddit", "error", err, "name", subreddit.Name)
		return err
	}
	
	log.Info("Subreddit created successfully", "id", subreddit.ID, "name", subreddit.Name)
	return nil
}

func (r *subredditRepository) Update(id string, updates *models.UpdateSubredditRequest) error {
	setParts := []string{}
	args := map[string]interface{}{
		"id":         id,
		"updated_at": time.Now(),
	}
	
	if updates.IsActive != nil {
		setParts = append(setParts, "is_active = :is_active")
		args["is_active"] = updates.IsActive
	}
	
	if len(setParts) == 0 {
		return nil // Nothing to update
	}
	
	setParts = append(setParts, "updated_at = :updated_at")
	
	query := `UPDATE subreddits SET ` + setParts[0]
	for i := 1; i < len(setParts); i++ {
		query += ", " + setParts[i]
	}
	query += ` WHERE id = :id`
	
	_, err := r.db.NamedExec(query, args)
	if err != nil {
		log.Error("Failed to update subreddit", "error", err, "id", id)
		return err
	}
	
	log.Info("Subreddit updated successfully", "id", id)
	return nil
}

func (r *subredditRepository) Delete(id string) error {
	query := `DELETE FROM subreddits WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("Failed to delete subreddit", "error", err, "id", id)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	log.Info("Subreddit deleted successfully", "id", id)
	return nil
}

func (r *subredditRepository) IncrementPostCount(id string) error {
	query := `UPDATE subreddits SET post_count = post_count + 1, updated_at = NOW() WHERE id = $1`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error("Failed to increment post count", "error", err, "id", id)
		return err
	}
	
	return nil
}

func (r *subredditRepository) UpdateLastPostedAt(id string, postedAt time.Time) error {
	query := `UPDATE subreddits SET last_posted_at = $1, updated_at = NOW() WHERE id = $2`
	
	_, err := r.db.Exec(query, postedAt, id)
	if err != nil {
		log.Error("Failed to update last posted at", "error", err, "id", id)
		return err
	}
	
	return nil
}