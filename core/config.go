package app

import (
	"github.com/BrockMekonnen/go-clean-starter/internal/_lib"
)

// AppConfig holds the configuration values
type AppConfig struct {
	AppName     string
	Environment string
	HTTP        ServerConfig
	Swagger     SwaggerConfig
	Database    DatabaseConfig
}

// SwaggerConfig holds Swagger documentation configurations
type SwaggerConfig struct {
	Title       string
	Version     string
	BasePath    string
	DocEndpoint string
}

// LoadConfig initializes and loads the application configuration
func LoadConfig() *AppConfig {
	_lib.LoadEnv()

	return &AppConfig{
		AppName:     "clean-go-api",
		Environment: _lib.GetEnvironment("development"),
		HTTP: ServerConfig{
			Host: _lib.GetEnvString("HOST", "127.0.0.1"),
			Port: _lib.GetEnvNumber("PORT", 3000),
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
			Username: _lib.GetEnvString("DB_USER", "birukmk"),
			Password: _lib.GetEnvString("DB_PASS", "112544"),
			Port:     _lib.GetEnvNumber("DB_PORT", 5432),
		},
	}
}
