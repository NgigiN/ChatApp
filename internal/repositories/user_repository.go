package repositories

import (
	"context"
	"database/sql"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, updated_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.IsActive)
	
	if err != nil {
		return errors.NewDatabaseError("failed to create user", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return errors.NewDatabaseError("failed to get user ID", err)
	}

	user.ID = int(id)
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, is_active
		FROM users WHERE id = ? AND is_active = true`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("user not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get user by ID", err)
	}

	return user, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, is_active
		FROM users WHERE username = ? AND is_active = true`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("user not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get user by username", err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password, created_at, updated_at, is_active
		FROM users WHERE email = ? AND is_active = true`
	
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("user not found", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get user by email", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, password = ?, updated_at = ?, is_active = ?
		WHERE id = ?`
	
	user.UpdatedAt = time.Now()
	
	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.Password, user.UpdatedAt, user.IsActive, user.ID)
	
	if err != nil {
		return errors.NewDatabaseError("failed to update user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("user not found", nil)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE users SET is_active = false, updated_at = ? WHERE id = ?`
	
	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return errors.NewDatabaseError("failed to delete user", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError("failed to get rows affected", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("user not found", nil)
	}

	return nil
}

func (r *userRepository) Exists(ctx context.Context, username, email string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM users 
		WHERE (username = ? OR email = ?) AND is_active = true`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, username, email).Scan(&count)
	if err != nil {
		return false, errors.NewDatabaseError("failed to check user existence", err)
	}

	return count > 0, nil
}
