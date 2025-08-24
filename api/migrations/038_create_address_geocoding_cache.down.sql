-- Remove address geocoding cache table
DROP INDEX IF EXISTS idx_address_geocoding_cache_geocoded_at;
DROP INDEX IF EXISTS idx_address_geocoding_cache_coordinates;
DROP INDEX IF EXISTS idx_address_geocoding_cache_normalized_address;
DROP TABLE IF EXISTS address_geocoding_cache;