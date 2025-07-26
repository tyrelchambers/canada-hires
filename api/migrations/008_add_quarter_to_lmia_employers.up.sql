ALTER TABLE lmia_employers ADD COLUMN quarter VARCHAR(10) NOT NULL DEFAULT '';
CREATE INDEX idx_lmia_employers_quarter ON lmia_employers(quarter);
CREATE INDEX idx_lmia_employers_quarter_year ON lmia_employers(quarter, year);