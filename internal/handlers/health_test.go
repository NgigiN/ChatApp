package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	logger := logger.New("info", "text")
	handler := NewHealthHandler(nil, nil, logger)
	
	router.GET("/health", handler.HealthCheck)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "healthy")
}

func TestLivenessCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	logger := logger.New("info", "text")
	handler := NewHealthHandler(nil, nil, logger)
	
	router.GET("/health/live", handler.LivenessCheck)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/live", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
}

func TestReadinessCheck(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	logger := logger.New("info", "text")
	handler := NewHealthHandler(nil, nil, logger)
	
	router.GET("/health/ready", handler.ReadinessCheck)

	// Test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health/ready", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
}
