package chat

import (
	"chat_app/internal/user"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/net/websocket"
)

var store = sessions.NewCookieStore([]byte("secret"))

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered: %v", r)
		}
		s.handleClientDisconnect(&Client{conn: conn})
		conn.Close()
	}()

	defaultUser := &user.User{
		Username: fmt.Sprintf("DefaultUser_%d", time.Now().Hour()),
	}

	session, _ := store.Get(conn.Request(), "session")
	if username, ok := session.Values["user"].(string); ok {
		defaultUser.Username = username
	}

	client := &Client{conn: conn, db: s.db, user: defaultUser}

	defaultRoom := "General"
	s.handleJoinRoom(client, &defaultRoom)

	for {
		var msg *Message
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			log.Printf("Error receiving message: %v", err)
			break
		}

		msg.Sender = client.user.Username
		msg.Timestamp = time.Now()

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

	if client.roomName != nil {
		if currentRoom, exists := s.rooms[*client.roomName]; exists {
			currentRoom.RemoveClient(client)
		}
	}

	newRoom, exists := s.rooms[*roomName]
	if !exists {
		log.Printf("Room %s does not exist - creating new room", *roomName)
		s.rooms[*roomName] = NewRoom(*roomName)
		newRoom = s.rooms[*roomName]
	}
	newRoom.AddClient(client)
	client.roomName = roomName

	messages, err := GetMessages(client.db, *roomName)
	if err != nil {
		log.Printf("Error fetching messages: %v", err)
		return
	}

	if err := websocket.JSON.Send(client.conn, messages); err != nil {
		log.Printf("Error sending message history: %v", err)
	}

	joinMessage := &Message{
		Type:    "system",
		Content: fmt.Sprintf("You have joined the room: %s", *roomName),
	}
	_ = websocket.JSON.Send(client.conn, joinMessage)
}

func (s *Server) handleMessage(client *Client, msg *Message) {
	if client == nil || client.roomName == nil {
		log.Println("Client not in a room, message ignored")
		return
	}

	roomName := *client.roomName
	msg.Sender = client.user.Username
	msg.Timestamp = time.Now()
	msg.Room = roomName

	if err := SaveMessage(client.db, msg); err != nil {
		log.Printf("Error saving message: %v", err)
		return
	}

	room, exists := s.rooms[roomName]
	if exists {
		room.Broadcast(msg)
	}
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
