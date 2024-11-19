package chat

import (
	"database/sql"
	"fmt"
	"time"
)

// Message represents a single chat message in the application.
type Message struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`      // Unique identifier for the message
	Sender    string    `json:"sender"`    // The sender's name or ID
	Content   string    `json:"content"`   // The text of the message
	Timestamp time.Time `json:"timestamp"` // When the message was sent
	Room      string    `json:"room"`      // The chat room the message belongs to
}

// NewMessage creates a new message instance.
func NewMessage(msgType, sender, content, room string) *Message {
	return &Message{
		Type:      msgType,
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
		Room:      room,
	}
}

// Format formats the message for display or logging.
func (m *Message) Format() string {
	return "[" + m.Timestamp.Format("15:04:05") + "] " + m.Sender + ": " + m.Content
}

// Validate checks if the message is valid (e.g., not empty).
func (m *Message) Validate() error {
	if m.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}
	return nil
}

func SaveMessage(db *sql.DB, msg *Message) error {
	query := "INSERT INTO messages (sender, content, timestamp, room) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, msg.Sender, msg.Content, msg.Timestamp, msg.Room)
	return err
}

func GetMessages(db *sql.DB, room string) ([]Message, error) {
	query := "SELECT id, sender, content, timestamp FROM messages WHERE room = ? ORDER BY timestamp"
	rows, err := db.Query(query, room)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &msg.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
