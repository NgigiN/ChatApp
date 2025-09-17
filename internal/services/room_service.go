package services

import (
	"context"
	"time"

	"chat_app/internal/models"
	"chat_app/internal/repositories"
	"chat_app/pkg/errors"
)

type roomService struct {
	roomRepo       repositories.RoomRepository
	roomMemberRepo repositories.RoomMemberRepository
}

func NewRoomService(roomRepo repositories.RoomRepository, roomMemberRepo repositories.RoomMemberRepository) RoomService {
	return &roomService{
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
	}
}

func (s *roomService) CreateRoom(ctx context.Context, userID int, name, description string, isPrivate bool) (*models.Room, error) {
	// Check if room name already exists
	exists, err := s.roomRepo.Exists(ctx, name)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to check room existence", err)
	}
	if exists {
		return nil, errors.NewConflictError("room name already exists", nil)
	}

	// Create room
	room := &models.Room{
		Name:        name,
		Description: description,
		IsPrivate:   isPrivate,
		CreatedBy:   userID,
	}

	if err := s.roomRepo.Create(ctx, room); err != nil {
		return nil, err
	}

	// Add creator as member
	member := &models.RoomMember{
		RoomID: room.ID,
		UserID: userID,
	}

	if err := s.roomMemberRepo.AddMember(ctx, member); err != nil {
		// If adding member fails, clean up the room
		s.roomRepo.Delete(ctx, room.ID)
		return nil, err
	}

	return room, nil
}

func (s *roomService) GetRoom(ctx context.Context, roomID int) (*models.Room, error) {
	return s.roomRepo.GetByID(ctx, roomID)
}

func (s *roomService) GetRoomByName(ctx context.Context, name string) (*models.Room, error) {
	return s.roomRepo.GetByName(ctx, name)
}

func (s *roomService) GetRooms(ctx context.Context, limit, offset int) ([]*models.Room, error) {
	return s.roomRepo.GetAll(ctx, limit, offset)
}

func (s *roomService) GetUserRooms(ctx context.Context, userID int) ([]*models.Room, error) {
	return s.roomMemberRepo.GetRoomsByUserID(ctx, userID)
}

func (s *roomService) UpdateRoom(ctx context.Context, roomID int, userID int, updates map[string]interface{}) (*models.Room, error) {
	// Get room
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator
	if room.CreatedBy != userID {
		return nil, errors.NewForbiddenError("only room creator can update room", nil)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok && name != "" {
		// Check if new name is already taken
		exists, err := s.roomRepo.Exists(ctx, name)
		if err != nil {
			return nil, errors.NewDatabaseError("failed to check room name", err)
		}
		if exists {
			return nil, errors.NewConflictError("room name already exists", nil)
		}
		room.Name = name
	}

	if description, ok := updates["description"].(string); ok {
		room.Description = description
	}

	if isPrivate, ok := updates["is_private"].(bool); ok {
		room.IsPrivate = isPrivate
	}

	room.UpdatedAt = time.Now()

	// Update room
	if err := s.roomRepo.Update(ctx, room); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *roomService) DeleteRoom(ctx context.Context, roomID int, userID int) error {
	// Get room
	room, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if room.CreatedBy != userID {
		return errors.NewForbiddenError("only room creator can delete room", nil)
	}

	return s.roomRepo.Delete(ctx, roomID)
}

func (s *roomService) JoinRoom(ctx context.Context, roomID, userID int) error {
	// Check if room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return err
	}

	// Check if user is already a member
	isMember, err := s.roomMemberRepo.IsMember(ctx, roomID, userID)
	if err != nil {
		return errors.NewDatabaseError("failed to check membership", err)
	}
	if isMember {
		return errors.NewConflictError("user is already a member", nil)
	}

	// Add member
	member := &models.RoomMember{
		RoomID: roomID,
		UserID: userID,
	}

	return s.roomMemberRepo.AddMember(ctx, member)
}

func (s *roomService) LeaveRoom(ctx context.Context, roomID, userID int) error {
	// Check if user is a member
	isMember, err := s.roomMemberRepo.IsMember(ctx, roomID, userID)
	if err != nil {
		return errors.NewDatabaseError("failed to check membership", err)
	}
	if !isMember {
		return errors.NewNotFoundError("user is not a member of this room", nil)
	}

	return s.roomMemberRepo.RemoveMember(ctx, roomID, userID)
}

func (s *roomService) GetRoomMembers(ctx context.Context, roomID int) ([]*models.User, error) {
	// This would require a join query - simplified for now
	// In a real implementation, you'd need to join with users table
	return []*models.User{}, nil
}
