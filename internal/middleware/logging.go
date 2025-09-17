package middleware

import (
	"time"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type LoggingMiddleware struct {
	logger *logger.Logger
}

func NewLoggingMiddleware(logger *logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

func (m *LoggingMiddleware) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get response status
		status := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get user agent
		userAgent := c.Request.UserAgent()

		// Get user ID if available
		userID, _ := c.Get("user_id")

		// Get username if available
		username, _ := c.Get("username")

		// Log the request
		m.logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", raw,
			"status", status,
			"latency", latency,
			"ip", clientIP,
			"user_agent", userAgent,
			"user_id", userID,
			"username", username,
		)

		// Log errors
		if status >= 400 {
			m.logger.Error("HTTP Error",
				"method", c.Request.Method,
				"path", path,
				"status", status,
				"latency", latency,
				"ip", clientIP,
				"user_id", userID,
				"username", username,
			)
		}
	}
}

func (m *LoggingMiddleware) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Log any errors that occurred
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				m.logger.Error("Request Error",
					"error", err.Error(),
					"type", err.Type,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"ip", c.ClientIP(),
				)
			}
		}
	}
}

func (m *LoggingMiddleware) WebSocketLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Log WebSocket connection attempt
		m.logger.Info("WebSocket Connection Attempt",
			"path", path,
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log WebSocket connection result
		m.logger.Info("WebSocket Connection Result",
			"path", path,
			"duration", duration,
			"ip", c.ClientIP(),
		)
	}
}
