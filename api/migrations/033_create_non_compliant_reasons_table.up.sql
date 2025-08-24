CREATE TABLE non_compliant_reasons (
    id SERIAL PRIMARY KEY,
    reason_code VARCHAR(10) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert common reason codes that appear in the data
INSERT INTO non_compliant_reasons (reason_code, description) VALUES
('5', 'Reason code 5'),
('6', 'Reason code 6'),
('15', 'Reason code 15')
ON CONFLICT (reason_code) DO NOTHING;

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_non_compliant_reasons_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_non_compliant_reasons_updated_at
    BEFORE UPDATE ON non_compliant_reasons
    FOR EACH ROW
    EXECUTE FUNCTION update_non_compliant_reasons_updated_at();