package user

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(db *sql.DB, user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, string(hashedPassword),
	)
	return err
}

func GetUser(db *sql.DB, username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func VerifyPassword(storedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
}
