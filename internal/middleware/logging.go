package middleware

import (
	"fmt"
	"time"

	"chat_app/internal/metrics"
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

		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()

		clientIP := c.ClientIP()

		userAgent := c.Request.UserAgent()

		userID, _ := c.Get("user_id")

		username, _ := c.Get("username")

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

		// Metrics
		metrics.HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, toStringStatus(status)).Inc()
		metrics.HTTPRequestDurationSeconds.WithLabelValues(c.Request.Method, path).Observe(latency.Seconds())

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

func toStringStatus(s int) string { return fmt.Sprintf("%d", s) }

func (m *LoggingMiddleware) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

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

		c.Next()

		duration := time.Since(start)

		m.logger.Info("WebSocket Connection Result",
			"path", path,
			"duration", duration,
			"ip", c.ClientIP(),
		)
	}
}
