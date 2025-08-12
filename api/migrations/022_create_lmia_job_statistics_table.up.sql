CREATE TABLE lmia_job_statistics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL,
    period_type VARCHAR(10) NOT NULL CHECK (period_type IN ('daily', 'monthly')),
    total_jobs INTEGER NOT NULL DEFAULT 0,
    unique_employers INTEGER NOT NULL DEFAULT 0,
    avg_salary_min DECIMAL(10,2),
    avg_salary_max DECIMAL(10,2),
    top_provinces JSONB,
    top_cities JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create unique index to prevent duplicate entries for same date and period
CREATE UNIQUE INDEX idx_lmia_job_statistics_date_period ON lmia_job_statistics(date, period_type);

-- Indexes for efficient querying
CREATE INDEX idx_lmia_job_statistics_date ON lmia_job_statistics(date);
CREATE INDEX idx_lmia_job_statistics_period_type ON lmia_job_statistics(period_type);
CREATE INDEX idx_lmia_job_statistics_total_jobs ON lmia_job_statistics(total_jobs);