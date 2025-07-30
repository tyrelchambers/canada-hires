package services

import (
	"canada-hires/helpers"
	"canada-hires/models"
	"canada-hires/repos"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
)

type RedditPostInfo struct {
	RedditPostID  *string
	RedditPostURL *string
}

type RedditService interface {
	PostJob(ctx context.Context, job *models.JobPosting) error
	PostJobWithConfig(ctx context.Context, job *models.JobPosting, config *models.RedditConfig, subreddit *models.Subreddit) (*RedditPostInfo, error)
	TestConnection(ctx context.Context) error
	GetDefaultConfig() *models.RedditConfig
	SetConfig(config *models.RedditConfig)
}

type redditService struct {
	authService         RedditAuthService
	config              *models.RedditConfig
	logger              *log.Logger
	jobRepo             repos.JobBankRepository
	subredditRepo       repos.SubredditRepository
	jobSubredditPostRepo repos.JobSubredditPostRepository
	httpClient          *http.Client
	userAgent           string
}

type RedditSubmitResponse struct {
	JSON struct {
		Errors [][]string `json:"errors"`
		Data   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"data"`
	} `json:"json"`
}

func NewRedditService(logger *log.Logger, jobRepo repos.JobBankRepository, subredditRepo repos.SubredditRepository, jobSubredditPostRepo repos.JobSubredditPostRepository) (RedditService, error) {
	// Create Reddit auth service
	authService, err := NewRedditAuthService(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Reddit auth service: %w", err)
	}

	userAgent := os.Getenv("REDDIT_USER_AGENT")
	if userAgent == "" {
		userAgent = "JobWatchCanada/1.0"
	}

	service := &redditService{
		authService:         authService,
		config:              models.DefaultRedditConfig(),
		logger:              logger,
		jobRepo:             jobRepo,
		subredditRepo:       subredditRepo,
		jobSubredditPostRepo: jobSubredditPostRepo,
		httpClient:          &http.Client{Timeout: 30 * time.Second},
		userAgent:           userAgent,
	}

	return service, nil
}

// PostJob posts a job to all active subreddits
func (rs *redditService) PostJob(ctx context.Context, job *models.JobPosting) error {
	// Get all active subreddits
	subreddits, err := rs.subredditRepo.GetActive()
	if err != nil {
		return fmt.Errorf("failed to get active subreddits: %w", err)
	}

	if len(subreddits) == 0 {
		rs.logger.Warn("No active subreddits found for posting")
		return nil
	}

	// Post to each active subreddit
	for _, subreddit := range subreddits {
		// Create a config for this specific subreddit
		config := rs.config

		rs.logger.Info("Posting job to subreddit",
			"subreddit", subreddit.Name,
			"job_title", job.Title,
			"employer", job.Employer,
		)

		redditPostInfo, err := rs.PostJobWithConfig(ctx, job, config, subreddit)
		if err != nil {
			rs.logger.Error("Failed to post job to subreddit",
				"error", err,
				"subreddit", subreddit.Name,
				"job_id", job.ID,
			)
			// Continue with other subreddits even if one fails
			continue
		}

		// Create job subreddit post record
		jobSubredditPost := &models.JobSubredditPost{
			ID:            uuid.New().String(),
			JobPostingID:  job.ID,
			SubredditID:   subreddit.ID,
			RedditPostID:  redditPostInfo.RedditPostID,
			RedditPostURL: redditPostInfo.RedditPostURL,
		}

		if err := rs.jobSubredditPostRepo.Create(jobSubredditPost); err != nil {
			rs.logger.Error("Failed to create job subreddit post record",
				"error", err,
				"job_id", job.ID,
				"subreddit_id", subreddit.ID,
			)
			// Don't fail the entire operation, just log the error
		}

		// Update subreddit statistics
		now := time.Now()
		if err := rs.subredditRepo.IncrementPostCount(subreddit.ID); err != nil {
			rs.logger.Error("Failed to increment post count", "error", err, "subreddit_id", subreddit.ID)
		}
		if err := rs.subredditRepo.UpdateLastPostedAt(subreddit.ID, now); err != nil {
			rs.logger.Error("Failed to update last posted at", "error", err, "subreddit_id", subreddit.ID)
		}

		rs.logger.Info("Successfully posted job to subreddit",
			"subreddit", subreddit.Name,
			"job_id", job.ID,
			"reddit_post_id", redditPostInfo.RedditPostID,
		)
	}

	return nil
}

// PostJobWithConfig posts a job using a specific configuration
func (rs *redditService) PostJobWithConfig(ctx context.Context, job *models.JobPosting, config *models.RedditConfig, subreddit *models.Subreddit) (*RedditPostInfo, error) {
	if !config.IsEnabled {
		rs.logger.Debug("Reddit posting disabled in configuration")
		return &RedditPostInfo{}, nil
	}

	if err := config.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("invalid Reddit configuration: %w", err)
	}

	// Generate post data from template
	postData := config.GeneratePostData(job)
	if postData == nil {
		return nil, fmt.Errorf("failed to generate post data")
	}

	// Override the subreddit with the specific one we're posting to
	postData.Subreddit = subreddit.Name

	// Debug: Print post structure in development mode
	if helpers.IsDev() {
		rs.logger.Info("=== DEBUG: Reddit Post Structure ===")
		rs.logger.Info("Post Title:", "title", postData.Title)
		rs.logger.Info("Post Subreddit:", "subreddit", postData.Subreddit)
		rs.logger.Info("Post Body:")
		fmt.Println(postData.Body)
		rs.logger.Info("=== End Reddit Post Structure ===")

		rs.logger.Debug("Not posting to reddit since we are in DEV mode")
		return &RedditPostInfo{}, nil
	}

	// Validate post data
	if err := rs.validatePostData(postData); err != nil {
		return nil, fmt.Errorf("invalid post data: %w", err)
	}

	rs.logger.Info("Posting job to Reddit",
		"subreddit", postData.Subreddit,
		"job_title", job.Title,
		"employer", job.Employer,
	)

	// Get access token before posting
	accessToken, err := rs.authService.GetAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Reddit access token: %w", err)
	}

	// Submit the post to Reddit using the API directly
	submitResp, err := rs.submitTextPost(ctx, accessToken, postData)
	if err != nil {
		return nil, fmt.Errorf("failed to submit post to Reddit: %w", err)
	}

	// Create Reddit post info from response
	redditPostInfo := &RedditPostInfo{}
	if submitResp.JSON.Data.Name != "" {
		redditPostInfo.RedditPostID = &submitResp.JSON.Data.Name
	}
	if submitResp.JSON.Data.URL != "" {
		redditPostInfo.RedditPostURL = &submitResp.JSON.Data.URL
	}

	rs.logger.Info("Successfully posted job to Reddit",
		"subreddit", postData.Subreddit,
		"job_id", job.ID,
		"post_title", postData.Title,
		"reddit_post_id", redditPostInfo.RedditPostID,
	)

	// Mark job as posted to Reddit
	if err := rs.jobRepo.UpdateJobPostingRedditStatus(job.ID, true); err != nil {
		rs.logger.Error("Failed to update Reddit posted status",
			"error", err,
			"job_id", job.ID,
		)
		// Don't return error as the Reddit posting was successful
	}

	return redditPostInfo, nil
}

// submitTextPost submits a text post to Reddit using the authenticated API
func (rs *redditService) submitTextPost(ctx context.Context, accessToken string, postData *models.RedditPostData) (*RedditSubmitResponse, error) {
	// Prepare form data
	data := url.Values{}
	data.Set("sr", postData.Subreddit)
	data.Set("kind", "self")
	data.Set("title", postData.Title)
	data.Set("text", postData.Body)
	data.Set("api_type", "json")

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://oauth.reddit.com/api/submit", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", rs.userAgent)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request
	resp, err := rs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var submitResp RedditSubmitResponse
	if err := json.NewDecoder(resp.Body).Decode(&submitResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for errors in the response
	if len(submitResp.JSON.Errors) > 0 {
		errMsg := "Reddit API errors: "
		for _, err := range submitResp.JSON.Errors {
			if len(err) > 0 {
				errMsg += strings.Join(err, " ") + "; "
			}
		}
		return nil, fmt.Errorf(errMsg)
	}

	// Log successful submission details
	if submitResp.JSON.Data.Name != "" {
		rs.logger.Info("Reddit post created successfully",
			"reddit_id", submitResp.JSON.Data.Name,
			"reddit_url", submitResp.JSON.Data.URL,
		)
	}

	return &submitResp, nil
}

// TestConnection tests the Reddit API connection
func (rs *redditService) TestConnection(ctx context.Context) error {
	rs.logger.Info("Testing Reddit API connection")

	// Test by getting an access token
	_, err := rs.authService.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("Reddit authentication failed: %w", err)
	}

	rs.logger.Info("Reddit API authentication successful")
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
