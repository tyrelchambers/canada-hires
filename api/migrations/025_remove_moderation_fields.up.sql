ALTER TABLE reports DROP COLUMN IF EXISTS status;
ALTER TABLE reports DROP COLUMN IF EXISTS moderated_by;
ALTER TABLE reports DROP COLUMN IF EXISTS moderation_notes;

DROP INDEX IF EXISTS idx_reports_status;