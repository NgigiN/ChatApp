package repositories

import (
	"chat_app/internal/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, username, email string) (bool, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *models.UserSession) error
	GetByToken(ctx context.Context, token string) (*models.UserSession, error)
	GetByUserID(ctx context.Context, userID int) ([]*models.UserSession, error)
	Update(ctx context.Context, session *models.UserSession) error
	Delete(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID int) error
	CleanupExpired(ctx context.Context) error
}

type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) error
	GetByID(ctx context.Context, id int) (*models.Room, error)
	GetByName(ctx context.Context, name string) (*models.Room, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Room, error)
	GetByUserID(ctx context.Context, userID int) ([]*models.Room, error)
	Update(ctx context.Context, room *models.Room) error
	Delete(ctx context.Context, id int) error
	Exists(ctx context.Context, name string) (bool, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id int) (*models.Message, error)
	GetByRoomID(ctx context.Context, roomID int, limit, offset int) ([]*models.Message, error)
	GetByRoomName(ctx context.Context, roomName string, limit, offset int) ([]*models.Message, error)
	GetRecent(ctx context.Context, roomID int, limit int) ([]*models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id int) error
	CountByRoomID(ctx context.Context, roomID int) (int64, error)
}

type RoomMemberRepository interface {
	AddMember(ctx context.Context, member *models.RoomMember) error
	RemoveMember(ctx context.Context, roomID, userID int) error
	GetMembers(ctx context.Context, roomID int) ([]*models.RoomMember, error)
	GetRoomsByUserID(ctx context.Context, userID int) ([]*models.Room, error)
	IsMember(ctx context.Context, roomID, userID int) (bool, error)
	GetMemberCount(ctx context.Context, roomID int) (int64, error)
}
