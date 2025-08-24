-- Drop the new unique index
DROP INDEX IF EXISTS non_compliant_employers_unique_idx;

-- Recreate the original unique constraint
ALTER TABLE non_compliant_employers ADD CONSTRAINT non_compliant_employers_business_operating_name_date_of_final_dec_key 
UNIQUE (business_operating_name, date_of_final_decision);