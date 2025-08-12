-- This migration removes all LMIA statistics data
-- Use with caution as this will delete all aggregated statistics

DELETE FROM lmia_job_statistics;