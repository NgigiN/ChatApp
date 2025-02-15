package user

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("secret"))

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := GetUser(db, credentials.Username)
		if err != nil {
			log.Printf("Error getting user: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := VerifyPassword(user.Password, credentials.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		session, _ := store.Get(c.Request, "session")
		session.Values["user"] = credentials.Username
		if err := session.Save(c.Request, c.Writer); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving session"})
			log.Printf("Error saving session: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": "some-jwt-or-session-token"})
	}
}
