package services

import (
    "context"
    "fmt"
    "time"

    "chat_app/internal/config"
    "chat_app/internal/models"
    "chat_app/internal/repositories"
    "chat_app/pkg/errors"
    "chat_app/pkg/utils"
    "github.com/redis/go-redis/v9"
)

type messageService struct {
    messageRepo    repositories.MessageRepository
    roomRepo       repositories.RoomRepository
    roomMemberRepo repositories.RoomMemberRepository
    cache          *redis.Client
}

func NewMessageService(messageRepo repositories.MessageRepository, roomRepo repositories.RoomRepository, roomMemberRepo repositories.RoomMemberRepository) MessageService {
    cfg := config.Load()
    redisClient := config.NewRedisClient(cfg.Redis)
    return &messageService{
        messageRepo:    messageRepo,
        roomRepo:       roomRepo,
        roomMemberRepo: roomMemberRepo,
        cache:          redisClient,
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

    if s.cache != nil && offset == 0 {
        cacheKey := fmt.Sprintf("room:%d:messages:recent:%d", roomID, limit)
        if data, err := s.cache.Get(ctx, cacheKey).Bytes(); err == nil && len(data) > 0 {
            var msgs []*models.Message
            if err := utils.MustUnmarshal(data, &msgs); err == nil {
                return msgs, nil
            }
        }

        msgs, err := s.messageRepo.GetByRoomID(ctx, roomID, limit, offset)
        if err != nil {
            return nil, err
        }
        _ = s.cache.Set(ctx, cacheKey, utils.MustMarshal(msgs), 30*time.Second).Err()
        return msgs, nil
    }

    return s.messageRepo.GetByRoomID(ctx, roomID, limit, offset)
}

func (s *messageService) GetRecentMessages(ctx context.Context, roomID int, limit int) ([]*models.Message, error) {
    // Check if room exists
    _, err := s.roomRepo.GetByID(ctx, roomID)
    if err != nil {
        return nil, err
    }

    if s.cache != nil {
        cacheKey := fmt.Sprintf("room:%d:messages:recent:%d", roomID, limit)
        if data, err := s.cache.Get(ctx, cacheKey).Bytes(); err == nil && len(data) > 0 {
            var msgs []*models.Message
            if err := utils.MustUnmarshal(data, &msgs); err == nil {
                return msgs, nil
            }
        }

        msgs, err := s.messageRepo.GetRecent(ctx, roomID, limit)
        if err != nil {
            return nil, err
        }
        _ = s.cache.Set(ctx, cacheKey, utils.MustMarshal(msgs), 30*time.Second).Err()
        return msgs, nil
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
