package model

import (
	"time"

	"gorm.io/gorm"
)

// Draft 草稿模型
type Draft struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Title      string         `json:"title" gorm:"type:varchar(255);not null"`
	Content    string         `json:"content" gorm:"type:text"`
	Summary    string         `json:"summary" gorm:"type:text"`
	Cover      string         `json:"cover" gorm:"type:varchar(255)"` // 文章封面图片URL
	CategoryID uint           `json:"category_id"`
	UserID     uint           `json:"user_id" gorm:"not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Category *Category `json:"category" gorm:"foreignKey:CategoryID"`
	User     *User     `json:"user" gorm:"foreignKey:UserID"`
	Tags     []Tag     `json:"tags" gorm:"many2many:draft_tags;"`
}
