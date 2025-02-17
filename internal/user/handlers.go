package user

import (
	"chat_app/database"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user database.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		if user.Role != "student" && user.Role != "teacher" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role, must be 'student' or 'teacher'"})
			return
		}

		if err := CreateUser(db, &user); err != nil {
			log.Printf("Error creating user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "role": user.Role})
	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		user, err := GetUser(db, credentials.Email)
		if err != nil {
			log.Printf("Error getting user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid credentials"})
			return
		}
		if err := VerifyPassword(user.Password, credentials.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Password"})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "role": user.Role, "user_id": user.ID})
	}
}
