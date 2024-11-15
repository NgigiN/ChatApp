package main

import (
	"chat_app/internal/chat" // Make sure to replace this with the actual import path of your chat package
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := chat.NewServer()

	http.Handle("/", http.FileServer(http.Dir("static")))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.Handler(server.HandleWS).ServeHTTP(w, r)
	})

	log.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
