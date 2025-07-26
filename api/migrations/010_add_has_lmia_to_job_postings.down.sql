-- Remove has_lmia column from job_postings table

DROP INDEX IF EXISTS idx_job_postings_has_lmia;

ALTER TABLE job_postings 
DROP COLUMN has_lmia;