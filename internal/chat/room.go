package chat

import (
	"sync"
	"golang.org/x/net/websocket"
)

type Room struct {
	name string
	clients map[*Client]bool
	mu sync.RWMutex
}

func NewRoom(name string) *Room {
	return &Room{
		name: name,
		clients: make(map[*Client]bool),
	}
}

func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[client] = true
}

func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, client)
}

func (r *Room) Broadcast(msg Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		websocket.JSON.Send(client.conn, msg)
	}}