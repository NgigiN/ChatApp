package chat

import (
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Room      string    `json:"room"`
}

func NewMessage(msgType, sender, content, room string) *Message {
	return &Message{
		Type:      msgType,
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
		Room:      room,
	}
}

func (m *Message) Format() string {
	return fmt.Sprintf("[%s] %s: %s", m.Timestamp.Format("15:04:05"), m.Sender, m.Content)
}

func (m *Message) Validate() error {
	if m.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}
	return nil
}

func SaveMessage(db *sql.DB, msg *Message) error {
	query := "INSERT INTO messages (room, sender, content, timestamp) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, msg.Room, msg.Sender, msg.Content, msg.Timestamp)
	return err
}

func GetMessages(db *sql.DB, room string) ([]*Message, error) {
	query := "SELECT id, sender, content, timestamp, room FROM messages WHERE room = ? ORDER BY timestamp DESC LIMIT 50"
	rows, err := db.Query(query, room)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		var timestamp []uint8
		if err := rows.Scan(&msg.ID, &msg.Sender, &msg.Content, &timestamp, &msg.Room); err != nil {
			return nil, err
		}

		msg.Timestamp, err = time.Parse("2006-01-02 15:04:05", string(timestamp))
		if err != nil {
			return nil, err
		}

		messages = append(messages, msg)
	}

	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
