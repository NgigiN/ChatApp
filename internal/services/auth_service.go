package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"chat_app/internal/models"
	"chat_app/internal/repositories"
	"chat_app/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
	jwtSecret   string
	jwtExpiry   time.Duration
}

func NewAuthService(userRepo repositories.UserRepository, sessionRepo repositories.SessionRepository, jwtSecret string, jwtExpiry time.Duration) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
		jwtExpiry:   jwtExpiry,
	}
}

func (s *authService) Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	exists, err := s.userRepo.Exists(ctx, req.Username, req.Email)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to check user existence", err)
	}
	if exists {
		return nil, errors.NewConflictError("user already exists", nil)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError("failed to hash password", err)
	}

	// Create user
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &models.UserSession{
		ID:        generateSessionID(),
		UserID:    user.ID,
		Token:     accessToken,
		ExpiresAt: time.Now().Add(s.jwtExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.NewUnauthorizedError("invalid credentials", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Create session
	session := &models.UserSession{
		ID:        generateSessionID(),
		UserID:    user.ID,
		Token:     accessToken,
		ExpiresAt: time.Now().Add(s.jwtExpiry),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token (simplified - in production, use proper JWT validation)
	session, err := s.sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid refresh token", err)
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Update session
	session.Token = accessToken
	session.ExpiresAt = time.Now().Add(s.jwtExpiry)
	if err := s.sessionRepo.Update(ctx, session); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*models.User, error) {
	session, err := s.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, errors.NewUnauthorizedError("invalid token", err)
	}

	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) generateTokens(userID int) (string, string, error) {
	// Simplified token generation - in production, use proper JWT
	accessToken := fmt.Sprintf("access_%d_%d", userID, time.Now().Unix())
	refreshToken := fmt.Sprintf("refresh_%d_%d", userID, time.Now().Unix())
	return accessToken, refreshToken, nil
}

func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
