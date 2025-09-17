package ws

import (
	"sync"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *subscription
	unregister chan *subscription
	broadcast  chan *messageEnvelope
	mu         sync.RWMutex
}

type subscription struct {
	client *Client
	room   string
}

type messageEnvelope struct {
	room string
	data []byte
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *subscription, 1024),
		unregister: make(chan *subscription, 1024),
		broadcast:  make(chan *messageEnvelope, 4096),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case sub := <-h.register:
			if _, ok := h.rooms[sub.room]; !ok {
				h.rooms[sub.room] = make(map[*Client]bool)
			}
			h.rooms[sub.room][sub.client] = true
		case sub := <-h.unregister:
			if clients, ok := h.rooms[sub.room]; ok {
				if _, exists := clients[sub.client]; exists {
					delete(clients, sub.client)
					close(sub.client.send)
					if len(clients) == 0 {
						delete(h.rooms, sub.room)
					}
				}
			}
		case msg := <-h.broadcast:
			if clients, ok := h.rooms[msg.room]; ok {
				for c := range clients {
					select {
					case c.send <- msg.data:
					default:
						// backpressure: drop slow client
						close(c.send)
						delete(clients, c)
					}
				}
			}
		}
	}
}

func (h *Hub) Join(room string, c *Client)  { h.register <- &subscription{client: c, room: room} }
func (h *Hub) Leave(room string, c *Client) { h.unregister <- &subscription{client: c, room: room} }
func (h *Hub) Broadcast(room string, payload []byte) {
	h.broadcast <- &messageEnvelope{room: room, data: payload}
}
