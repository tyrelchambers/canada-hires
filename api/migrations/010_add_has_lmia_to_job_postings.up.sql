-- Add has_lmia column to job_postings table

ALTER TABLE job_postings 
ADD COLUMN has_lmia BOOLEAN NOT NULL DEFAULT FALSE;

-- Add index for LMIA filtering
CREATE INDEX idx_job_postings_has_lmia ON job_postings(has_lmia);