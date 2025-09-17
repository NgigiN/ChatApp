package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		room := c.Query("room")
		if room == "" {
			room = "General"
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), room: room}
		hub.Join(room, client)

		go client.writePump()
		go client.readPump()
	}
}
