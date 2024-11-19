package main

import (
	"chat_app/internal/chat"
	"chat_app/internal/user"
	"database/sql"
	"log"
	"net/http"

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
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/login", user.LoginHandler(db))
	http.HandleFunc("/register", user.RegisterHandler(db))

	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		server.HandleWS(conn)
	}))

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
