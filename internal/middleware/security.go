package middleware

import (
	"net/http"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type SecurityMiddleware struct {
	logger *logger.Logger
}

func NewSecurityMiddleware(logger *logger.Logger) *SecurityMiddleware {
	return &SecurityMiddleware{logger: logger}
}

func (m *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")

		c.Header("X-Content-Type-Options", "nosniff")

		c.Header("X-XSS-Protection", "1; mode=block")

		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self' ws: wss:;")

		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}

func (m *SecurityMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}

		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func (m *SecurityMiddleware) RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			m.logger.Warn("Request too large", "size", c.Request.ContentLength, "max", maxSize)
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "Request too large"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *SecurityMiddleware) IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		allowed := false
		for _, allowedIP := range allowedIPs {
			if clientIP == allowedIP {
				allowed = true
				break
			}
		}

		if !allowed {
			m.logger.Warn("IP not whitelisted", "ip", clientIP)
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *SecurityMiddleware) BlockSuspiciousRequests() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		path := c.Request.URL.Path
		suspiciousUserAgents := []string{
			"sqlmap",
			"nikto",
			"nmap",
			"masscan",
			"zap",
			"burp",
		}

		for _, suspicious := range suspiciousUserAgents {
			if containsIgnoreCase(userAgent, suspicious) {
				m.logger.Warn("Suspicious user agent blocked", "user_agent", userAgent, "path", path)
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		suspiciousPaths := []string{
			"/admin",
			"/wp-admin",
			"/phpmyadmin",
			"/.env",
			"/config",
			"/backup",
		}

		for _, suspicious := range suspiciousPaths {
			if containsIgnoreCase(path, suspicious) {
				m.logger.Warn("Suspicious path blocked", "path", path, "ip", c.ClientIP())
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr)))
}
