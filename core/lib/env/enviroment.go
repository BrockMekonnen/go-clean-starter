package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
	"github.com/joho/godotenv"
)

// Load environment variables based on GO_APP_ENV
func LoadEnv(logger *logger.Log) {
	env := os.Getenv("GO_APP_ENV")
	if env == "" {
		env = "development"
		err := os.Setenv("GO_APP_ENV", env)
		if err != nil {
			logger.Fatal("Failed to load Env", err)
		}
	}

	var envFile string
	if env == "production" {
		envFile = ".env"
	} else if _, err := os.Stat(fmt.Sprintf(".env.%s.local", env)); err == nil {
		envFile = fmt.Sprintf(".env.%s.local", env)
	} else {
		envFile = fmt.Sprintf(".env.%s", env)
	}

	_ = godotenv.Load(envFile)
}

// Environment types
var validEnvironments = []string{"development", "production", "test"}

// EnvironmentConfig struct
type EnvironmentConfig struct {
	Environment string
}

// GetEnvironment validates and retrieves the NODE_ENV variable
func GetEnvironment(defaultValue string) string {
	env := os.Getenv("GO_APP_ENV")
	if env == "" {
		env = defaultValue
		err := os.Setenv("GO_APP_ENV", env)
		if err != nil {
			fmt.Print("error setting GO_APP_ENV")
		}
	}

	for _, validEnv := range validEnvironments {
		if env == validEnv {
			return env
		}
	}

	panic(fmt.Sprintf("Invalid NODE_ENV value. Accepted values: %v", validEnvironments))
}

// GetEnvString retrieves a string environment variable with a fallback
func GetEnvString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" && defaultValue == "" {
		panic(fmt.Sprintf("Required environment variable %s is undefined and has no default", key))
	}
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvNumber retrieves a number environment variable with a fallback
func GetEnvNumber(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == 0 {
			panic(fmt.Sprintf("Required environment variable %s is undefined and has no default", key))
		}
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("Invalid number format for %s: %s", key, value))
	}
	return intValue
}
