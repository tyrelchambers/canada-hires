-- Add geolocation fields to lmia_employers table
ALTER TABLE lmia_employers 
ADD COLUMN latitude DECIMAL(10, 7),
ADD COLUMN longitude DECIMAL(10, 7),
ADD COLUMN geocoded_at TIMESTAMP WITH TIME ZONE;

-- Create index for geolocation queries
CREATE INDEX idx_lmia_employers_geolocation ON lmia_employers(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;