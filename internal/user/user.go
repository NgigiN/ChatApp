package user

import (
	"chat_app/database"

	"gorm.io/gorm"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(db *gorm.DB, user *database.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	return db.Create(user).Error
}

func GetUser(db *gorm.DB, email string) (*database.User, error) {
	var user database.User
	if err := db.Where("email=?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func VerifyPassword(storedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
}
