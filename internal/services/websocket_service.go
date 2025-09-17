package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"chat_app/internal/models"
	"chat_app/pkg/errors"
	"golang.org/x/net/websocket"
)

type WebSocketConnection struct {
	Conn     *websocket.Conn
	User     *models.User
	RoomName string
}

type websocketService struct {
	connections map[string]*WebSocketConnection // connection ID -> connection
	rooms       map[string]map[string]*WebSocketConnection // room name -> connection ID -> connection
	mu          sync.RWMutex
}

func NewWebSocketService() WebSocketService {
	return &websocketService{
		connections: make(map[string]*WebSocketConnection),
		rooms:       make(map[string]map[string]*WebSocketConnection),
	}
}

func (s *websocketService) HandleConnection(ctx context.Context, conn interface{}, user *models.User) error {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return errors.NewInvalidInputError("invalid connection type", nil)
	}

	connectionID := generateConnectionID()
	wsConnection := &WebSocketConnection{
		Conn: wsConn,
		User: user,
	}

	s.mu.Lock()
	s.connections[connectionID] = wsConnection
	s.mu.Unlock()

	// Handle connection cleanup on disconnect
	defer func() {
		s.mu.Lock()
		delete(s.connections, connectionID)
		if wsConnection.RoomName != "" {
			if roomConnections, exists := s.rooms[wsConnection.RoomName]; exists {
				delete(roomConnections, connectionID)
				if len(roomConnections) == 0 {
					delete(s.rooms, wsConnection.RoomName)
				}
			}
		}
		s.mu.Unlock()
	}()

	// Handle incoming messages
	for {
		var msg models.WebSocketMessage
		if err := websocket.JSON.Receive(wsConn, &msg); err != nil {
			log.Printf("WebSocket receive error: %v", err)
			break
		}

		// Process message based on type
		switch msg.Type {
		case "join":
			if roomName, ok := msg.Data["room"].(string); ok {
				if err := s.JoinRoom(ctx, user.ID, roomName); err != nil {
					log.Printf("Error joining room: %v", err)
				}
			}
		case "message":
			if content, ok := msg.Data["content"].(string); ok {
				message := &models.Message{
					RoomID:   0, // Will be set by message service
					UserID:   user.ID,
					Username: user.Username,
					Content:  content,
					Type:     "message",
				}
				if err := s.BroadcastMessage(ctx, message); err != nil {
					log.Printf("Error broadcasting message: %v", err)
				}
			}
		}
	}

	return nil
}

func (s *websocketService) JoinRoom(ctx context.Context, userID int, roomName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find user's connection
	var connection *WebSocketConnection
	for _, conn := range s.connections {
		if conn.User.ID == userID {
			connection = conn
			break
		}
	}

	if connection == nil {
		return errors.NewNotFoundError("user connection not found", nil)
	}

	// Remove from current room if any
	if connection.RoomName != "" {
		if roomConnections, exists := s.rooms[connection.RoomName]; exists {
			// Find connection ID
			for id, conn := range roomConnections {
				if conn == connection {
					delete(roomConnections, id)
					break
				}
			}
		}
	}

	// Add to new room
	connection.RoomName = roomName
	if s.rooms[roomName] == nil {
		s.rooms[roomName] = make(map[string]*WebSocketConnection)
	}
	// Find connection ID
	for id, conn := range s.connections {
		if conn == connection {
			s.rooms[roomName][id] = connection
			break
		}
	}

	// Notify user of successful join
	joinMsg := &models.WebSocketMessage{
		Type:    "system",
		Content: "You have joined the room: " + roomName,
	}
	websocket.JSON.Send(connection.Conn, joinMsg)

	return nil
}

func (s *websocketService) LeaveRoom(ctx context.Context, userID int, roomName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find user's connection
	var connection *WebSocketConnection
	var connectionID string
	for id, conn := range s.connections {
		if conn.User.ID == userID {
			connection = conn
			connectionID = id
			break
		}
	}

	if connection == nil {
		return errors.NewNotFoundError("user connection not found", nil)
	}

	// Remove from room
	if roomConnections, exists := s.rooms[roomName]; exists {
		delete(roomConnections, connectionID)
		if len(roomConnections) == 0 {
			delete(s.rooms, roomName)
		}
	}

	connection.RoomName = ""

	return nil
}

func (s *websocketService) BroadcastMessage(ctx context.Context, message *models.Message) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Find room connections - need to get room name from room ID
	// This is a simplified implementation - in production, you'd need to map room ID to room name
	var roomName string
	for name, connections := range s.rooms {
		for _, conn := range connections {
			if conn.RoomName == name {
				roomName = name
				break
			}
		}
		if roomName != "" {
			break
		}
	}
	
	roomConnections, exists := s.rooms[roomName]
	if !exists {
		return errors.NewNotFoundError("room not found", nil)
	}

	// Broadcast to all connections in the room
	wsMsg := &models.WebSocketMessage{
		Type:      "message",
		Content:   message.Content,
		Sender:    message.Username,
		Timestamp: message.CreatedAt,
	}

	for _, conn := range roomConnections {
		go func(c *WebSocketConnection) {
			if err := websocket.JSON.Send(c.Conn, wsMsg); err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		}(conn)
	}

	return nil
}

func (s *websocketService) BroadcastToRoom(ctx context.Context, roomName string, message *models.WebSocketMessage) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roomConnections, exists := s.rooms[roomName]
	if !exists {
		return errors.NewNotFoundError("room not found", nil)
	}

	for _, conn := range roomConnections {
		go func(c *WebSocketConnection) {
			if err := websocket.JSON.Send(c.Conn, message); err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		}(conn)
	}

	return nil
}

func (s *websocketService) GetConnectedUsers(ctx context.Context, roomName string) ([]*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	roomConnections, exists := s.rooms[roomName]
	if !exists {
		return []*models.User{}, nil
	}

	var users []*models.User
	for _, conn := range roomConnections {
		users = append(users, conn.User)
	}

	return users, nil
}

func generateConnectionID() string {
	// Simplified ID generation
	return fmt.Sprintf("conn_%d", time.Now().UnixNano())
}
