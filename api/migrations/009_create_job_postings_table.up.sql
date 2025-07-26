CREATE TABLE job_scraping_runs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT,
    total_pages INTEGER DEFAULT 0,
    jobs_scraped INTEGER DEFAULT 0,
    jobs_stored INTEGER DEFAULT 0,
    last_page_scraped INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE job_postings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_bank_id VARCHAR(255) UNIQUE NOT NULL,
    title TEXT NOT NULL,
    employer VARCHAR(500) NOT NULL,
    location VARCHAR(200) NOT NULL,
    province VARCHAR(2),
    city VARCHAR(100),
    salary_min DECIMAL(10,2),
    salary_max DECIMAL(10,2),
    salary_type VARCHAR(20),
    posting_date TIMESTAMP WITH TIME ZONE,
    url TEXT NOT NULL,
    is_tfw BOOLEAN DEFAULT TRUE,
    description TEXT,
    scraping_run_id UUID NOT NULL REFERENCES job_scraping_runs(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for efficient querying
CREATE INDEX idx_job_postings_employer ON job_postings(employer);
CREATE INDEX idx_job_postings_location ON job_postings(location);
CREATE INDEX idx_job_postings_province ON job_postings(province);
CREATE INDEX idx_job_postings_salary_min ON job_postings(salary_min);
CREATE INDEX idx_job_postings_posting_date ON job_postings(posting_date);
CREATE INDEX idx_job_postings_is_tfw ON job_postings(is_tfw);
CREATE INDEX idx_job_postings_scraping_run_id ON job_postings(scraping_run_id);
CREATE INDEX idx_job_scraping_runs_status ON job_scraping_runs(status);
CREATE INDEX idx_job_scraping_runs_started_at ON job_scraping_runs(started_at);