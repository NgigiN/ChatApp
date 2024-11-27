// internal/chat/server.go
package chat

import (
	"chat_app/internal/user"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn     *websocket.Conn
	roomName *string
	db       *sql.DB
	user     *user.User
}

type Server struct {
	rooms map[string]*Room
	mu    sync.RWMutex
	db    *sql.DB
}

func NewServer(db *sql.DB) *Server {
	s := &Server{
		rooms: make(map[string]*Room),
		db:    db,
	}

	// Initialize default rooms
	roomNames := []string{"General", "Gaming", "Technology", "Movies", "Music", "Random"}
	for _, name := range roomNames {
		s.rooms[name] = NewRoom(name)
	}
	return s
}

func (s *Server) HandleWS(conn *websocket.Conn) {
	client := &Client{conn: conn, db: s.db}
	defer func() {
		s.handleClientDisconnect(client)
		conn.Close()
	}()

	for {
		var msg *Message
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		switch msg.Type {
		case "join":
			s.handleJoinRoom(client, &msg.Room)
		case "message":
			if client.roomName != nil {
				s.handleMessage(client, msg)
			}
		}
	}
}

func (s *Server) handleJoinRoom(client *Client, roomName *string) {
	if roomName == nil {
		log.Println("handleJoinRoom: roomName is nil")
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the client from their current room
	if client.roomName != nil {
		currentRoom, exists := s.rooms[*client.roomName]
		if exists {
			currentRoom.RemoveClient(client)
		}
	}

	// Join the new room
	newRoom, exists := s.rooms[*roomName]
	if !exists {
		log.Printf("Room %s does not exist", *roomName)
		return
	}

	newRoom.AddClient(client)
	client.roomName = roomName

	response := &Message{
		Type:    "join",
		Room:    *roomName,
		Content: fmt.Sprintf("Joined room %s", *roomName),
	}
	_ = websocket.JSON.Send(client.conn, response)
}

func (s *Server) handleMessage(client *Client, msg *Message) {
	if client == nil || client.roomName == nil {
		log.Println("Client not associated with any room, message ignored")
		return
	}

	roomName := *client.roomName
	msg.Sender = client.user.Username
	msg.Timestamp = time.Now()
	msg.Room = roomName

	if err := SaveMessage(client.db, msg); err != nil {
		log.Printf("Error saving message: %v", err)
	}

	room, exists := s.rooms[roomName]
	if !exists || room == nil {
		log.Printf("Room %s does not exist, cannot broadcast message", roomName)
		return
	}

	room.Broadcast(msg)
}

func (s *Server) handleClientDisconnect(client *Client) {
	if client == nil || client.roomName == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[*client.roomName]
	if exists && client.user != nil {
		room.RemoveClient(client)

		disconnectMsg := &Message{
			Type:    "system",
			Room:    *client.roomName,
			Content: fmt.Sprintf("User %s has left the room", client.user.Username),
		}
		room.Broadcast(disconnectMsg)
	}
}
