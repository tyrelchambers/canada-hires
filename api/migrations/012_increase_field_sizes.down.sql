-- Rollback field size increases

-- Revert province field size (this may fail if data is too long)
-- ALTER TABLE job_postings 
-- ALTER COLUMN province TYPE VARCHAR(2);

-- Revert city field size (this may fail if data is too long)
-- ALTER TABLE job_postings 
-- ALTER COLUMN city TYPE VARCHAR(100);

-- Revert salary_type field size (this may fail if data is too long)
-- ALTER TABLE job_postings 
-- ALTER COLUMN salary_type TYPE VARCHAR(20);