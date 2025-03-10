package dto

import "time"

// CreateDraftRequest 创建草稿请求
type CreateDraftRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// UpdateDraftRequest 更新草稿请求
type UpdateDraftRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	Summary    string `json:"summary"`
	CategoryID uint   `json:"category_id"`
	TagIDs     []uint `json:"tag_ids"`
}

// DraftResponse 草稿响应
type DraftResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Summary    string    `json:"summary"`
	CategoryID uint      `json:"category_id"`
	Category   string    `json:"category"`
	Tags       []TagInfo `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DraftListQuery 草稿列表查询参数
type DraftListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Search   string `form:"search"`
}
