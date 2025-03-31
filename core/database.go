package app

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DatabaseConfig holds PostgreSQL connection details
type DatabaseConfig struct {
	Database string
	Host     string
	Username string
	Password string
	Port     int
}

// DatabaseRegistry holds the database connection instance
type DatabaseRegistry struct {
	DB *gorm.DB
}

// NewDatabase initializes a new PostgreSQL connection using GORM
func NewDatabase(config DatabaseConfig) (*DatabaseRegistry, func(), error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host, config.Username, config.Password, config.Database, config.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		return nil, nil, err
	}

	log.Println("Connected to PostgreSQL successfully.")

	// Cleanup function to close the database connection when shutting down
	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		log.Println("PostgreSQL connection closed.")
	}

	return &DatabaseRegistry{DB: db}, cleanup, nil
}