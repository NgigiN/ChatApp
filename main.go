package main

import (
	"chat_app/internal/chat"
	"chat_app/internal/user"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/websocket"
)

func main() {
	db, err := sql.Open("mysql", "chat:chat@tcp(127.0.0.1:3306)/chat_app")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	server := chat.NewServer(db)

	r := gin.Default()
	r.Static("/static", "./static")
	r.POST("/login", user.LoginHandler(db))
	r.POST("/register", user.RegisterHandler(db))

	r.GET("/ws", func(c *gin.Context) {
		websocket.Handler(func(conn *websocket.Conn) {
			server.HandleWS(conn)
		}).ServeHTTP(c.Writer, c.Request)
	})

	log.Println("Server started at :8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
