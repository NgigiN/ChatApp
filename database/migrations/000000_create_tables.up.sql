-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE message_type AS ENUM ('message', 'system', 'join');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE
);

CREATE TABLE rooms (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    type message_type DEFAULT 'message',
    sender VARCHAR(100) NOT NULL,
    room VARCHAR(100) NOT NULL,
    room_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_room FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    CONSTRAINT content_not_empty CHECK (length(trim(content)) > 0)
);

CREATE TABLE user_room_activity (
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    room_id BIGINT REFERENCES rooms(id) ON DELETE CASCADE,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, room_id)
);

-- Create a view for active users in rooms
CREATE VIEW active_room_users AS
SELECT
    r.name as room_name,
    u.first_name,
    u.role,
    ura.last_activity
FROM user_room_activity ura
JOIN users u ON u.id = ura.user_id
JOIN rooms r ON r.id = ura.room_id
WHERE ura.last_activity > (CURRENT_TIMESTAMP - INTERVAL '15 minutes');

-- Create a view for recent messages with user details
CREATE VIEW recent_messages AS
SELECT
    m.id,
    m.content,
    m.type,
    m.sender,
    m.room,
    m.created_at as timestamp,
    r.name as room_name
FROM messages m
JOIN rooms r ON r.id = m.room_id
WHERE m.created_at > (CURRENT_TIMESTAMP - INTERVAL '24 hours')
ORDER BY m.created_at DESC;
