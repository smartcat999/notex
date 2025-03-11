package dto

import "time"

// TagInfo 标签信息
type TagInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Slug       string `json:"slug"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
	UserID     uint   `json:"-"` // 内部使用，不从请求参数中绑定
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Slug       string `json:"slug"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
	Status     string `json:"status" binding:"omitempty,oneof=draft published"`
}

// PostResponse 文章响应
type PostResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Summary      string    `json:"summary"`
	Cover        string    `json:"cover"`
	Slug         string    `json:"slug"`
	CategoryID   uint      `json:"category_id"`
	Category     string    `json:"category"`
	Tags         []TagInfo `json:"tags"`
	Status       string    `json:"status"`
	Views        int64     `json:"views"`
	CommentCount int64     `json:"comment_count"`
	PublishedAt  time.Time `json:"published_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PostListQuery 文章列表查询参数
type PostListQuery struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=10"`
	Status     string `form:"status"`
	CategoryID uint   `form:"category_id"`
	TagID      uint   `form:"tag_id"`
	Search     string `form:"search"`
	Sort       string `form:"sort"`
	User       string `form:"user"` // 用于过滤特定用户的文章，值为 "current" 时表示当前用户
	UserID     uint   `form:"-"`    // 内部使用，不从请求参数中绑定
}

// ArchiveResponse 文章归档响应
type ArchiveResponse struct {
	Date  string `json:"date"`  // 归档日期，格式：YYYY-MM
	Count int64  `json:"count"` // 文章数量
}
