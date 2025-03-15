package dto

import "time"

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ParentID  *uint  `json:"parent_id,omitempty"`   // 父评论ID，用于回复
	ReplyToID *uint  `json:"reply_to_id,omitempty"` // 具体回复的评论ID
}

// CommentRequest represents the request body for creating a comment
type CommentRequest struct {
	Content   string `json:"content" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
	ReplyToID *uint  `json:"reply_to_id"` // ID of the comment being replied to
}

// CommentResponse represents a comment with user information
type CommentResponse struct {
	ID         uint      `json:"id"`
	Content    string    `json:"content"`
	PostID     uint      `json:"post_id"`
	PostTitle  string    `json:"post_title,omitempty"`
	UserID     uint      `json:"user_id"`
	ParentID   *uint     `json:"parent_id,omitempty"`
	User       *UserInfo `json:"user"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ReplyCount int       `json:"reply_count"` // 子评论数量

	// 回复相关
	Parent   *CommentBrief   `json:"parent,omitempty"`   // Parent comment if this is a reply
	ReplyTo  *CommentBrief   `json:"reply_to,omitempty"` // The specific comment being replied to
	Children []*CommentBrief `json:"children,omitempty"` // Child comments/replies
}

// CommentBrief represents a brief version of a comment
type CommentBrief struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	User      *UserInfo `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	ReplyTo   *struct { // 添加被回复评论的信息
		ID   uint      `json:"id"`
		User *UserInfo `json:"user"`
	} `json:"reply_to,omitempty"`
}

// UserInfo 评论用户信息
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// CommentListQuery 评论列表查询参数
type CommentListQuery struct {
	Page         int  `form:"page,default=1"`
	PageSize     int  `form:"page_size,default=10"`
	WithChildren bool `form:"with_children"` // 是否包含子评论
}
