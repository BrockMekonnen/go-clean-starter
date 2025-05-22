package infrastructure

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint      `gorm:"primaryKey;not null;column:id" json:"id"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	UserId    uint      `gorm:"not null;index;column:user_id" json:"user_id"`
	PostId    uint      `gorm:"not null;index;column:post_id" json:"post_id"`
	Status    string    `gorm:"type:varchar(20);not null;index" json:"status"`
	CreatedAt time.Time `gorm:"not null;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;column:updated_at" json:"updated_at"`
	Version   int       `gorm:"type:int;not null;default:0" json:"version"`
	Deleted   bool      `gorm:"type:boolean;default:false;column:deleted"`
}

func InitCommentTable(db *gorm.DB) error {
	return db.AutoMigrate(&Comment{})
}
