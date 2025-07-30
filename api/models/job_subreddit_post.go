package models

import (
	"time"
)

type JobSubredditPost struct {
	ID            string     `json:"id" db:"id"`
	JobPostingID  string     `json:"job_posting_id" db:"job_posting_id"`
	SubredditID   string     `json:"subreddit_id" db:"subreddit_id"`
	RedditPostID  *string    `json:"reddit_post_id" db:"reddit_post_id"`     // Reddit's post ID (e.g., "abc123")
	RedditPostURL *string    `json:"reddit_post_url" db:"reddit_post_url"`   // Full URL to the Reddit post
	PostedAt      time.Time  `json:"posted_at" db:"posted_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	
	// Joined data when querying with job details
	SubredditName *string `json:"subreddit_name,omitempty" db:"subreddit_name"`
}

type CreateJobSubredditPostRequest struct {
	JobPostingID  string  `json:"job_posting_id" validate:"required"`
	SubredditID   string  `json:"subreddit_id" validate:"required"`
	RedditPostID  *string `json:"reddit_post_id,omitempty"`
	RedditPostURL *string `json:"reddit_post_url,omitempty"`
}