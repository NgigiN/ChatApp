package chat

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Room    string `json:"room"`
	Content string `json:"content"`
	Sender  string `json:"sender"`
}

type Client struct {
	conn *websocket.Conn
	room string
}

type Server struct {
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewServer() *Server {
	s := &Server{
		rooms: make(map[string]*Room),
	}

	roomNames := []string{"303", "311", "314", "332", "355", "362"}
	for _, name := range roomNames {
		s.rooms[name] = NewRoom(name)
	}

	return s
}

func (s *Server) HandleWS(conn *websocket.Conn) {
	client := &Client{conn: conn}
	defer conn.Close()

	for {
		var msg Message
		err := websocket.JSON.Receive(conn, &msg)
		if err != nil {
			log.Printf("Error receiving message: %v", err)
			s.handleClientDisconnect(client)
			break
		}

		switch msg.Type {
		case "join":
			s.handleJoinRoom(client, msg.Room)
		case "message":
			s.handleMessage(client, msg)
		}
	}

}

func (s *Server) handleJoinRoom(client *Client, roomName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.room != "" {
		if room, exists := s.rooms[client.room]; exists {
			room.RemoveClient(client)
		}
	}

	if room, exists := s.rooms[roomName]; exists {
		room.AddClient(client)
		client.room = roomName

		response := Message{
			Type:    "join",
			Room:    roomName,
			Content: fmt.Sprintf("Joined room %s", roomName),
		}
		websocket.JSON.Send(client.conn, response)
	}
}

func (s *Server) handleMessage(client *Client, msg Message) {
	if room, exists := s.rooms[client.room]; exists {
		room.Broadcast(msg)
	}
}

func (s *Server) handleClientDisconnect(client *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client.room != "" {
		if room, exists := s.rooms[client.room]; exists {
			room.RemoveClient(client)
		}
	}
}
