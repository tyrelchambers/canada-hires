CREATE TABLE IF NOT EXISTS scraper_jobs (
    id SERIAL PRIMARY KEY,
    job_type VARCHAR(50) NOT NULL DEFAULT 'lmia_scraper',
    last_run_at TIMESTAMP WITH TIME ZONE,
    next_scheduled_run TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index for efficient querying
CREATE INDEX idx_scraper_jobs_type_last_run ON scraper_jobs(job_type, last_run_at);

-- Insert initial record for LMIA scraper job tracking
INSERT INTO scraper_jobs (job_type, next_scheduled_run, status) 
VALUES ('lmia_scraper', NOW(), 'pending')
ON CONFLICT DO NOTHING;