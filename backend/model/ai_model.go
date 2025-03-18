package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// AIModel 表示AI模型配置
type AIModel struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Provider    string    `gorm:"size:50;not null" json:"provider"`
	ModelID     string    `gorm:"size:100;not null" json:"modelId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	IsPaid      bool      `gorm:"default:false" json:"isPaid"`
	IsEnabled   bool      `gorm:"default:true" json:"isEnabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName 指定AIModel的表名
func (AIModel) TableName() string {
	return "ai_models"
}

// AIProvider 表示AI提供商配置
type AIProvider struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProviderID  string    `gorm:"size:50;not null;uniqueIndex" json:"providerId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"size:500" json:"description"`
	HasEndpoint bool      `gorm:"default:false" json:"hasEndpoint"`
	IsEnabled   bool      `gorm:"default:true" json:"isEnabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName 指定AIProvider的表名
func (AIProvider) TableName() string {
	return "ai_providers"
}

// JSONMap 用于存储JSON格式的数据
type JSONMap map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// AIUserSetting 表示用户的AI设置
type AIUserSetting struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"userId"`
	ProviderID    string    `gorm:"size:50;not null" json:"providerId"`
	APIKey        string    `gorm:"size:500" json:"apiKey"`
	Endpoint      string    `gorm:"size:500" json:"endpoint"`
	EnabledModels JSONMap   `gorm:"type:json" json:"enabledModels"`
	ModelParams   JSONMap   `gorm:"type:json" json:"modelParams"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TableName 指定AIUserSetting的表名
func (AIUserSetting) TableName() string {
	return "ai_user_settings"
}

// AIDefaultSetting 表示用户的默认AI设置
type AIDefaultSetting struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;uniqueIndex" json:"userId"`
	DefaultModel string    `gorm:"size:100" json:"defaultModel"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName 指定AIDefaultSetting的表名
func (AIDefaultSetting) TableName() string {
	return "ai_default_settings"
}
