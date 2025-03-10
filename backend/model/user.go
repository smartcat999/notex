package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleUser   = "user"
	RoleEditor = "editor"
	RoleAdmin  = "admin"
)

// User 用户模型
type User struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Username      string         `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Email         string         `json:"email" gorm:"size:100;not null;uniqueIndex"`
	Password      string         `json:"-" gorm:"size:100;not null"`                      // json:"-" 确保密码不会被序列化
	Role          string         `json:"role" gorm:"size:20;not null;default:'user'"`     // admin, editor, user
	Status        string         `json:"status" gorm:"size:20;not null;default:'active'"` // active, inactive, banned
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	Bio           string         `json:"bio" gorm:"size:500"`
	Avatar        string         `json:"avatar" gorm:"size:255"`
	LastLogin     time.Time      `json:"last_login"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// SetPassword 设置用户密码（加密）
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsAdmin 检查用户是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEditor 检查用户是否是编辑
func (u *User) IsEditor() bool {
	return u.Role == "editor" || u.Role == "admin"
}

// IsActive 检查用户是否处于活动状态
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// BeforeCreate 在创建用户前设置默认角色
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}
