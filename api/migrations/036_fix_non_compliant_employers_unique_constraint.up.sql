-- Drop the existing unique constraint
ALTER TABLE non_compliant_employers DROP CONSTRAINT IF EXISTS non_compliant_employers_business_operating_name_date_of_final_dec_key;

-- Create a new unique constraint that properly handles NULL values
-- This uses a unique index with COALESCE to handle NULL dates
CREATE UNIQUE INDEX non_compliant_employers_unique_idx 
ON non_compliant_employers (business_operating_name, COALESCE(date_of_final_decision, '1900-01-01'::date));