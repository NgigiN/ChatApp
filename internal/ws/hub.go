package ws

import (
	"context"
	"sync"

	"chat_app/internal/metrics"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	rooms      map[string]map[*Client]bool
	register   chan *subscription
	unregister chan *subscription
	broadcast  chan *messageEnvelope
	mu         sync.RWMutex
	pubsub     *redis.Client
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
		// update room state in Redis (optional)
		h.updateRoomState(sub.room, 1)
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
		// update room state in Redis (optional)
		h.updateRoomState(sub.room, -1)
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
				metrics.MessagesBroadcastTotal.WithLabelValues(msg.room).Inc()
				if h.pubsub != nil {
					_ = h.pubsub.Publish(context.Background(), "chat:"+msg.room, msg.data).Err()
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

// EnableRedis enables cross-instance broadcasting via Redis Pub/Sub and starts a subscriber.
func (h *Hub) EnableRedis(client *redis.Client) {
	h.pubsub = client
	go func() {
		if client == nil {
			return
		}
		ctx := context.Background()
		// subscribe to all chat:* channels
		p := client.PSubscribe(ctx, "chat:*")
		ch := p.Channel()
		for msg := range ch {
			room := msg.Channel[len("chat:"):]
			h.broadcast <- &messageEnvelope{room: room, data: []byte(msg.Payload)}
		}
	}()
}

func (h *Hub) updateRoomState(room string, delta int64) {
    if h.pubsub == nil {
        return
    }
    ctx := context.Background()
    // Maintain a set of rooms
    _ = h.pubsub.SAdd(ctx, "rooms", room).Err()
    // Track member counts per room
    count, err := h.pubsub.HIncrBy(ctx, "room:members", room, delta).Result()
    if err == nil && count <= 0 {
        _ = h.pubsub.HDel(ctx, "room:members", room).Err()
        _ = h.pubsub.SRem(ctx, "rooms", room).Err()
    }
}
