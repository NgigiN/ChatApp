package database

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName string `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(100);not null" json:"last_name"`
	Phone     string `gorm:"type:varchar(100);not null" json:"phone"`
	Email     string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password  string `gorm:"type:varchar(100);not null" json:"password"`
	Role      string `gorm:"type:varchar(100);not null" json:"role"`
}

type Message struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Room      string    `json:"room"`
}

type Room struct {
	gorm.Model
	Name    string                   `gorm:"type:varchar(100);unique;not null" json:"name"`
	Clients map[*websocket.Conn]bool `gorm:"-" json:"-"`
	mu      sync.RWMutex             `gorm:"-" json:"-"`
}

func NewRoom(name string) *Room {
	return &Room{
		Name:    name,
		Clients: make(map[*websocket.Conn]bool),
	}
}

func InitdB() (*gorm.DB, error) {
	dsn := "host=localhost user=chat password=chat dbname=chat_app port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Message{}, &Room{})

	return db, nil
}
