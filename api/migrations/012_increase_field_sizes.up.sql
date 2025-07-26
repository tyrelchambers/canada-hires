-- Increase field sizes for better compatibility with scraper data

-- Increase province field to handle full province names
ALTER TABLE job_postings 
ALTER COLUMN province TYPE VARCHAR(50);

-- Increase city field size for longer city names
ALTER TABLE job_postings 
ALTER COLUMN city TYPE VARCHAR(150);

-- Increase salary_type field size for longer descriptions
ALTER TABLE job_postings 
ALTER COLUMN salary_type TYPE VARCHAR(50);