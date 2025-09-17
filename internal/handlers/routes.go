package handlers

import (
	"chat_app/internal/config"
	"chat_app/internal/middleware"
	"chat_app/internal/ws"
	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db interface{}, redis interface{}, logger *logger.Logger) {
	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(nil, logger) // TODO: inject auth service
	validationMiddleware := middleware.NewValidationMiddleware(logger)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(100, 60, logger) // 100 requests per minute
	securityMiddleware := middleware.NewSecurityMiddleware(logger)
	loggingMiddleware := middleware.NewLoggingMiddleware(logger)

	// Apply global middleware
	router.Use(loggingMiddleware.RequestLogger())
	router.Use(securityMiddleware.SecurityHeaders())
	router.Use(securityMiddleware.CORS())
	router.Use(securityMiddleware.RequestSizeLimit(10 * 1024 * 1024)) // 10MB limit
	router.Use(securityMiddleware.BlockSuspiciousRequests())

	// Health check endpoints
	healthHandler := NewHealthHandler(nil, nil, logger) // TODO: inject db and redis
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.GET("/health/live", healthHandler.LivenessCheck)
	router.GET("/metrics", healthHandler.Metrics)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes
		public := v1.Group("/")
		{
			public.POST("/register", validationMiddleware.ValidateUsername(), validationMiddleware.ValidateEmail(), validationMiddleware.ValidatePassword())
			public.POST("/login", validationMiddleware.ValidateUsername(), validationMiddleware.ValidatePassword())
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(authMiddleware.RequireAuth())
		protected.Use(rateLimitMiddleware.RateLimitPerUser())
		{
			// User routes
			protected.GET("/profile")
			protected.PUT("/profile")
			protected.DELETE("/profile")
			protected.POST("/change-password", validationMiddleware.ValidatePassword())

			// Room routes
			rooms := protected.Group("/rooms")
			rooms.Use(rateLimitMiddleware.RateLimit())
			{
				rooms.GET("/")
				rooms.POST("/", validationMiddleware.ValidateMessage())
				rooms.GET("/:id")
				rooms.PUT("/:id")
				rooms.DELETE("/:id")
				rooms.POST("/:id/join")
				rooms.DELETE("/:id/leave")
				rooms.GET("/:id/members")
			}

			// Message routes
			messages := protected.Group("/messages")
			messages.Use(rateLimitMiddleware.RateLimitPerRoom())
			{
				messages.GET("/")
				messages.POST("/", validationMiddleware.ValidateMessage())
				messages.GET("/:id")
				messages.PUT("/:id")
				messages.DELETE("/:id")
			}
		}
	}

	// WebSocket endpoint with Redis Pub/Sub for cross-instance broadcasting
	hub := ws.NewHub()
	cfg := config.Load()
	redisClient := config.NewRedisClient(cfg.Redis)
	hub.EnableRedis(redisClient)
	go hub.Run()
	router.GET("/ws", ws.ServeWS(hub))

	// Static files
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
}
