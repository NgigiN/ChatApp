package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db     *sql.DB
	redis  *redis.Client
	logger *logger.Logger
}

func NewHealthHandler(db *sql.DB, redis *redis.Client, logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "healthy"
	httpStatus := http.StatusOK
	checks := make(map[string]string)

	// Check database
	if err := h.checkDatabase(); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	} else {
		checks["database"] = "healthy"
	}

	// Check Redis
	if err := h.checkRedis(); err != nil {
		checks["redis"] = "unhealthy: " + err.Error()
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	} else {
		checks["redis"] = "healthy"
	}

	response := gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
		"version":   "1.0.0",
	}

	c.JSON(httpStatus, response)
}

func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	status := "ready"
	httpStatus := http.StatusOK
	checks := make(map[string]string)

	// Check if all critical services are ready
	if err := h.checkDatabase(); err != nil {
		checks["database"] = "not ready: " + err.Error()
		status = "not ready"
		httpStatus = http.StatusServiceUnavailable
	} else {
		checks["database"] = "ready"
	}

	response := gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	}

	c.JSON(httpStatus, response)
}

func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	// Simple liveness check - if the server is responding, it's alive
	c.JSON(http.StatusOK, gin.H{
		"status":    "alive",
		"timestamp": time.Now().UTC(),
	})
}

func (h *HealthHandler) checkDatabase() error {
	if h.db == nil {
		return nil // Database is optional in tests
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.db.PingContext(ctx)
}

func (h *HealthHandler) checkRedis() error {
	if h.redis == nil {
		return nil // Redis is optional
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.redis.Ping(ctx).Err()
}

func (h *HealthHandler) Metrics(c *gin.Context) {
	// Basic metrics endpoint
	stats := gin.H{
		"timestamp": time.Now().UTC(),
		"database": gin.H{
			"max_open_conns": h.db.Stats().MaxOpenConnections,
			"open_conns":     h.db.Stats().OpenConnections,
			"in_use":         h.db.Stats().InUse,
			"idle":           h.db.Stats().Idle,
		},
	}

	if h.redis != nil {
		info, err := h.redis.Info(context.Background()).Result()
		if err == nil {
			stats["redis"] = gin.H{
				"status": "connected",
				"info":   info,
			}
		} else {
			stats["redis"] = gin.H{
				"status": "disconnected",
				"error":  err.Error(),
			}
		}
	}

	c.JSON(http.StatusOK, stats)
}
