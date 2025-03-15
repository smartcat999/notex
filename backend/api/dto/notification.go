package dto

import "time"

// NotificationResponse 通知响应
type NotificationResponse struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`       // 通知类型：comment, reply
	Content   string    `json:"content"`    // 通知内容
	PostID    uint      `json:"post_id"`    // 相关文章ID
	CommentID uint      `json:"comment_id"` // 相关评论ID
	IsRead    bool      `json:"is_read"`    // 是否已读
	CreatedAt time.Time `json:"created_at"` // 创建时间

	// 发送者信息
	SenderID       uint   `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	SenderAvatar   string `json:"sender_avatar"`

	// 文章信息
	PostTitle string `json:"post_title"`
}

// PostBrief 文章简要信息
type PostBrief struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

// CommentInfo 评论信息
type CommentInfo struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
}

// NotificationListQuery 通知列表查询参数
type NotificationListQuery struct {
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=10"`
	Unread   *bool  `form:"unread,omitempty"` // 是否只查询未读通知
	Type     string `form:"type,omitempty"`   // 通知类型过滤：post_comment, comment_reply
}
