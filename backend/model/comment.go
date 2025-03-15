package model

import (
	"time"

	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    uint           `json:"post_id" gorm:"not null"`
	ParentID  *uint          `json:"parent_id"`
	ReplyToID *uint          `json:"reply_to_id"`
	Status    string         `json:"status" gorm:"size:20;not null;default:'active'"` // active, hidden, deleted
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Post     *Post     `json:"post" gorm:"foreignKey:PostID"`
	Parent   *Comment  `json:"parent" gorm:"foreignKey:ParentID"`
	ReplyTo  *Comment  `json:"reply_to" gorm:"foreignKey:ReplyToID"`
	Children []Comment `json:"children" gorm:"foreignKey:ParentID"`
}

// TableName specifies the table name for Comment model
func (Comment) TableName() string {
	return "comments"
}
