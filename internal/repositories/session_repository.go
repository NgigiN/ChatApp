package repositories

import (
	"context"
	"database/sql"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions (id, user_id, token, expires_at, created_at, is_active)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	now := time.Now()
	session.CreatedAt = now
	session.IsActive = true

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.Token, session.ExpiresAt, session.CreatedAt, session.IsActive)
	
	if err != nil {
		return errors.NewDatabaseError("failed to create session", err)
	}

	return nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (*models.UserSession, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_active
		FROM user_sessions 
		WHERE token = ? AND is_active = true AND expires_at > NOW()`
	
	session := &models.UserSession{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt, &session.IsActive)
	
	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("session not found or expired", err)
	}
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get session by token", err)
	}

	return session, nil
}

func (r *sessionRepository) GetByUserID(ctx context.Context, userID int) ([]*models.UserSession, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, is_active
		FROM user_sessions 
		WHERE user_id = ? AND is_active = true AND expires_at > NOW()
		ORDER BY created_at DESC`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to get sessions by user ID", err)
	}
	defer rows.Close()

	var sessions []*models.UserSession
	for rows.Next() {
		session := &models.UserSession{}
		err := rows.Scan(&session.ID, &session.UserID, &session.Token, &session.ExpiresAt, &session.CreatedAt, &session.IsActive)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to scan session", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *models.UserSession) error {
	query := `
		UPDATE user_sessions 
		SET expires_at = ?, is_active = ?
		WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, session.ExpiresAt, session.IsActive, session.ID)
	if err != nil {
		return errors.NewDatabaseError("failed to update session", err)
	}

	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	query := `UPDATE user_sessions SET is_active = false WHERE token = ?`
	
	_, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return errors.NewDatabaseError("failed to delete session", err)
	}

	return nil
}

func (r *sessionRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `UPDATE user_sessions SET is_active = false WHERE user_id = ?`
	
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return errors.NewDatabaseError("failed to delete sessions by user ID", err)
	}

	return nil
}

func (r *sessionRepository) CleanupExpired(ctx context.Context) error {
	query := `UPDATE user_sessions SET is_active = false WHERE expires_at < NOW()`
	
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return errors.NewDatabaseError("failed to cleanup expired sessions", err)
	}

	return nil
}
