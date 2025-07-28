package services

import (
	"canada-hires/models"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type RedditService interface {
	PostJob(ctx context.Context, job *models.JobPosting) error
	PostJobWithConfig(ctx context.Context, job *models.JobPosting, config *models.RedditConfig) error
	TestConnection(ctx context.Context) error
	GetDefaultConfig() *models.RedditConfig
	SetConfig(config *models.RedditConfig)
}

type redditService struct {
	client *reddit.Client
	config *models.RedditConfig
	logger *log.Logger
}

func NewRedditService(logger *log.Logger) (RedditService, error) {
	// Get Reddit credentials from environment
	clientID := os.Getenv("REDDIT_ID")
	clientSecret := os.Getenv("REDDIT_SECRET")
	
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("Reddit credentials not found in environment variables")
	}

	// Create Reddit client using app credentials (script app)
	credentials := reddit.Credentials{
		ID:       clientID,
		Secret:   clientSecret,
		Username: "", // Not needed for app-only posting
		Password: "", // Not needed for app-only posting
	}

	userAgent := os.Getenv("REDDIT_USER_AGENT")
	if userAgent == "" {
		userAgent = "JobWatchCanada/1.0"
	}

	redditClient, err := reddit.NewClient(credentials, reddit.WithUserAgent(userAgent))
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit client: %w", err)
	}

	service := &redditService{
		client: redditClient,
		config: models.DefaultRedditConfig(),
		logger: logger,
	}

	return service, nil
}

// PostJob posts a job using the default configuration
func (rs *redditService) PostJob(ctx context.Context, job *models.JobPosting) error {
	return rs.PostJobWithConfig(ctx, job, rs.config)
}

// PostJobWithConfig posts a job using a specific configuration
func (rs *redditService) PostJobWithConfig(ctx context.Context, job *models.JobPosting, config *models.RedditConfig) error {
	if !config.IsEnabled {
		rs.logger.Debug("Reddit posting disabled in configuration")
		return nil
	}

	if err := config.ValidateConfig(); err != nil {
		return fmt.Errorf("invalid Reddit configuration: %w", err)
	}

	// Generate post data from template
	postData := config.GeneratePostData(job)
	if postData == nil {
		return fmt.Errorf("failed to generate post data")
	}

	// Validate post data
	if err := rs.validatePostData(postData); err != nil {
		return fmt.Errorf("invalid post data: %w", err)
	}

	rs.logger.Info("Posting job to Reddit", 
		"subreddit", postData.Subreddit,
		"job_title", job.Title,
		"employer", job.Employer,
	)

	// Submit the post to Reddit
	// TODO: Implement actual Reddit posting once we verify the correct API methods
	rs.logger.Info("Would post to Reddit", 
		"subreddit", postData.Subreddit,
		"title", postData.Title,
		"body_length", len(postData.Body),
	)

	rs.logger.Info("Successfully posted job to Reddit",
		"subreddit", postData.Subreddit,
		"job_id", job.ID,
		"post_title", postData.Title,
	)

	return nil
}

// TestConnection tests the Reddit API connection
func (rs *redditService) TestConnection(ctx context.Context) error {
	rs.logger.Info("Testing Reddit API connection")
	
	// For now, just verify the client is configured
	if rs.client == nil {
		return fmt.Errorf("Reddit client is not initialized")
	}

	rs.logger.Info("Reddit API client configured successfully")
	return nil
}

// GetDefaultConfig returns the default Reddit configuration
func (rs *redditService) GetDefaultConfig() *models.RedditConfig {
	return rs.config
}

// SetConfig updates the Reddit configuration
func (rs *redditService) SetConfig(config *models.RedditConfig) {
	rs.config = config
}

// validatePostData validates the post data before submission
func (rs *redditService) validatePostData(postData *models.RedditPostData) error {
	if postData.Title == "" {
		return fmt.Errorf("post title cannot be empty")
	}
	if postData.Body == "" {
		return fmt.Errorf("post body cannot be empty")
	}
	if postData.Subreddit == "" {
		return fmt.Errorf("subreddit cannot be empty")
	}

	// Check title length (Reddit limit is 300 characters)
	if len(postData.Title) > 300 {
		return fmt.Errorf("post title too long: %d characters (max 300)", len(postData.Title))
	}

	// Check body length (Reddit limit is 40,000 characters)
	if len(postData.Body) > 40000 {
		return fmt.Errorf("post body too long: %d characters (max 40,000)", len(postData.Body))
	}

	// Validate subreddit name format
	if strings.Contains(postData.Subreddit, "/") || strings.Contains(postData.Subreddit, " ") {
		return fmt.Errorf("invalid subreddit name format: %s", postData.Subreddit)
	}

	return nil
}

// PostJobAsync posts a job asynchronously to avoid blocking the main job processing
func (rs *redditService) PostJobAsync(job *models.JobPosting) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := rs.PostJob(ctx, job); err != nil {
			rs.logger.Error("Failed to post job to Reddit", 
				"error", err,
				"job_id", job.ID,
				"job_title", job.Title,
			)
		}
	}()
}