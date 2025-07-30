-- Remove Reddit approval workflow fields from job_postings table

-- Drop indexes
DROP INDEX IF EXISTS idx_job_postings_reddit_approval_status;
DROP INDEX IF EXISTS idx_job_postings_reddit_approved_at;

-- Drop constraint
ALTER TABLE job_postings 
DROP CONSTRAINT IF EXISTS chk_reddit_approval_status;

-- Drop columns
ALTER TABLE job_postings 
DROP COLUMN IF EXISTS reddit_approval_status,
DROP COLUMN IF EXISTS reddit_approved_by,
DROP COLUMN IF EXISTS reddit_approved_at,
DROP COLUMN IF EXISTS reddit_rejection_reason;