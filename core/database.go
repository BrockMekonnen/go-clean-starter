package core

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/BrockMekonnen/go-clean-starter/core/lib/logger"
)


// DatabaseProvider defines the interface for database access
type DatabaseProvider interface {
	GetDB() *gorm.DB
}

// DatabaseRegistry holds the database connection instance
type DatabaseRegistry struct {
	DB *gorm.DB
}

func (d *DatabaseRegistry) GetDB() *gorm.DB {
    return d.DB
}

// NewDatabase initializes a new PostgreSQL connection using GORM
func NewDatabase(config DatabaseConfig, logger logger.Log) (*DatabaseRegistry, func(), error) {
	
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host, config.Username, config.Password, config.Database, config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to PostgreSQL:", err)
		return nil, nil, err
	}

	logger.Info("Connected to PostgreSQL successfully.")

	// Cleanup function to close the database connection when shutting down
	cleanup := func() {
		sqlDB, _ := db.DB()
		err = sqlDB.Close()
		if err != nil {
			logger.Info("Failed to close PostgresSQL connection")
		}
		logger.Info("PostgreSQL connection closed.")
	}

	return &DatabaseRegistry{db}, cleanup, nil
}

