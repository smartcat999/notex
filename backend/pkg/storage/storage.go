package storage

import (
	"mime/multipart"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 阿里云 OSS
	StorageTypeCOS   StorageType = "cos"   // 腾讯云 COS
	StorageTypeMinio StorageType = "minio" // MinIO 对象存储
)

// UploadResult 上传结果
type UploadResult struct {
	URL          string `json:"url"`           // 文件访问URL
	ThumbnailURL string `json:"thumbnail_url"` // 缩略图URL（仅图片）
	Filename     string `json:"filename"`      // 文件名
	Size         int64  `json:"size"`          // 文件大小
	Type         string `json:"type"`          // 文件类型
}

// StorageConfig 存储配置
type StorageConfig struct {
	Type          StorageType `json:"type"`           // 存储类型
	URLPrefix     string      `json:"url_prefix"`     // URL前缀
	MaxSize       int64       `json:"max_size"`       // 最大文件大小
	AllowedTypes  []string    `json:"allowed_types"`  // 允许的文件类型
	ThumbnailSize int         `json:"thumbnail_size"` // 缩略图大小

	// 本地存储配置
	LocalConfig struct {
		UploadDir string `json:"upload_dir"` // 上传目录
	} `json:"local_config"`

	// 对象存储通用配置
	ObjectStorage struct {
		Endpoint     string `json:"endpoint"`      // 存储服务地址
		AccessKey    string `json:"access_key"`    // 访问密钥ID
		AccessSecret string `json:"access_secret"` // 访问密钥密码
		BucketName   string `json:"bucket_name"`   // 存储桶名称
		Region       string `json:"region"`        // 地域
	} `json:"object_storage"`
}

// Storage 存储接口
type Storage interface {
	// Upload 上传文件
	Upload(file multipart.File, header *multipart.FileHeader) (*UploadResult, error)

	// Delete 删除文件
	Delete(fileURL string) error

	// GetUploadConfig 获取上传配置（返回给前端）
	GetUploadConfig() map[string]interface{}
}

// StorageProvider 存储提供者接口
type StorageProvider interface {
	// NewStorage 创建存储实例
	NewStorage(config *StorageConfig) (Storage, error)
}
