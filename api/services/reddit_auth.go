package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type RedditAuthService interface {
	GetAccessToken(ctx context.Context) (string, error)
	RefreshTokenIfNeeded(ctx context.Context) error
}

type redditAuthService struct {
	clientID     string
	clientSecret string
	username     string
	password     string
	userAgent    string
	logger       *log.Logger
	
	// Token cache
	accessToken string
	tokenType   string
	expiresAt   time.Time
}

type RedditTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error,omitempty"`
}

func NewRedditAuthService(logger *log.Logger) (RedditAuthService, error) {
	clientID := os.Getenv("REDDIT_ID")
	clientSecret := os.Getenv("REDDIT_SECRET")
	username := os.Getenv("REDDIT_USERNAME")
	password := os.Getenv("REDDIT_PASSWORD")
	
	if clientID == "" || clientSecret == "" || username == "" || password == "" {
		return nil, fmt.Errorf("missing required Reddit credentials: REDDIT_ID, REDDIT_SECRET, REDDIT_USERNAME, REDDIT_PASSWORD")
	}
	
	userAgent := os.Getenv("REDDIT_USER_AGENT")
	if userAgent == "" {
		userAgent = "JobWatchCanada/1.0"
	}
	
	return &redditAuthService{
		clientID:     clientID,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
		userAgent:    userAgent,
		logger:       logger,
	}, nil
}

// GetAccessToken returns a valid access token, fetching a new one if needed
func (r *redditAuthService) GetAccessToken(ctx context.Context) (string, error) {
	// Check if we have a valid cached token
	if r.accessToken != "" && time.Now().Before(r.expiresAt.Add(-30*time.Second)) {
		r.logger.Debug("Using cached Reddit access token")
		return r.accessToken, nil
	}
	
	// Fetch a new token
	if err := r.fetchNewToken(ctx); err != nil {
		return "", err
	}
	
	return r.accessToken, nil
}

// RefreshTokenIfNeeded refreshes the token if it's expired or about to expire
func (r *redditAuthService) RefreshTokenIfNeeded(ctx context.Context) error {
	// Refresh if token is empty or expires within 30 seconds
	if r.accessToken == "" || time.Now().After(r.expiresAt.Add(-30*time.Second)) {
		return r.fetchNewToken(ctx)
	}
	return nil
}

// fetchNewToken gets a new access token from Reddit API
func (r *redditAuthService) fetchNewToken(ctx context.Context) error {
	r.logger.Info("Fetching new Reddit access token")
	
	// Prepare the request data
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", r.username)
	data.Set("password", r.password)
	
	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", r.userAgent)
	req.SetBasicAuth(r.clientID, r.clientSecret)
	
	// Make the request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()
	
	// Parse the response
	var tokenResp RedditTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	
	// Check for errors in the response
	if tokenResp.Error != "" {
		return fmt.Errorf("Reddit API error: %s", tokenResp.Error)
	}
	
	if tokenResp.AccessToken == "" {
		return fmt.Errorf("no access token in response")
	}
	
	// Cache the token
	r.accessToken = tokenResp.AccessToken
	r.tokenType = tokenResp.TokenType
	r.expiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	
	r.logger.Info("Successfully obtained Reddit access token", 
		"token_type", tokenResp.TokenType,
		"expires_in", tokenResp.ExpiresIn,
		"expires_at", r.expiresAt.Format(time.RFC3339),
	)
	
	return nil
}