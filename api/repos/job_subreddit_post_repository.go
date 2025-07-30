package repos

import (
	"canada-hires/models"
	"database/sql"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"
)

type JobSubredditPostRepository interface {
	Create(jobSubredditPost *models.JobSubredditPost) error
	GetByJobID(jobID string) ([]*models.JobSubredditPost, error)
	GetBySubredditID(subredditID string) ([]*models.JobSubredditPost, error)
	UpdateRedditPostInfo(id string, redditPostID string, redditPostURL string) error
}

type jobSubredditPostRepository struct {
	db *sqlx.DB
}

func NewJobSubredditPostRepository(db *sqlx.DB) JobSubredditPostRepository {
	return &jobSubredditPostRepository{db: db}
}

func (r *jobSubredditPostRepository) Create(jobSubredditPost *models.JobSubredditPost) error {
	query := `
		INSERT INTO job_subreddit_posts (id, job_posting_id, subreddit_id, reddit_post_id, reddit_post_url, posted_at, created_at)
		VALUES (:id, :job_posting_id, :subreddit_id, :reddit_post_id, :reddit_post_url, :posted_at, :created_at)
	`
	
	now := time.Now()
	jobSubredditPost.PostedAt = now
	jobSubredditPost.CreatedAt = now
	
	_, err := r.db.NamedExec(query, jobSubredditPost)
	if err != nil {
		log.Error("Failed to create job subreddit post", 
			"error", err, 
			"job_id", jobSubredditPost.JobPostingID, 
			"subreddit_id", jobSubredditPost.SubredditID)
		return err
	}
	
	log.Info("Job subreddit post created successfully", 
		"id", jobSubredditPost.ID, 
		"job_id", jobSubredditPost.JobPostingID, 
		"subreddit_id", jobSubredditPost.SubredditID)
	return nil
}

func (r *jobSubredditPostRepository) GetByJobID(jobID string) ([]*models.JobSubredditPost, error) {
	posts := []*models.JobSubredditPost{}
	query := `
		SELECT jsp.id, jsp.job_posting_id, jsp.subreddit_id, jsp.reddit_post_id, 
		       jsp.reddit_post_url, jsp.posted_at, jsp.created_at, s.name as subreddit_name
		FROM job_subreddit_posts jsp
		LEFT JOIN subreddits s ON jsp.subreddit_id = s.id
		WHERE jsp.job_posting_id = $1
		ORDER BY jsp.posted_at DESC
	`
	
	err := r.db.Select(&posts, query, jobID)
	if err != nil {
		log.Error("Failed to get job subreddit posts by job ID", "error", err, "job_id", jobID)
		return nil, err
	}
	
	return posts, nil
}

func (r *jobSubredditPostRepository) GetBySubredditID(subredditID string) ([]*models.JobSubredditPost, error) {
	posts := []*models.JobSubredditPost{}
	query := `
		SELECT jsp.id, jsp.job_posting_id, jsp.subreddit_id, jsp.reddit_post_id, 
		       jsp.reddit_post_url, jsp.posted_at, jsp.created_at, s.name as subreddit_name
		FROM job_subreddit_posts jsp
		LEFT JOIN subreddits s ON jsp.subreddit_id = s.id
		WHERE jsp.subreddit_id = $1
		ORDER BY jsp.posted_at DESC
	`
	
	err := r.db.Select(&posts, query, subredditID)
	if err != nil {
		log.Error("Failed to get job subreddit posts by subreddit ID", "error", err, "subreddit_id", subredditID)
		return nil, err
	}
	
	return posts, nil
}

func (r *jobSubredditPostRepository) UpdateRedditPostInfo(id string, redditPostID string, redditPostURL string) error {
	query := `
		UPDATE job_subreddit_posts 
		SET reddit_post_id = $1, reddit_post_url = $2
		WHERE id = $3
	`
	
	result, err := r.db.Exec(query, redditPostID, redditPostURL, id)
	if err != nil {
		log.Error("Failed to update Reddit post info", 
			"error", err, 
			"id", id, 
			"reddit_post_id", redditPostID)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	log.Info("Reddit post info updated successfully", 
		"id", id, 
		"reddit_post_id", redditPostID, 
		"reddit_post_url", redditPostURL)
	return nil
}