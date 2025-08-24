-- Create postal_codes table to store unique postal codes with coordinates
CREATE TABLE postal_codes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    postal_code VARCHAR(10) NOT NULL UNIQUE,
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX idx_postal_codes_postal_code ON postal_codes(postal_code);
CREATE INDEX idx_postal_codes_geolocation ON postal_codes(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;