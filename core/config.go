package core

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/env"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
)

// AppConfig holds the configuration values
type AppConfig struct {
	AppName     string
	Environment string
	StartedAt   time.Time
	HTTP        ServerConfig
	Swagger     SwaggerConfig
	Database    DatabaseConfig
	Encryption EncryptionConfig
}

// SwaggerConfig holds Swagger documentation configurations
type SwaggerConfig struct {
	Title       string
	Version     string
	BasePath    string
	DocEndpoint string
}

type EncryptionConfig struct {
	JWTKey  string
}

// LoadConfig initializes and loads the application configuration
func LoadConfig(logger *logger.Log) *AppConfig {
	_lib.LoadEnv(logger)

	return &AppConfig{
		AppName:     "clean-go-api",
		Environment: _lib.GetEnvironment("development"),
		StartedAt:   time.Now(),
		HTTP: ServerConfig{
			Host: _lib.GetEnvString("HTTP_HOST", "127.0.0.1"),
			Port: _lib.GetEnvNumber("HTTP_PORT", 3000),
			Cors: true,
		},
		Swagger: SwaggerConfig{
			Title:       "Template API",
			Version:     "1.0.0",
			BasePath:    "/api",
			DocEndpoint: "/api-docs",
		},
		Database: DatabaseConfig{
			Database: _lib.GetEnvString("DB_NAME", "go-clean"),
			Host:     _lib.GetEnvString("DB_HOST", "127.0.0.1"),
			Username: _lib.GetEnvString("DB_USER", "postgres"),
			Password: _lib.GetEnvString("DB_PASS", "password"),
			Port:     _lib.GetEnvNumber("DB_PORT", 5432),
		},
		Encryption: EncryptionConfig{
			JWTKey: _lib.GetEnvString("JWT_SECRET_KEY", "d6a6a047d84d6884"),
		},
	}
}
