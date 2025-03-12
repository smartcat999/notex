package storage

import (
	"mime/multipart"
	"notex/pkg/types"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 阿里云 OSS
	StorageTypeCOS   StorageType = "cos"   // 腾讯云 COS
	StorageTypeMinio StorageType = "minio" // MinIO 对象存储
)

// UploadConfig 上传配置基础结构
type UploadConfig struct {
	Type         StorageType `json:"type"`         // 存储类型
	MaxSize      int64       `json:"maxSize"`      // 最大文件大小
	AllowedTypes []string    `json:"allowedTypes"` // 允许的文件类型
	DirectUpload bool        `json:"directUpload"` // 是否直传
	UploadURL    string      `json:"uploadUrl"`    // 上传URL
}

// LocalUploadConfig 本地存储上传配置
type LocalUploadConfig struct {
	UploadConfig
}

// OSSUploadConfig 阿里云OSS上传配置
type OSSUploadConfig struct {
	UploadConfig
	Region    string            `json:"region"`    // 地域
	Bucket    string            `json:"bucket"`    // Bucket名称
	Endpoint  string            `json:"endpoint"`  // 访问域名
	Headers   map[string]string `json:"headers"`   // 请求头
	Callback  string            `json:"callback"`  // 回调地址
	URLPrefix string            `json:"urlPrefix"` // 访问URL前缀
}

// COSUploadConfig 腾讯云COS上传配置
type COSUploadConfig struct {
	UploadConfig
	Region    string            `json:"region"`    // 地域
	Bucket    string            `json:"bucket"`    // Bucket名称
	Endpoint  string            `json:"endpoint"`  // 访问域名
	Headers   map[string]string `json:"headers"`   // 请求头
	URLPrefix string            `json:"urlPrefix"` // 访问URL前缀
}

// MinioUploadConfig MinIO上传配置
type MinioUploadConfig struct {
	UploadConfig
	Bucket    string `json:"bucket"`    // Bucket名称
	Endpoint  string `json:"endpoint"`  // 访问域名
	Region    string `json:"region"`    // 地域
	URLPrefix string `json:"urlPrefix"` // 访问URL前缀
}

// UploadResult 上传结果
type UploadResult struct {
	URL          string `json:"url"`           // 文件访问URL
	ThumbnailURL string `json:"thumbnail_url"` // 缩略图URL（仅图片）
	Filename     string `json:"filename"`      // 文件名
	Size         int64  `json:"size"`          // 文件大小
	Type         string `json:"type"`          // 文件类型
}

// UploadCredentials 上传凭证
type UploadCredentials struct {
	AccessKeyId     string            `json:"accessKeyId,omitempty"`     // 临时访问密钥ID
	AccessKeySecret string            `json:"accessKeySecret,omitempty"` // 临时访问密钥密码
	SecurityToken   string            `json:"securityToken,omitempty"`   // 安全令牌
	Authorization   string            `json:"authorization,omitempty"`   // 签名
	FileKey         string            `json:"fileKey,omitempty"`         // 文件路径
	Expires         int64             `json:"expires,omitempty"`         // 过期时间（秒）
	Extra           map[string]string `json:"extra,omitempty"`           // 额外参数
}

// Storage 存储接口
type Storage interface {
	// Upload 上传文件
	Upload(file multipart.File, header *multipart.FileHeader) (*UploadResult, error)

	// Delete 删除文件
	Delete(fileURL string) error

	// GetUploadConfig 获取上传配置（返回给前端）
	GetUploadConfig() interface{}

	// GetType 获取存储类型
	GetType() StorageType

	// GetCredentials 获取上传凭证
	GetCredentials(filename string, contentType string) (*UploadCredentials, error)
}

// StorageProvider 存储提供者接口
type StorageProvider interface {
	// NewStorage 创建存储实例
	NewStorage(config *types.StorageConfig) (Storage, error)
}
