-- Remove reddit_posted column from job_postings table

DROP INDEX IF EXISTS idx_job_postings_reddit_posted;

ALTER TABLE job_postings 
DROP COLUMN reddit_posted;