package handlers

import (
	"context"
	"chat_app/internal/models"
)

// Mock services for development - replace with real implementations
type mockRoomService struct{}

func (m *mockRoomService) CreateRoom(ctx context.Context, userID int, name, description string, isPrivate bool) (*models.Room, error) {
	return &models.Room{
		ID:          1,
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		CreatedBy:   userID,
	}, nil
}

func (m *mockRoomService) GetRoom(ctx context.Context, roomID int) (*models.Room, error) {
	return &models.Room{
		ID:          roomID,
		Name:        "Test Room",
		Description: "Test Description",
		IsPrivate:   false,
		CreatedBy:   1,
	}, nil
}

func (m *mockRoomService) GetRoomByName(ctx context.Context, name string) (*models.Room, error) {
	return &models.Room{
		ID:          1,
		Name:        name,
		Description: "Test Description",
		IsPrivate:   false,
		CreatedBy:   1,
	}, nil
}

func (m *mockRoomService) GetRooms(ctx context.Context, limit, offset int) ([]*models.Room, error) {
	return []*models.Room{}, nil
}

func (m *mockRoomService) GetUserRooms(ctx context.Context, userID int) ([]*models.Room, error) {
	return []*models.Room{
		{
			ID:          1,
			Name:        "Math 101",
			Description: "Mathematics for beginners",
			IsPrivate:   false,
			CreatedBy:   userID,
		},
		{
			ID:          2,
			Name:        "Physics Lab",
			Description: "Advanced Physics Laboratory",
			IsPrivate:   true,
			CreatedBy:   userID,
		},
	}, nil
}

func (m *mockRoomService) UpdateRoom(ctx context.Context, roomID int, userID int, updates map[string]interface{}) (*models.Room, error) {
	return &models.Room{
		ID:          roomID,
		Name:        "Updated Room",
		Description: "Updated Description",
		IsPrivate:   false,
		CreatedBy:   userID,
	}, nil
}

func (m *mockRoomService) DeleteRoom(ctx context.Context, roomID int, userID int) error {
	return nil
}

func (m *mockRoomService) JoinRoom(ctx context.Context, roomID, userID int) error {
	return nil
}

func (m *mockRoomService) LeaveRoom(ctx context.Context, roomID, userID int) error {
	return nil
}

func (m *mockRoomService) GetRoomMembers(ctx context.Context, roomID int) ([]*models.User, error) {
	return []*models.User{
		{
			ID:       1,
			Username: "john_doe",
			Email:    "john@example.com",
		},
		{
			ID:       2,
			Username: "jane_smith",
			Email:    "jane@example.com",
		},
	}, nil
}

type mockUserService struct{}

func (m *mockUserService) GetProfile(ctx context.Context, userID int) (*models.User, error) {
	return &models.User{
		ID:       userID,
		Username: "test_user",
		Email:    "test@example.com",
	}, nil
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) (*models.User, error) {
	return &models.User{
		ID:       userID,
		Username: "updated_user",
		Email:    "updated@example.com",
	}, nil
}

func (m *mockUserService) DeactivateAccount(ctx context.Context, userID int) error {
	return nil
}

func (m *mockUserService) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error {
	return nil
}

func (m *mockUserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return &models.User{
		ID:       1,
		Username: username,
		Email:    username + "@example.com",
	}, nil
}

func (m *mockUserService) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	return []*models.User{
		{
			ID:       1,
			Username: "user1",
			Email:    "user1@example.com",
		},
		{
			ID:       2,
			Username: "user2",
			Email:    "user2@example.com",
		},
	}, nil
}
