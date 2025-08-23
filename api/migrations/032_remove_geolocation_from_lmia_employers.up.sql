-- Remove latitude and longitude columns from lmia_employers table
ALTER TABLE lmia_employers DROP COLUMN IF EXISTS latitude;
ALTER TABLE lmia_employers DROP COLUMN IF EXISTS longitude;