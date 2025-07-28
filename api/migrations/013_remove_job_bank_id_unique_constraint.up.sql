-- Remove the unique constraint on job_bank_id as we now use url as the unique identifier
ALTER TABLE job_postings DROP CONSTRAINT IF EXISTS job_postings_job_bank_id_key;