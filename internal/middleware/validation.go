package middleware

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ValidationMiddleware struct {
	logger *logger.Logger
}

func NewValidationMiddleware(logger *logger.Logger) *ValidationMiddleware {
	return &ValidationMiddleware{logger: logger}
}

func (m *ValidationMiddleware) SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		if values, err := url.ParseQuery(c.Request.URL.RawQuery); err == nil {
			for key, values := range values {
				for i, value := range values {
					values[i] = m.sanitizeString(value)
				}
				for i, value := range values {
					values[i] = m.sanitizeString(value)
				}
				if len(values) > 0 {
					c.Request.URL.RawQuery = strings.ReplaceAll(c.Request.URL.RawQuery, key+"="+values[0], key+"="+m.sanitizeString(values[0]))
				}
			}

			// Sanitize form data
			if err := c.Request.ParseForm(); err == nil {
				for _, values := range c.Request.PostForm {
					for i, value := range values {
						values[i] = m.sanitizeString(value)
					}
				}
			}

			c.Next()
		}
	}
}

func (m *ValidationMiddleware) ValidateUsername() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		if !m.isValidUsername(req.Username) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid username format"})
			c.Abort()
			return
		}

		c.Set("validated_username", req.Username)
		c.Next()
	}
}

func (m *ValidationMiddleware) ValidateEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		if !m.isValidEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			c.Abort()
			return
		}

		c.Set("validated_email", req.Email)
		c.Next()
	}
}

func (m *ValidationMiddleware) ValidatePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		if !m.isValidPassword(req.Password) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain at least one letter and one number"})
			c.Abort()
			return
		}

		c.Set("validated_password", req.Password)
		c.Next()
	}
}

func (m *ValidationMiddleware) ValidateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Content string `json:"content" binding:"required"`
			Room    string `json:"room" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		if !m.isValidMessage(req.Content) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message content"})
			c.Abort()
			return
		}

		if !m.isValidRoomName(req.Room) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid room name"})
			c.Abort()
			return
		}

		c.Set("validated_content", req.Content)
		c.Set("validated_room", req.Room)
		c.Next()
	}
}

func (m *ValidationMiddleware) sanitizeString(input string) string {
	input = strings.TrimSpace(input)

	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	input = htmlTagRegex.ReplaceAllString(input, "")

	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")

	jsProtocolRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsProtocolRegex.ReplaceAllString(input, "")

	sqlInjectionRegex := regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)`)
	input = sqlInjectionRegex.ReplaceAllString(input, "")

	return input
}

func (m *ValidationMiddleware) isValidUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{3,50}$`)
	return usernameRegex.MatchString(username)
}

func (m *ValidationMiddleware) isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func (m *ValidationMiddleware) isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	return hasLetter && hasNumber
}

func (m *ValidationMiddleware) isValidMessage(content string) bool {
	content = strings.TrimSpace(content)
	if len(content) == 0 || len(content) > 1000 {
		return false
	}

	whitespaceRegex := regexp.MustCompile(`\s{10,}`)
	return !whitespaceRegex.MatchString(content)
}

func (m *ValidationMiddleware) isValidRoomName(roomName string) bool {
	roomNameRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]{1,100}$`)
	return roomNameRegex.MatchString(roomName)
}
