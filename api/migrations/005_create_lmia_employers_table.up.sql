CREATE TABLE lmia_employers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_id UUID NOT NULL REFERENCES lmia_resources(id) ON DELETE CASCADE,
    
    -- These are the ONLY 8 columns from the actual LMIA CSV files
    province_territory TEXT,           -- "Province/Territory"
    program_stream TEXT,              -- "Program Stream" 
    employer TEXT NOT NULL,           -- "Employer"
    address TEXT,                     -- "Address"
    occupation TEXT,                  -- "Occupation"
    incorporate_status TEXT,          -- "Incorporate Status"
    approved_lmias INTEGER,           -- "Approved LMIAs"
    approved_positions INTEGER,       -- "Approved Positions"
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_lmia_employers_resource_id ON lmia_employers(resource_id);
CREATE INDEX idx_lmia_employers_employer ON lmia_employers(employer);
CREATE INDEX idx_lmia_employers_province_territory ON lmia_employers(province_territory);
CREATE INDEX idx_lmia_employers_program_stream ON lmia_employers(program_stream);