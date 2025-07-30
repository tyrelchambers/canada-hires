package models

import (
	"time"
)

type Subreddit struct {
	ID           string     `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`                   // Subreddit name (without r/)
	IsActive     bool       `json:"is_active" db:"is_active"`         // Whether posts should be made to this subreddit
	PostCount    int        `json:"post_count" db:"post_count"`       // Number of posts made to this subreddit
	LastPostedAt *time.Time `json:"last_posted_at" db:"last_posted_at"` // When the last post was made
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateSubredditRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	IsActive *bool  `json:"is_active,omitempty"`
}

type UpdateSubredditRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
}