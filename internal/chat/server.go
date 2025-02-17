package chat

import (
	"chat_app/database"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var store = sessions.NewCookieStore([]byte("secret"))

// ChatRoom extends the database Room with websocket functionality
type ChatRoom struct {
	database.Room
	Clients map[*websocket.Conn]*Client
	mu      sync.RWMutex
}

type Client struct {
	Conn     *websocket.Conn
	RoomName string
	DB       *gorm.DB
	User     *database.User
}

type Server struct {
	rooms map[string]*ChatRoom
	mu    sync.RWMutex
	db    *gorm.DB
}

type Message struct {
	database.Message
}

func NewServer(db *gorm.DB) *Server {
	s := &Server{
		rooms: make(map[string]*ChatRoom),
		db:    db,
	}

	// Initialize default rooms
	roomNames := []string{"General", "Gaming", "Technology", "Movies", "Music", "Random"}
	for _, name := range roomNames {
		var dbRoom database.Room
		if err := db.FirstOrCreate(&dbRoom, database.Room{Name: name}).Error; err != nil {
			log.Printf("Error creating room %s: %v", name, err)
			continue
		}

		// Create chat room
		s.rooms[name] = &ChatRoom{
			Room:    dbRoom,
			Clients: make(map[*websocket.Conn]*Client),
		}
	}

	return s
}

func (s *Server) HandleWS(conn *websocket.Conn, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
		}
		s.handleClientDisconnect(&Client{Conn: conn})
		conn.Close()
	}()

	// Default user
	defaultUser := &database.User{
		FirstName: fmt.Sprintf("User_%d", time.Now().Unix()),
	}

	// Get session
	session, _ := store.Get(r, "session")
	if userID, ok := session.Values["user_id"].(uint); ok {
		var user database.User
		if err := s.db.First(&user, userID).Error; err == nil {
			defaultUser = &user
		}
	}

	client := &Client{
		Conn: conn,
		DB:   s.db,
		User: defaultUser,
	}

	defaultRoom := "General"
	s.handleJoinRoom(client, defaultRoom)

	for {
		var msg database.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		msg.Sender = client.User.FirstName
		msg.Timestamp = time.Now()

		switch msg.Type {
		case "join":
			s.handleJoinRoom(client, msg.Room)
		case "message":
			if client.RoomName != "" {
				s.handleMessage(client, &msg)
			}
		}
	}
}

func (s *Server) handleJoinRoom(client *Client, roomName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove client from current room
	if client.RoomName != "" {
		if currentRoom, exists := s.rooms[client.RoomName]; exists {
			currentRoom.removeClient(client)
		}
	}

	// Join new room
	newRoom, exists := s.rooms[roomName]
	if !exists {
		log.Printf("Room %s does not exist", roomName)
		return
	}

	newRoom.addClient(client)
	client.RoomName = roomName

	// Fetch recent messages
	var messages []database.Message
	if err := s.db.Where("room_id = ?", newRoom.ID).Order("created_at desc").Limit(50).Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages: %v", err)
	} else {
		if err := client.Conn.WriteJSON(messages); err != nil {
			log.Printf("Error sending message history: %v", err)
		}
	}

	// Send join notification
	joinMsg := database.Message{
		Type:      "system",
		Content:   fmt.Sprintf("You have joined the room: %s", roomName),
		Sender:    client.User.FirstName,
		Room:      newRoom.Name,
		Timestamp: time.Now(),
	}
	client.Conn.WriteJSON(joinMsg)
}

func (cr *ChatRoom) addClient(client *Client) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	cr.Clients[client.Conn] = client
	log.Printf("Client added to room %s. Total clients: %d", cr.Name, len(cr.Clients))
}

func (cr *ChatRoom) removeClient(client *Client) {
	cr.mu.Lock()
	defer cr.mu.Unlock()
	delete(cr.Clients, client.Conn)
	log.Printf("Client removed from room %s. Total clients: %d", cr.Name, len(cr.Clients))
}

func (cr *ChatRoom) broadcast(msg *database.Message) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	for _, client := range cr.Clients {
		go func(c *Client) {
			if err := c.Conn.WriteJSON(msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
			}
		}(client)
	}
}

func (s *Server) handleMessage(client *Client, msg *database.Message) {
	if client.RoomName == "" {
		log.Println("Client not in a room, message ignored")
		return
	}

	room, exists := s.rooms[client.RoomName]
	if !exists {
		log.Printf("Room %s does not exist", client.RoomName)
		return
	}

	// Save message to database
	if err := s.db.Create(msg).Error; err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	// Broadcast to room
	room.broadcast(msg)
}

func (s *Server) handleClientDisconnect(client *Client) {
	if client.RoomName == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	room, exists := s.rooms[client.RoomName]
	if exists && client.User != nil {
		room.removeClient(client)

		disconnectMsg := &database.Message{
			Type:      "system",
			Content:   fmt.Sprintf("User %s has left the room", client.User.FirstName),
			Sender:    client.User.FirstName,
			Room:      room.Name,
			Timestamp: time.Now(),
		}
		room.broadcast(disconnectMsg)
	}
}
