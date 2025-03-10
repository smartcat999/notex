package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"user_id"`
	Username   string         `json:"username"`
	Action     string         `json:"action"`      // 操作类型：login, logout, create, update, delete
	Resource   string         `json:"resource"`    // 资源类型：user, note, etc.
	ResourceID string         `json:"resource_id"` // 资源ID
	Details    string         `json:"details"`     // 操作详情，JSON格式
	IP         string         `json:"ip"`          // 操作者IP
	UserAgent  string         `json:"user_agent"`  // 用户代理
	Status     string         `json:"status"`      // success, failed
	Error      string         `json:"error"`       // 错误信息
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
