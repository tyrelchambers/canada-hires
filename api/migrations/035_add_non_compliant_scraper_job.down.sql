-- Remove non-compliant scraper job record
DELETE FROM scraper_jobs WHERE job_type = 'non_compliant_scraper';