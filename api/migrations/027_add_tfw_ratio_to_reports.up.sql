-- Add tfw_ratio column to reports table
ALTER TABLE reports ADD COLUMN tfw_ratio VARCHAR(10) CHECK (tfw_ratio IN ('few', 'many', 'most', 'all'));

-- Migrate existing confidence_level data to tfw_ratio
UPDATE reports 
SET tfw_ratio = CASE 
    WHEN confidence_level IS NULL THEN NULL
    WHEN confidence_level BETWEEN 1 AND 3 THEN 'few'
    WHEN confidence_level BETWEEN 4 AND 6 THEN 'many'
    WHEN confidence_level BETWEEN 7 AND 9 THEN 'most'
    WHEN confidence_level = 10 THEN 'all'
    ELSE 'many' -- fallback for any edge cases
END
WHERE confidence_level IS NOT NULL;

-- Create index on new column for performance
CREATE INDEX idx_reports_tfw_ratio ON reports(tfw_ratio);