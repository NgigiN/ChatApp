package middleware

import (
	"net/http"
	"strings"

	"chat_app/internal/models"
	"chat_app/internal/services"
	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService services.AuthService
	logger      *logger.Logger
}

func NewAuthMiddleware(authService services.AuthService, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			m.logger.Warn("No token provided")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.logger.Warn("Invalid token", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Add user to context
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)

		c.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := m.extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		user, err := m.authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			m.logger.Warn("Invalid token in optional auth", "error", err)
			c.Next()
			return
		}

		// Add user to context if valid
		c.Set("user", user)
		c.Set("user_id", user.ID)
		c.Set("username", user.Username)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRoomAccess(roomService services.RoomService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		roomName := c.Param("room")
		if roomName == "" {
			roomName = c.Query("room")
		}

		if roomName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Room name required"})
			c.Abort()
			return
		}

		// Check if user has access to the room
		room, err := roomService.GetRoomByName(c.Request.Context(), roomName)
		if err != nil {
			m.logger.Warn("Room not found", "room", roomName, "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
			c.Abort()
			return
		}

		// For private rooms, check if user is a member
		if room.IsPrivate {
			// This would need to be implemented in room service
			// For now, we'll allow access if user is authenticated
			m.logger.Info("Private room access", "room", roomName, "user_id", userID)
		}

		c.Set("room", room)
		c.Next()
	}
}

func (m *AuthMiddleware) extractToken(c *gin.Context) string {
	// Try to get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try to get token from query parameter
	token := c.Query("token")
	if token != "" {
		return token
	}

	// Try to get token from cookie
	cookie, err := c.Cookie("auth_token")
	if err == nil && cookie != "" {
		return cookie
	}

	return ""
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// This is a simplified check - in production, you'd have proper role-based access
		// For now, we'll check if username contains "admin"
		userModel := user.(*models.User)
		if !strings.Contains(strings.ToLower(userModel.Username), "admin") {
			m.logger.Warn("Admin access denied", "username", userModel.Username)
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
