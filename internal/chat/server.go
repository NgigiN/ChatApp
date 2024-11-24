package chat

import (
	"chat_app/internal/user"
	"database/sql"
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn *websocket.Conn
	room string
	db   *sql.DB
	user *user.User
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
	log.Printf("New client connected: %s", conn.RemoteAddr())
	defer func() {
		s.handleClientDisconnect(client)
		conn.Close()
	}()

	for {
		var msg Message
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		log.Printf("Received message of type %s in room %s from client %s", msg.Type, msg.Room, conn.RemoteAddr())

		switch msg.Type {
		case "join":
			s.handleJoinRoom(client, msg.Room)
		case "message":
			if client.room != "" {
				s.handleMessage(client, msg)
				msg.Room = client.room
				msg.Sender = conn.RemoteAddr().String()
				log.Printf("Handling message in room %s: %s", client.room, msg.Content)

				if room, exists := s.rooms[client.room]; exists {
					room.Broadcast(msg)
				} else {
					log.Printf("Room %s not found", client.room)
				}
			} else {
				log.Printf("Client not in any room")
			}
		}
	}
}

func (s *Server) handleJoinRoom(client *Client, roomName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove the client from their current room if they have one
	if client.room != "" {
		if room, exists := s.rooms[client.room]; exists {
			room.RemoveClient(client)
		}
	}

	// Check if the requested room exists
	if room, exists := s.rooms[roomName]; exists {
		room.AddClient(client)
		client.room = roomName

		response := Message{
			Type:    "join",
			Room:    roomName,
			Content: fmt.Sprintf("Joined room %s", roomName),
		}
		// Send confirmation of joining the room
		if err := websocket.JSON.Send(client.conn, response); err != nil {
			log.Printf("Error sending join confirmation: %v", err)
			return
		}

		// Notify others in the room
		joinMsg := Message{
			Type:    "system",
			Room:    roomName,
			Content: "New user joined the room",
		}
		log.Printf("Broadcasting join message to room %s", roomName)
		room.Broadcast(joinMsg)
	} else {
		log.Printf("Room %s does not exist", roomName)
	}
}

func (s *Server) handleMessage(client *Client, msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if room, exists := s.rooms[client.room]; exists {
		// Use authenticated user's username as the sender
		msg.Sender = client.user.Username

		// Save message in the database
		_, err := s.db.Exec(
			"INSERT INTO messages (room, sender, content, timestamp) VALUES (?, ?, ?, ?)",
			client.room, client.user.Username, msg.Content, msg.Timestamp,
		)
		if err != nil {
			log.Printf("Error saving message to the database: %v", err)
			return
		}

		// Broadcast message to other clients in the room
		log.Printf("Broadcasting message in room %s from user %s: %s", client.room, client.user.Username, msg.Content)
		room.Broadcast(msg)
	}
}

func (s *Server) handleClientDisconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.room != "" {
		if room, exists := s.rooms[client.room]; exists {
			room.RemoveClient(client)
			// Broadcast disconnect message
			disconnectMsg := Message{
				Type:    "system",
				Room:    client.room,
				Content: "A user has left the room",
			}
			room.Broadcast(disconnectMsg)
		}
	}
}
