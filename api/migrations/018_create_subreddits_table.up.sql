CREATE TABLE IF NOT EXISTS subreddits (
    id VARCHAR(36) PRIMARY KEY DEFAULT (gen_random_uuid()::text),
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    post_count INTEGER DEFAULT 0,
    last_posted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index for efficient querying
CREATE INDEX idx_subreddits_active ON subreddits(is_active);
CREATE INDEX idx_subreddits_name ON subreddits(name);