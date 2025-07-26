-- Rollback changes for job_postings table

-- Remove salary_raw column
ALTER TABLE job_postings 
DROP COLUMN IF EXISTS salary_raw;

-- Remove unique constraint on url
ALTER TABLE job_postings 
DROP CONSTRAINT IF EXISTS unique_job_posting_url;

-- Make job_bank_id NOT NULL again (this may fail if there are NULL values)
-- ALTER TABLE job_postings 
-- ALTER COLUMN job_bank_id SET NOT NULL;