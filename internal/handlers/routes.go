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

	// Initialize services (simplified - in production, use dependency injection)
	// TODO: Properly initialize services with repositories
	roomService := &mockRoomService{}
	userService := &mockUserService{}

	// Initialize handlers
	roomHandlers := NewRoomHandlers(roomService, userService)
	moderationHandlers := NewModerationHandlers(roomService, userService)
	inviteHandlers := NewInviteHandlers(roomService, userService)

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
				// Room management
				rooms.GET("/", roomHandlers.GetUserRooms)              // Get user's rooms
				rooms.POST("/", roomHandlers.CreateRoom)               // Create new room
				rooms.GET("/:id", roomHandlers.GetRoom)                // Get room details
				rooms.PUT("/:id", roomHandlers.UpdateRoom)             // Update room
				rooms.DELETE("/:id", roomHandlers.DeleteRoom)          // Delete room
				rooms.POST("/:id/join", roomHandlers.JoinRoom)         // Join room
				rooms.DELETE("/:id/leave", roomHandlers.LeaveRoom)     // Leave room
				rooms.GET("/:id/members", roomHandlers.GetRoomMembers) // Get room members

				// Moderation routes
				moderation := rooms.Group("/:id/moderation")
				{
					moderation.POST("/remove", moderationHandlers.RemoveUser)             // Remove user from room
					moderation.POST("/reset", moderationHandlers.ResetRoom)               // Reset room (remove all members)
					moderation.GET("/permissions", moderationHandlers.GetRoomPermissions) // Get user permissions
				}

				// Invite routes (for private rooms)
				invites := rooms.Group("/:id/invites")
				{
					invites.POST("/", inviteHandlers.InviteUser)              // Invite single user
					invites.POST("/bulk", inviteHandlers.InviteMultipleUsers) // Invite multiple users
					invites.GET("/users", inviteHandlers.GetInvitableUsers)   // Get invitable users
				}
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
