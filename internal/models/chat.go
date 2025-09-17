package models

import (
	"time"
)

type Room struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	IsPrivate   bool      `json:"is_private" db:"is_private"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	IsActive    bool      `json:"is_active" db:"is_active"`
}

type Message struct {
	ID        int       `json:"id" db:"id"`
	RoomID    int       `json:"room_id" db:"room_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type RoomMember struct {
	ID       int       `json:"id" db:"id"`
	RoomID   int       `json:"room_id" db:"room_id"`
	UserID   int       `json:"user_id" db:"user_id"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
	IsActive bool      `json:"is_active" db:"is_active"`
}

type WebSocketMessage struct {
	Type      string                 `json:"type"`
	Room      string                 `json:"room,omitempty"`
	Content   string                 `json:"content,omitempty"`
	Sender    string                 `json:"sender,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

type JoinRoomRequest struct {
	RoomName string `json:"room_name" validate:"required,min=1,max=100"`
}

type SendMessageRequest struct {
	Room    string `json:"room" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=1000"`
	Type    string `json:"type,omitempty"`
}
