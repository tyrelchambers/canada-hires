-- Remove tfw_ratio column and index
DROP INDEX IF EXISTS idx_reports_tfw_ratio;
ALTER TABLE reports DROP COLUMN IF EXISTS tfw_ratio;