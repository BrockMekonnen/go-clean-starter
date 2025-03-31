package infrastructure

import (
	"time"
	"gorm.io/gorm"
)

// UserSchema defines the PostgreSQL user table structure
type UserSchema struct {
	Id              uint           `gorm:"primaryKey;autoIncrement"`
	FirstName       string         `gorm:"size:100;not null"`
	LastName        string         `gorm:"size:100;not null"`
	Phone           string         `gorm:"size:20;uniqueIndex;not null"`
	Email           string         `gorm:"size:255;uniqueIndex;not null"`
	Password        string         `gorm:"size:255;not null"`
	Roles           []string       `gorm:"type:text[];not null"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	Version         int            `gorm:"default:1"`
}


// InitUserTable initializes the user table using GORM's AutoMigrate
func InitUserTable(db *gorm.DB) error {
	// Register the model and auto-migrate
	err := db.AutoMigrate(&UserSchema{})
	if err != nil {
		return err
	}

	// Add any additional indexes or constraints not handled by GORM tags
	err = db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`).Error

	return err
}