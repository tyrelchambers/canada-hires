CREATE TABLE non_compliant_employers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_operating_name TEXT NOT NULL,
    business_legal_name TEXT,
    address TEXT,
    date_of_final_decision DATE,
    penalty_amount INTEGER,
    penalty_currency VARCHAR(10) DEFAULT 'CAD',
    status VARCHAR(50),
    scraped_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint to prevent duplicate entries (same business name and decision date)
    UNIQUE(business_operating_name, date_of_final_decision)
);

-- Junction table for many-to-many relationship between employers and reasons
CREATE TABLE non_compliant_employer_reasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employer_id UUID NOT NULL REFERENCES non_compliant_employers(id) ON DELETE CASCADE,
    reason_id INTEGER NOT NULL REFERENCES non_compliant_reasons(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(employer_id, reason_id)
);

-- Indexes for better performance
CREATE INDEX idx_non_compliant_employers_business_name ON non_compliant_employers(business_operating_name);
CREATE INDEX idx_non_compliant_employers_decision_date ON non_compliant_employers(date_of_final_decision);
CREATE INDEX idx_non_compliant_employers_status ON non_compliant_employers(status);
CREATE INDEX idx_non_compliant_employers_scraped_at ON non_compliant_employers(scraped_at);
CREATE INDEX idx_non_compliant_employer_reasons_employer_id ON non_compliant_employer_reasons(employer_id);
CREATE INDEX idx_non_compliant_employer_reasons_reason_id ON non_compliant_employer_reasons(reason_id);

-- Create updated_at trigger for employers table
CREATE OR REPLACE FUNCTION update_non_compliant_employers_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_non_compliant_employers_updated_at
    BEFORE UPDATE ON non_compliant_employers
    FOR EACH ROW
    EXECUTE FUNCTION update_non_compliant_employers_updated_at();