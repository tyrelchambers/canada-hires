-- Add reddit_posted column to job_postings table

ALTER TABLE job_postings 
ADD COLUMN reddit_posted BOOLEAN NOT NULL DEFAULT FALSE;

-- Add index for Reddit posting filtering
CREATE INDEX idx_job_postings_reddit_posted ON job_postings(reddit_posted);