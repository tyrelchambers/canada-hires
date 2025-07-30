-- Clean up salary_raw data by removing "Salary" prefix
UPDATE job_postings 
SET salary_raw = REPLACE(salary_raw, 'Salary', '')
WHERE salary_raw LIKE 'Salary%';