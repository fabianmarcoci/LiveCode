-- Remove default value from updated_at
ALTER TABLE users ALTER COLUMN updated_at DROP DEFAULT;

-- Drop trigger
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();