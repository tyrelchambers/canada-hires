-- Recreate the junction table
CREATE TABLE non_compliant_employer_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employer_id UUID NOT NULL REFERENCES non_compliant_employers(id) ON DELETE CASCADE,
    reason_id INTEGER NOT NULL REFERENCES non_compliant_reasons(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employer_id, reason_id)
);

-- Recreate indexes
CREATE INDEX idx_non_compliant_employer_reasons_employer_id ON non_compliant_employer_reasons(employer_id);
CREATE INDEX idx_non_compliant_employer_reasons_reason_id ON non_compliant_employer_reasons(reason_id);

-- Remove the reason_codes column from non_compliant_employers
ALTER TABLE non_compliant_employers DROP COLUMN IF EXISTS reason_codes;