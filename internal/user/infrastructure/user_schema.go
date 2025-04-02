package infrastructure

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type User struct {
	Id        uint           `gorm:"primaryKey;not null;column:id"`
	FirstName string         `gorm:"size:100;not null;column:first_name"`
	LastName  string         `gorm:"size:100;not null;column:last_name"`
	Phone     string         `gorm:"size:20;not null;column:phone"`
	Email     string         `gorm:"size:255;not null;column:email"`
	Password  string         `gorm:"size:255;not null;column:password"`
	Roles     pgtype.JSONB       `gorm:"type:jsonb;not null;column:roles"`
	CreatedAt time.Time      `gorm:"autoCreateTime;column:created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;column:updated_at"`
	Version   int            `gorm:"default:1;column:version"`
}

func InitUserTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}