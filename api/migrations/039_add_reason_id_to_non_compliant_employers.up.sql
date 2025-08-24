-- Add reason_codes array column to non_compliant_employers table
ALTER TABLE non_compliant_employers 
ADD COLUMN reason_codes VARCHAR(10)[] DEFAULT '{}';

-- Create index for better performance on reason_codes array queries
CREATE INDEX idx_non_compliant_employers_reason_codes ON non_compliant_employers USING GIN(reason_codes);

-- Drop the junction table as it's no longer needed
DROP TABLE IF EXISTS non_compliant_employer_reasons;