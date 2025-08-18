ALTER TABLE reports ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'flagged'));
ALTER TABLE reports ADD COLUMN moderated_by UUID REFERENCES users(id);
ALTER TABLE reports ADD COLUMN moderation_notes TEXT;

CREATE INDEX idx_reports_status ON reports(status);