package chat

import (
	"log"
	"sync"

	"golang.org/x/net/websocket"
)

type Room struct {
	name    string
	clients map[*Client]bool
	mu      sync.RWMutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[client] = true
	log.Printf("Client added to room %s. Total clients: %d", r.name, len(r.clients))
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, client)
	log.Printf("Client removed from room %s. Total clients: %d", r.name, len(r.clients))
}

func (r *Room) Broadcast(msg Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	log.Printf("Broadcasting to %d clients in room %s", len(r.clients), r.name)
	for client := range r.clients {
		go func(c *Client) {
			if err := websocket.JSON.Send(c.conn, msg); err != nil {
				log.Printf("Error sending message to client: %v", err)
				return
			}
		}(client)
	}
}
