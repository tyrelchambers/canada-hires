-- Add back latitude and longitude columns to lmia_employers table (rollback)
ALTER TABLE lmia_employers ADD COLUMN IF NOT EXISTS latitude DECIMAL(10, 7);
ALTER TABLE lmia_employers ADD COLUMN IF NOT EXISTS longitude DECIMAL(10, 7);