-- Add postal_code column to lmia_employers table
ALTER TABLE lmia_employers ADD COLUMN postal_code VARCHAR(10);

-- Create index for efficient querying
CREATE INDEX idx_lmia_employers_postal_code ON lmia_employers (postal_code) WHERE postal_code IS NOT NULL;