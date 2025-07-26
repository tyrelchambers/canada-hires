DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
ALTER TABLE IF EXISTS sessions DROP CONSTRAINT IF EXISTS fk_sessions_user_id;
DROP TABLE IF EXISTS sessions;