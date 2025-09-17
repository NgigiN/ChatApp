package middleware

import (
	"net/http"
	"sync"
	"time"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int
	window   time.Duration
}

type RateLimitMiddleware struct {
	limiter *RateLimiter
	logger  *logger.Logger
}

func NewRateLimitMiddleware(limit int, window time.Duration, logger *logger.Logger) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: &RateLimiter{
			requests: make(map[string][]time.Time),
			limit:    limit,
			window:   window,
		},
		logger: logger,
	}
}

func (m *RateLimitMiddleware) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !m.limiter.Allow(clientIP) {
			m.logger.Warn("Rate limit exceeded", "ip", clientIP)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": m.limiter.GetRetryAfter(clientIP),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *RateLimitMiddleware) RateLimitPerUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		key := "user_" + string(rune(userID.(int)))

		if !m.limiter.Allow(key) {
			m.logger.Warn("User rate limit exceeded", "user_id", userID)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": m.limiter.GetRetryAfter(key),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *RateLimitMiddleware) RateLimitPerRoom() gin.HandlerFunc {
	return func(c *gin.Context) {
		roomName := c.Param("room")
		if roomName == "" {
			roomName = c.Query("room")
		}

		if roomName == "" {
			c.Next()
			return
		}

		key := "room_" + roomName

		if !m.limiter.Allow(key) {
			m.logger.Warn("Room rate limit exceeded", "room", roomName)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Room rate limit exceeded",
				"retry_after": m.limiter.GetRetryAfter(key),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean up old requests
	if requests, exists := rl.requests[key]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[key] = validRequests
	}

	// Check if we're under the limit
	if len(rl.requests[key]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

func (rl *RateLimiter) GetRetryAfter(key string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	requests, exists := rl.requests[key]
	if !exists || len(requests) == 0 {
		return 0
	}

	// Find the oldest request within the window
	now := time.Now()
	cutoff := now.Add(-rl.window)

	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			retryAfter := int(rl.window.Seconds() - now.Sub(reqTime).Seconds())
			if retryAfter < 0 {
				return 0
			}
			return retryAfter
		}
	}

	return 0
}

func (rl *RateLimiter) Cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	for key, requests := range rl.requests {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requests, key)
		} else {
			rl.requests[key] = validRequests
		}
	}
}
