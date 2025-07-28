-- Re-add the unique constraint on job_bank_id (if rolling back)
ALTER TABLE job_postings ADD CONSTRAINT job_postings_job_bank_id_key UNIQUE(job_bank_id);