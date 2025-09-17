package repositories

import (
	"context"
	"database/sql"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
)

type roomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) Create(ctx context.Context, room *models.Room) error {
	query := `
		INSERT INTO rooms (name, description, is_private, created_by, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	room.CreatedAt = now
	room.UpdatedAt = now
	room.IsActive = true

	result, err := r.db.ExecContext(ctx, query,
		room.Name, room.Description, room.IsPrivate, room.CreatedBy, room.CreatedAt, room.UpdatedAt, room.IsActive)

	if err != nil {
		return errors.NewDatabaseError("failed to create room", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewDatabaseError("failed to get room ID", err)
	}

	room.ID = int(id)
	return nil
}

func (r *roomRepository) GetByID(ctx context.Context, id int) (*models.Room, error) {
	query := `
		SELECT id, name, description, is_private, created_by, created_at, updated_at, is_active
		FROM rooms WHERE id = ? AND is_active = true`

	room := &models.Room{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID, &room.Name, &room.Description, &room.IsPrivate,
		&room.CreatedBy, &room.CreatedAt, &room.UpdatedAt, &room.IsActive)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("room not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get room by ID", err)
	}

	return room, nil
}

func (r *roomRepository) GetByName(ctx context.Context, name string) (*models.Room, error) {
	query := `
		SELECT id, name, description, is_private, created_by, created_at, updated_at, is_active
		FROM rooms WHERE name = ? AND is_active = true`

	room := &models.Room{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&room.ID, &room.Name, &room.Description, &room.IsPrivate,
		&room.CreatedBy, &room.CreatedAt, &room.UpdatedAt, &room.IsActive)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("room not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get room by name", err)
	}

	return room, nil
}

func (r *roomRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Room, error) {
	query := `
		SELECT id, name, description, is_private, created_by, created_at, updated_at, is_active
		FROM rooms
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get rooms", err)
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

func (r *roomRepository) GetByUserID(ctx context.Context, userID int) ([]*models.Room, error) {
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

func (r *roomRepository) Update(ctx context.Context, room *models.Room) error {
	query := `
		UPDATE rooms
		SET name = ?, description = ?, is_private = ?, updated_at = ?, is_active = ?
		WHERE id = ?`

	room.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		room.Name, room.Description, room.IsPrivate, room.UpdatedAt, room.IsActive, room.ID)

	if err != nil {
		return errors.NewDatabaseError("failed to update room", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("room not found", nil)
	}

	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE rooms SET is_active = false, updated_at = ? WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return errors.NewDatabaseError("failed to delete room", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("room not found", nil)
	}

	return nil
}

func (r *roomRepository) Exists(ctx context.Context, name string) (bool, error) {
	query := `SELECT COUNT(*) FROM rooms WHERE name = ? AND is_active = true`

	var count int
	err := r.db.QueryRowContext(ctx, query, name).Scan(&count)
	if err != nil {
		return false, errors.NewDatabaseError("failed to check room existence", err)
	}

	return count > 0, nil
}
