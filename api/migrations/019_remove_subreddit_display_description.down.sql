-- Add back display_name and description columns to subreddits table
ALTER TABLE subreddits ADD COLUMN display_name VARCHAR(100);
ALTER TABLE subreddits ADD COLUMN description TEXT;