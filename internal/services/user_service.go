package services

import (
	"context"
	"time"

	"chat_app/internal/models"
	"chat_app/internal/repositories"
	"chat_app/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if username, ok := updates["username"].(string); ok && username != "" {
		// Check if username is already taken
		existingUser, err := s.userRepo.GetByUsername(ctx, username)
		if err == nil && existingUser.ID != userID {
			return nil, errors.NewConflictError("username already taken", nil)
		}
		user.Username = username
	}

	if email, ok := updates["email"].(string); ok && email != "" {
		// Check if email is already taken
		existingUser, err := s.userRepo.GetByEmail(ctx, email)
		if err == nil && existingUser.ID != userID {
			return nil, errors.NewConflictError("email already taken", nil)
		}
		user.Email = email
	}

	user.UpdatedAt = time.Now()

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	// Remove password from response
	user.Password = ""
	return user, nil
}

func (s *userService) DeactivateAccount(ctx context.Context, userID int) error {
	return s.userRepo.Delete(ctx, userID)
}

func (s *userService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	// Get user with password
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.NewUnauthorizedError("invalid current password", err)
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewInternalError("failed to hash new password", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s *userService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	// This would need to be implemented in the repository
	// For now, return empty slice
	return []*models.User{}, nil
}
