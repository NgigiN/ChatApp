package services

import (
	"context"
	"time"

	"chat_app/internal/models"
	"chat_app/internal/repositories"
	"chat_app/pkg/errors"
)

type messageService struct {
	messageRepo    repositories.MessageRepository
	roomRepo       repositories.RoomRepository
	roomMemberRepo repositories.RoomMemberRepository
}

func NewMessageService(messageRepo repositories.MessageRepository, roomRepo repositories.RoomRepository, roomMemberRepo repositories.RoomMemberRepository) MessageService {
	return &messageService{
		messageRepo:    messageRepo,
		roomRepo:       roomRepo,
		roomMemberRepo: roomMemberRepo,
	}
}

func (s *messageService) SendMessage(ctx context.Context, userID int, req *models.SendMessageRequest) (*models.Message, error) {
	// Get room by name
	room, err := s.roomRepo.GetByName(ctx, req.Room)
	if err != nil {
		return nil, err
	}

	// Check if user is a member of the room
	isMember, err := s.roomMemberRepo.IsMember(ctx, room.ID, userID)
	if err != nil {
		return nil, errors.NewDatabaseError("failed to check room membership", err)
	}
	if !isMember {
		return nil, errors.NewForbiddenError("user is not a member of this room", nil)
	}

	// Get user for username (simplified - in production, get from context or user service)
	// For now, we'll use a placeholder
	username := "user" // This should come from user service or context

	// Create message
	message := &models.Message{
		RoomID:   room.ID,
		UserID:   userID,
		Username: username,
		Content:  req.Content,
		Type:     req.Type,
	}

	if message.Type == "" {
		message.Type = "message"
	}

	// Save message
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *messageService) GetMessages(ctx context.Context, roomID int, limit, offset int) ([]*models.Message, error) {
	// Check if room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.GetByRoomID(ctx, roomID, limit, offset)
}

func (s *messageService) GetRecentMessages(ctx context.Context, roomID int, limit int) ([]*models.Message, error) {
	// Check if room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	return s.messageRepo.GetRecent(ctx, roomID, limit)
}

func (s *messageService) EditMessage(ctx context.Context, messageID, userID int, content string) (*models.Message, error) {
	// Get message
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	// Check if user is the message author
	if message.UserID != userID {
		return nil, errors.NewForbiddenError("only message author can edit message", nil)
	}

	// Update message
	message.Content = content
	message.UpdatedAt = time.Now()

	if err := s.messageRepo.Update(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *messageService) DeleteMessage(ctx context.Context, messageID, userID int) error {
	// Get message
	message, err := s.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}

	// Check if user is the message author
	if message.UserID != userID {
		return errors.NewForbiddenError("only message author can delete message", nil)
	}

	return s.messageRepo.Delete(ctx, messageID)
}

func (s *messageService) GetMessage(ctx context.Context, messageID int) (*models.Message, error) {
	return s.messageRepo.GetByID(ctx, messageID)
}
