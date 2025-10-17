package main

import (
	"chat_app/internal/chat"
	"chat_app/internal/user"
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func main() {
	// Load .env (if present) so environment variables are available
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found or couldn't be loaded; falling back to OS environment variables")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
		log.Fatal("Database credentials are not set in environment variables")
	}
	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ")/" + dbName
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create tables if they don't exist
	createTables(db)

	server := chat.NewServer(db)

	r := gin.Default()
	r.Static("/static", "./static")
	r.POST("/login", user.LoginHandler(db))
	r.POST("/register", user.RegisterHandler(db))

	// API endpoints for rooms and messages
	r.GET("/api/rooms", getRoomsHandler(db))
	r.POST("/api/rooms", createRoomHandler(db))
	r.GET("/api/rooms/:id/messages", getMessagesHandler(db))
	r.POST("/api/rooms/:id/messages", createMessageHandler(db))

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

func createTables(db *sql.DB) {
	// Create rooms table
	roomsTable := `
	CREATE TABLE IF NOT EXISTS rooms (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE,
		description TEXT,
		is_private BOOLEAN DEFAULT FALSE,
		created_by INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

	if _, err := db.Exec(roomsTable); err != nil {
		log.Fatalf("Error creating rooms table: %v", err)
	}

	// Create messages table
	messagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INT AUTO_INCREMENT PRIMARY KEY,
		room_id INT NOT NULL,
		sender VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
	)`

	if _, err := db.Exec(messagesTable); err != nil {
		log.Fatalf("Error creating messages table: %v", err)
	}

	// Create room_members table
	roomMembersTable := `
	CREATE TABLE IF NOT EXISTS room_members (
		id INT AUTO_INCREMENT PRIMARY KEY,
		room_id INT NOT NULL,
		username VARCHAR(255) NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
		UNIQUE KEY unique_room_member (room_id, username)
	)`

	if _, err := db.Exec(roomMembersTable); err != nil {
		log.Fatalf("Error creating room_members table: %v", err)
	}

	// Insert default rooms
	insertDefaultRooms(db)
}

func insertDefaultRooms(db *sql.DB) {
	defaultRooms := []struct {
		name        string
		description string
		isPrivate   bool
	}{
		{"Math 101 - Calculus", "Introduction to Calculus and Differential Equations", false},
		{"Physics Lab", "Advanced Physics Laboratory Sessions", true},
		{"Chemistry Study Group", "Organic Chemistry Study and Discussion", false},
		{"General Discussion", "General chat for all students", false},
	}

	for _, room := range defaultRooms {
		_, err := db.Exec(`
			INSERT IGNORE INTO rooms (name, description, is_private, created_by)
			VALUES (?, ?, ?, 1)
		`, room.name, room.description, room.isPrivate)
		if err != nil {
			log.Printf("Error inserting default room %s: %v", room.name, err)
		}
	}
}

// API Handlers
func getRoomsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT r.id, r.name, r.description, r.is_private, r.created_at,
			       COUNT(rm.username) as member_count
			FROM rooms r
			LEFT JOIN room_members rm ON r.id = rm.room_id
			GROUP BY r.id, r.name, r.description, r.is_private, r.created_at
			ORDER BY r.created_at DESC
		`)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch rooms"})
			return
		}
		defer rows.Close()

		var rooms []map[string]interface{}
		for rows.Next() {
			var id int
			var name, description string
			var isPrivate bool
			var createdAt string
			var memberCount int

			err := rows.Scan(&id, &name, &description, &isPrivate, &createdAt, &memberCount)
			if err != nil {
				log.Printf("Error scanning room: %v", err)
				continue
			}

			rooms = append(rooms, map[string]interface{}{
				"id":           id,
				"name":         name,
				"description":  description,
				"is_private":   isPrivate,
				"member_count": memberCount,
				"created_at":   createdAt,
			})
		}

		c.JSON(200, rooms)
	}
}

func createRoomHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var room struct {
			Name        string `json:"name" binding:"required"`
			Description string `json:"description"`
			IsPrivate   bool   `json:"is_private"`
		}

		if err := c.ShouldBindJSON(&room); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		result, err := db.Exec(`
			INSERT INTO rooms (name, description, is_private, created_by)
			VALUES (?, ?, ?, 1)
		`, room.Name, room.Description, room.IsPrivate)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create room"})
			return
		}

		roomID, _ := result.LastInsertId()
		c.JSON(201, gin.H{
			"id":           roomID,
			"name":         room.Name,
			"description":  room.Description,
			"is_private":   room.IsPrivate,
			"member_count": 0,
		})
	}
}

func getMessagesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")

		// Get room name from room ID
		var roomName string
		err := db.QueryRow("SELECT name FROM rooms WHERE id = ?", roomID).Scan(&roomName)
		if err != nil {
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		}

		rows, err := db.Query(`
			SELECT sender, content, timestamp
			FROM messages
			WHERE room = ?
			ORDER BY timestamp ASC
		`, roomName)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch messages"})
			return
		}
		defer rows.Close()

		var messages []map[string]interface{}
		for rows.Next() {
			var sender, content, timestamp string
			err := rows.Scan(&sender, &content, &timestamp)
			if err != nil {
				log.Printf("Error scanning message: %v", err)
				continue
			}

			messages = append(messages, map[string]interface{}{
				"sender":     sender,
				"content":    content,
				"created_at": timestamp,
			})
		}

		c.JSON(200, messages)
	}
}

func createMessageHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roomID := c.Param("id")

		// Get room name from room ID
		var roomName string
		err := db.QueryRow("SELECT name FROM rooms WHERE id = ?", roomID).Scan(&roomName)
		if err != nil {
			c.JSON(404, gin.H{"error": "Room not found"})
			return
		}

		var message struct {
			Sender  string `json:"sender" binding:"required"`
			Content string `json:"content" binding:"required"`
		}

		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		_, err = db.Exec(`
			INSERT INTO messages (room, sender, content, timestamp)
			VALUES (?, ?, ?, NOW())
		`, roomName, message.Sender, message.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to save message"})
			return
		}

		c.JSON(201, gin.H{"message": "Message saved"})
	}
}
