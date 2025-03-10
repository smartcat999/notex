package model

import "time"

// VerificationType 验证类型
type VerificationType string

const (
	VerificationTypeEmail         VerificationType = "email"          // 邮箱验证
	VerificationTypePasswordReset VerificationType = "password_reset" // 密码重置
)

// Verification 验证码模型
type Verification struct {
	ID        uint             `json:"id" gorm:"primaryKey"`
	Email     string           `json:"email" gorm:"not null"`
	Code      string           `json:"code" gorm:"not null"`
	Type      VerificationType `json:"type" gorm:"not null"`
	Used      bool             `json:"used" gorm:"default:false"`
	ExpiresAt time.Time        `json:"expires_at"`
	CreatedAt time.Time        `json:"created_at"`
}

// IsExpired 检查验证码是否过期
func (v *Verification) IsExpired() bool {
	return time.Now().After(v.ExpiresAt)
}

// IsValid 检查验证码是否有效
func (v *Verification) IsValid() bool {
	return !v.Used && !v.IsExpired()
}
