-- Create address geocoding cache table to store address -> coordinates mappings
CREATE TABLE address_geocoding_cache (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    normalized_address TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    confidence DOUBLE PRECISION,
    geocoded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create unique index on normalized address to prevent duplicates
CREATE UNIQUE INDEX idx_address_geocoding_cache_normalized_address ON address_geocoding_cache (normalized_address);

-- Create index on coordinates for spatial queries
CREATE INDEX idx_address_geocoding_cache_coordinates ON address_geocoding_cache (latitude, longitude);

-- Create index on geocoded_at for cleanup/maintenance
CREATE INDEX idx_address_geocoding_cache_geocoded_at ON address_geocoding_cache (geocoded_at);