package infrastructure

import (
	commentInfra "github.com/BrockMekonnen/go-clean-starter/internal/comment/infrastructure"
	userInfra "github.com/BrockMekonnen/go-clean-starter/internal/user/infrastructure"
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID          uint                   `gorm:"primaryKey;not null;column:id"`
	Title       string                 `gorm:"type:varchar(255);not null" json:"title"`
	Content     string                 `gorm:"type:text;not null" json:"content"`
	UserId      uint                   `gorm:"not null;index" json:"user_id"`
	User        userInfra.User         `gorm:"foreignKey:UserId"`
	Comments    []commentInfra.Comment `gorm:"foreignKey:PostId"`
	State       string                 `gorm:"type:varchar(20);not null;index" json:"state"`
	PublishedAt *time.Time             `gorm:"type:timestamp" json:"posted_at,omitempty"`
	CreatedAt   time.Time              `gorm:"type:timestamp;not null;column:created_at"`
	UpdatedAt   time.Time              `gorm:"type:timestamp;not null;column:updated_at"`
	Version     int                    `gorm:"type:int;not null;default:0"`
	Deleted     bool                   `gorm:"type:boolean;default:false;column:deleted"`
}

func InitPostTable(db *gorm.DB) error {
	return db.AutoMigrate(&Post{})
}
