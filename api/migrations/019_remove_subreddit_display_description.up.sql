-- Remove display_name and description columns from subreddits table
ALTER TABLE subreddits DROP COLUMN IF EXISTS display_name;
ALTER TABLE subreddits DROP COLUMN IF EXISTS description;