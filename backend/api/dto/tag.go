package dto

import "time"

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

// TagResponse 标签响应
type TagResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	PostCount int64     `json:"post_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TagListQuery 标签列表查询参数
type TagListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Search   string `form:"search"`
}
