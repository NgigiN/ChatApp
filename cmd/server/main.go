package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat_app/internal/config"
	"chat_app/internal/handlers"
	"chat_app/internal/ws"
	"chat_app/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	var envFile = flag.String("env", "", "Environment file to load (e.g., .env, env.dev)")
	flag.Parse()

	var cfg *config.Config
	if *envFile != "" {
		cfg = config.LoadFromFile(*envFile)
	} else {
		cfg = config.Load()
	}

	logger := logger.New(cfg.Logging.Level, cfg.Logging.Format)

	db, err := config.NewDatabaseConnection(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
	defer config.CloseDatabase(db)

	router := gin.New()
	router.Use(gin.Recovery())

	handlers.SetupRoutes(router, db, nil, logger)

	redisClient := config.NewRedisClient(cfg.Redis)
	_ = redisClient
	_ = ws.NewHub

	server := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.Infof("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
