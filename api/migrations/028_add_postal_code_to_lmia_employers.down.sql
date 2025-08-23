-- Remove postal_code column from lmia_employers table
DROP INDEX IF EXISTS idx_lmia_employers_postal_code;
ALTER TABLE lmia_employers DROP COLUMN IF EXISTS postal_code;