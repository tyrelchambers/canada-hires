CREATE TABLE lmia_resources (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    resource_id VARCHAR(255) UNIQUE NOT NULL,
    name TEXT NOT NULL,
    quarter VARCHAR(10) NOT NULL,
    year INTEGER NOT NULL,
    url TEXT NOT NULL,
    format VARCHAR(10) NOT NULL,
    language VARCHAR(5) NOT NULL DEFAULT 'en',
    size_bytes BIGINT,
    last_modified TIMESTAMP WITH TIME ZONE,
    date_published TIMESTAMP WITH TIME ZONE,
    downloaded_at TIMESTAMP WITH TIME ZONE,
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_lmia_resources_quarter_year ON lmia_resources(quarter, year);
CREATE INDEX idx_lmia_resources_language ON lmia_resources(language);
CREATE INDEX idx_lmia_resources_processed_at ON lmia_resources(processed_at);