package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	SecretKey     string
	Expiration    time.Duration
	RefreshExpiry time.Duration
}

type LoggingConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8000"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", "30s"),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", "30s"),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", "120s"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "127.0.0.1"),
			Port:         getEnv("DB_PORT", "3306"),
			User:         getEnv("DB_USER", "chat"),
			Password:     getEnv("DB_PASSWORD", "chat"),
			DBName:       getEnv("DB_NAME", "chat_app"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 5),
			MaxLifetime:  getDurationEnv("DB_MAX_LIFETIME", "5m"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			SecretKey:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration:    getDurationEnv("JWT_EXPIRATION", "24h"),
			RefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", "168h"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func NewDatabaseConnection(cfg DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.MaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func CloseDatabase(db *sql.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
	dbIndex := cfg.DB
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       dbIndex,
	})
	return client
}
