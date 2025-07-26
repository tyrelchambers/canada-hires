ALTER TABLE lmia_employers ADD COLUMN year INTEGER NOT NULL DEFAULT 0;
CREATE INDEX idx_lmia_employers_year ON lmia_employers(year);