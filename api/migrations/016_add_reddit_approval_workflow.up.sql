-- Add Reddit approval workflow fields to job_postings table

ALTER TABLE job_postings 
ADD COLUMN reddit_approval_status VARCHAR(20) NOT NULL DEFAULT 'pending',
ADD COLUMN reddit_approved_by VARCHAR(255),
ADD COLUMN reddit_approved_at TIMESTAMP,
ADD COLUMN reddit_rejection_reason TEXT;

-- Add index for approval status filtering
CREATE INDEX idx_job_postings_reddit_approval_status ON job_postings(reddit_approval_status);

-- Add index for approved jobs
CREATE INDEX idx_job_postings_reddit_approved_at ON job_postings(reddit_approved_at);

-- Add constraint to ensure valid approval statuses
ALTER TABLE job_postings 
ADD CONSTRAINT chk_reddit_approval_status 
CHECK (reddit_approval_status IN ('pending', 'approved', 'rejected'));

-- Comment on the new columns
COMMENT ON COLUMN job_postings.reddit_approval_status IS 'Status of Reddit posting approval: pending, approved, rejected';
COMMENT ON COLUMN job_postings.reddit_approved_by IS 'User ID or email of admin who approved/rejected the posting';
COMMENT ON COLUMN job_postings.reddit_approved_at IS 'Timestamp when approval decision was made';
COMMENT ON COLUMN job_postings.reddit_rejection_reason IS 'Reason provided when job was rejected for Reddit posting';