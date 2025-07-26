-- Update job_postings table to better support the new scraper data format

-- Make job_bank_id nullable and add unique constraint on url instead
ALTER TABLE job_postings 
ALTER COLUMN job_bank_id DROP NOT NULL;

-- Add unique constraint on url since that's what we'll use as the unique identifier
ALTER TABLE job_postings 
ADD CONSTRAINT unique_job_posting_url UNIQUE(url);

-- Add salary_raw column to store the original salary string from scraper
ALTER TABLE job_postings 
ADD COLUMN salary_raw TEXT;

-- Add index for salary_raw column for searching
CREATE INDEX idx_job_postings_salary_raw ON job_postings(salary_raw);