-- Create table to track which subreddits each job was posted to
CREATE TABLE IF NOT EXISTS job_subreddit_posts (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    job_posting_id VARCHAR(36) NOT NULL,
    subreddit_id VARCHAR(36) NOT NULL,
    reddit_post_id VARCHAR(50), -- Reddit post ID (e.g., "abc123")
    reddit_post_url TEXT, -- Full URL to the Reddit post
    posted_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Foreign key constraints
    FOREIGN KEY (subreddit_id) REFERENCES subreddits(id) ON DELETE CASCADE,

    -- Unique constraint to prevent duplicate posts
    UNIQUE(job_posting_id, subreddit_id)
);

-- Create indexes for efficient querying
CREATE INDEX idx_job_subreddit_posts_job_id ON job_subreddit_posts(job_posting_id);
CREATE INDEX idx_job_subreddit_posts_subreddit_id ON job_subreddit_posts(subreddit_id);
CREATE INDEX idx_job_subreddit_posts_posted_at ON job_subreddit_posts(posted_at DESC);
