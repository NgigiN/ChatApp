package repositories

import (
	"context"
	"database/sql"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
)

type roomMemberRepository struct {
	db *sql.DB
}

func NewRoomMemberRepository(db *sql.DB) RoomMemberRepository {
	return &roomMemberRepository{db: db}
}

func (r *roomMemberRepository) AddMember(ctx context.Context, member *models.RoomMember) error {
	query := `
		INSERT INTO room_members (room_id, user_id, joined_at, is_active)
		VALUES (?, ?, ?, ?)`

	now := time.Now()
	member.JoinedAt = now
	member.IsActive = true

	_, err := r.db.ExecContext(ctx, query, member.RoomID, member.UserID, member.JoinedAt, member.IsActive)
	if err != nil {
		return errors.NewDatabaseError("failed to add room member", err)
	}

	return nil
}

func (r *roomMemberRepository) RemoveMember(ctx context.Context, roomID, userID int) error {
	query := `UPDATE room_members SET is_active = false WHERE room_id = ? AND user_id = ?`

	result, err := r.db.ExecContext(ctx, query, roomID, userID)
	if err != nil {
		return errors.NewDatabaseError("failed to remove room member", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("room member not found", nil)
	}

	return nil
}

func (r *roomMemberRepository) GetMembers(ctx context.Context, roomID int) ([]*models.RoomMember, error) {
	query := `
		SELECT id, room_id, user_id, joined_at, is_active
		FROM room_members
		WHERE room_id = ? AND is_active = true
		ORDER BY joined_at ASC`

	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get room members", err)
	}
	defer rows.Close()

	var members []*models.RoomMember
	for rows.Next() {
		member := &models.RoomMember{}
		err := rows.Scan(&member.ID, &member.RoomID, &member.UserID, &member.JoinedAt, &member.IsActive)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan room member", err)
		}
		members = append(members, member)
	}

	return members, nil
}

func (r *roomMemberRepository) GetRoomsByUserID(ctx context.Context, userID int) ([]*models.Room, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_private, r.created_by, r.created_at, r.updated_at, r.is_active
		FROM rooms r
		INNER JOIN room_members rm ON r.id = rm.room_id
		WHERE rm.user_id = ? AND r.is_active = true AND rm.is_active = true
		ORDER BY r.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get user rooms", err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.ID, &room.Name, &room.Description, &room.IsPrivate,
			&room.CreatedBy, &room.CreatedAt, &room.UpdatedAt, &room.IsActive)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan room", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *roomMemberRepository) IsMember(ctx context.Context, roomID, userID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM room_members
		WHERE room_id = ? AND user_id = ? AND is_active = true`

	var count int
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan(&count)
	if err != nil {
		return false, errors.NewDatabaseError("failed to check room membership", err)
	}

	return count > 0, nil
}

func (r *roomMemberRepository) GetMemberCount(ctx context.Context, roomID int) (int64, error) {
	query := `SELECT COUNT(*) FROM room_members WHERE room_id = ? AND is_active = true`

	var count int64
	err := r.db.QueryRowContext(ctx, query, roomID).Scan(&count)
	if err != nil {
		return 0, errors.NewDatabaseError("failed to get member count", err)
	}

	return count, nil
}
