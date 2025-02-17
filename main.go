package main

import (
	"chat_app/internal/chat"
	"chat_app/internal/user"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lpernett/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	server := chat.NewServer(db)

	r := gin.Default()
	r.Static("/static", "./static")
	r.POST("/login", user.LoginHandler(db))
	r.POST("/register", user.RegisterHandler(db))

	r.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Failed to set websocket upgrade: ", err)
			return
		}
		server.HandleWS(conn, c.Request)
	})

	log.Println("Server started at :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
