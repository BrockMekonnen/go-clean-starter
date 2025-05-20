package core

import (
	"time"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/BrockMekonnen/go-clean-starter/core/lib/env"
)

// AppConfig holds the configuration values
type AppConfig struct {
	AppName     string
	Environment string
	StartedAt   time.Time
	Swagger     SwaggerConfig
	HTTP        ServerConfig
	Database    DatabaseConfig
	Encryption  EncryptionConfig
}

// SwaggerConfig holds Swagger documentation configurations
type SwaggerConfig struct {
	Title       string
	Version     string
	BasePath    string
	DocEndpoint string
}

type EncryptionConfig struct {
	JWTKey   string
	HashSalt string
}

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	Host string
	Port int
	Cors bool
}

// DatabaseConfig holds PostgreSQL connection details
type DatabaseConfig struct {
	Database string
	Host     string
	Username string
	Password string
	Port     int
}

// LoadConfig initializes and loads the application configuration
func LoadConfig(logger *logger.Log) *AppConfig {
	env.LoadEnv(logger)

	return &AppConfig{
		AppName:     "clean-go-api",
		Environment: env.GetEnvironment("development"),
		StartedAt:   time.Now(),
		Swagger: SwaggerConfig{
			Title:       "Template API",
			Version:     "1.0.0",
			BasePath:    "/api",
			DocEndpoint: "/api-docs",
		},
		HTTP: ServerConfig{
			Host: env.GetEnvString("HTTP_HOST", "0.0.0.0"),
			Port: env.GetEnvNumber("HTTP_PORT", 9090),
			Cors: true,
		},
		Database: DatabaseConfig{
			Database: env.GetEnvString("DB_NAME", "go-clean"),
			Host:     env.GetEnvString("DB_HOST", "127.0.0.1"),
			Username: env.GetEnvString("DB_USER", "postgres"),
			Password: env.GetEnvString("DB_PASS", "password"),
			Port:     env.GetEnvNumber("DB_PORT", 5432),
		},
		Encryption: EncryptionConfig{
			JWTKey:   env.GetEnvString("JWT_SECRET_KEY", "d6a6a047d84d6884"),
			HashSalt: env.GetEnvString("HASH_SALT", "d84d6884d6a6a047"),
		},
	}
}
