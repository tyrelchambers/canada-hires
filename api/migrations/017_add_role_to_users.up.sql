-- Add role field to users table

ALTER TABLE users 
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user';

-- Add constraint to ensure valid roles
ALTER TABLE users 
ADD CONSTRAINT chk_user_role 
CHECK (role IN ('user', 'admin'));

-- Add index for role filtering
CREATE INDEX idx_users_role ON users(role);

-- Comment on the new column
COMMENT ON COLUMN users.role IS 'User role: user or admin';