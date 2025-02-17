-- Drop indexes
DROP INDEX IF EXISTS idx_messages_room_id_timestamp;
DROP INDEX IF EXISTS idx_messages_sender;
DROP INDEX IF EXISTS idx_user_room_activity_last_activity;
DROP INDEX IF EXISTS idx_users_email;

-- Drop triggers
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_rooms_updated_at ON rooms;
DROP TRIGGER IF EXISTS message_update_user_activity ON messages;

-- Drop trigger functions
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
DROP FUNCTION IF EXISTS update_user_room_activity() CASCADE;

-- Drop views
DROP VIEW IF EXISTS active_room_users;
DROP VIEW IF EXISTS recent_messages;

-- Drop tables
DROP TABLE IF EXISTS user_room_activity;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS rooms;
DROP TABLE IF EXISTS users;

-- Drop type
DROP TYPE IF EXISTS message_type;

-- Disable extensions (optional, only if you want to remove the extensions)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
-- DROP EXTENSION IF EXISTS "pgcrypto";