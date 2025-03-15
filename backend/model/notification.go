package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	NotificationTypeCommentReply = "comment_reply"
	NotificationTypePostComment  = "post_comment"
)

// Notification 通知模型
type Notification struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Type      string         `json:"type" gorm:"size:50;not null"` // 通知类型
	UserID    uint           `json:"user_id" gorm:"not null"`      // 接收通知的用户ID
	ActorID   uint           `json:"actor_id" gorm:"not null"`     // 触发通知的用户ID
	PostID    *uint          `json:"post_id"`                      // 相关文章ID
	CommentID *uint          `json:"comment_id"`                   // 相关评论ID
	Content   string         `json:"content" gorm:"type:text"`     // 通知内容
	Read      bool           `json:"read" gorm:"default:false"`    // 是否已读
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	User    *User    `json:"user" gorm:"foreignKey:UserID"`       // 接收通知的用户
	Actor   *User    `json:"actor" gorm:"foreignKey:ActorID"`     // 触发通知的用户
	Post    *Post    `json:"post" gorm:"foreignKey:PostID"`       // 相关文章
	Comment *Comment `json:"comment" gorm:"foreignKey:CommentID"` // 相关评论
}
