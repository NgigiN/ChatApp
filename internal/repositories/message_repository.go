package repositories

import (
	"context"
	"database/sql"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	query := `
		INSERT INTO messages (room_id, user_id, username, content, type, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	result, err := r.db.ExecContext(ctx, query,
		message.RoomID, message.UserID, message.Username, message.Content, message.Type, message.CreatedAt, message.UpdatedAt)
	
	if err != nil {
		return errors.NewDatabaseError("failed to create message", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewDatabaseError("failed to get message ID", err)
	}

	message.ID = int(id)
	return nil
}

func (r *messageRepository) GetByID(ctx context.Context, id int) (*models.Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, type, created_at, updated_at
		FROM messages WHERE id = ?`
	
	message := &models.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID, &message.RoomID, &message.UserID, &message.Username,
		&message.Content, &message.Type, &message.CreatedAt, &message.UpdatedAt)
	
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("message not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get message by ID", err)
	}

	return message, nil
}

func (r *messageRepository) GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, type, created_at, updated_at
		FROM messages 
		WHERE room_id = ? 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`
	
	rows, err := r.db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get messages by room ID", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(&message.ID, &message.RoomID, &message.UserID, &message.Username,
			&message.Content, &message.Type, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan message", err)
		}
		messages = append(messages, message)
	}

	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) GetByRoomName(ctx context.Context, roomName string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT m.id, m.room_id, m.user_id, m.username, m.content, m.type, m.created_at, m.updated_at
		FROM messages m
		INNER JOIN rooms r ON m.room_id = r.id
		WHERE r.name = ? AND r.is_active = true
		ORDER BY m.created_at DESC 
		LIMIT ? OFFSET ?`
	
	rows, err := r.db.QueryContext(ctx, query, roomName, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get messages by room name", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(&message.ID, &message.RoomID, &message.UserID, &message.Username,
			&message.Content, &message.Type, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan message", err)
		}
		messages = append(messages, message)
	}

	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) GetRecent(ctx context.Context, roomID int, limit int) ([]*models.Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, type, created_at, updated_at
		FROM messages 
		WHERE room_id = ? 
		ORDER BY created_at DESC 
		LIMIT ?`
	
	rows, err := r.db.QueryContext(ctx, query, roomID, limit)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get recent messages", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		message := &models.Message{}
		err := rows.Scan(&message.ID, &message.RoomID, &message.UserID, &message.Username,
			&message.Content, &message.Type, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan message", err)
		}
		messages = append(messages, message)
	}

	// Reverse to show oldest first
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	query := `
		UPDATE messages 
		SET content = ?, updated_at = ?
		WHERE id = ?`
	
	message.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query, message.Content, message.UpdatedAt, message.ID)
	if err != nil {
		return errors.NewDatabaseError("failed to update message", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("message not found", nil)
	}

	return nil
}

func (r *messageRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM messages WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewDatabaseError("failed to delete message", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("message not found", nil)
	}

	return nil
}

func (r *messageRepository) CountByRoomID(ctx context.Context, roomID int) (int64, error) {
	query := `SELECT COUNT(*) FROM messages WHERE room_id = ?`
	
	var count int64
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&count)
	if err != nil {
		return 0, errors.NewDatabaseError("failed to count messages", err)
	}

	return count, nil
}
