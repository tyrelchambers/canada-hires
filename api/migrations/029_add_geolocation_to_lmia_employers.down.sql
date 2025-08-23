-- Remove geolocation fields from lmia_employers table
DROP INDEX IF EXISTS idx_lmia_employers_geolocation;

ALTER TABLE lmia_employers 
DROP COLUMN IF EXISTS latitude,
DROP COLUMN IF EXISTS longitude,
DROP COLUMN IF EXISTS geocoded_at;