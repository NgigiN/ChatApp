package services

import (
	"chat_app/internal/models"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}

type UserService interface {
	GetProfile(ctx context.Context, userID int) (*models.User, error)
	UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) (*models.User, error)
	DeactivateAccount(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) error
}

type RoomService interface {
	CreateRoom(ctx context.Context, userID int, name, description string, isPrivate bool) (*models.Room, error)
	GetRoom(ctx context.Context, roomID int) (*models.Room, error)
	GetRoomByName(ctx context.Context, name string) (*models.Room, error)
	GetRooms(ctx context.Context, limit, offset int) ([]*models.Room, error)
	GetUserRooms(ctx context.Context, userID int) ([]*models.Room, error)
	UpdateRoom(ctx context.Context, roomID int, userID int, updates map[string]interface{}) (*models.Room, error)
	DeleteRoom(ctx context.Context, roomID int, userID int) error
	JoinRoom(ctx context.Context, roomID, userID int) error
	LeaveRoom(ctx context.Context, roomID, userID int) error
	GetRoomMembers(ctx context.Context, roomID int) ([]*models.User, error)
}

type MessageService interface {
	SendMessage(ctx context.Context, userID int, req *models.SendMessageRequest) (*models.Message, error)
	GetMessages(ctx context.Context, roomID int, limit, offset int) ([]*models.Message, error)
	GetRecentMessages(ctx context.Context, roomID int, limit int) ([]*models.Message, error)
	EditMessage(ctx context.Context, messageID, userID int, content string) (*models.Message, error)
	DeleteMessage(ctx context.Context, messageID, userID int) error
	GetMessage(ctx context.Context, messageID int) (*models.Message, error)
}

type WebSocketService interface {
	HandleConnection(ctx context.Context, conn interface{}, user *models.User) error
	JoinRoom(ctx context.Context, userID int, roomName string) error
	LeaveRoom(ctx context.Context, userID int, roomName string) error
	BroadcastMessage(ctx context.Context, message *models.Message) error
	BroadcastToRoom(ctx context.Context, roomName string, message *models.WebSocketMessage) error
	GetConnectedUsers(ctx context.Context, roomName string) ([]*models.User, error)
}
