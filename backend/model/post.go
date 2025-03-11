package model

import "time"

// Post 表示一篇文章/笔记
type Post struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Content     string    `json:"content" gorm:"type:text"`
	Summary     string    `json:"summary" gorm:"type:text"`
	Cover       string    `json:"cover" gorm:"type:varchar(255)"` // 文章封面图片URL
	Slug        string    `json:"slug" gorm:"uniqueIndex"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	CategoryID  uint      `json:"category_id"`
	Category    Category  `json:"category" gorm:"foreignKey:CategoryID"`
	Tags        []Tag     `json:"tags" gorm:"many2many:post_tags;"`
	Status      string    `json:"status" gorm:"default:'draft'"` // draft, published
	Views       int64     `json:"views" gorm:"default:0"`        // 浏览量
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category 表示文章分类
type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Tag 表示文章标签
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
