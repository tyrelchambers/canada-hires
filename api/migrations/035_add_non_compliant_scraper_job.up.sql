-- Insert initial record for non-compliant scraper job tracking
INSERT INTO scraper_jobs (job_type, next_scheduled_run, status) 
VALUES ('non_compliant_scraper', NOW(), 'pending')
ON CONFLICT DO NOTHING;