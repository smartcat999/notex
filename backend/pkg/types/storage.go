package types

import "fmt"

// StorageConfig 存储配置
type StorageConfig struct {
	Type          string   `yaml:"type" json:"type"`                     // 存储类型：local, oss, cos, minio
	MaxSize       int64    `yaml:"max_size" json:"max_size"`             // 最大文件大小（字节）
	AllowedTypes  []string `yaml:"allowed_types" json:"allowed_types"`   // 允许的文件类型
	ThumbnailSize int      `yaml:"thumbnail_size" json:"thumbnail_size"` // 缩略图大小

	// 本地存储配置
	Local struct {
		UploadDir string `yaml:"upload_dir" json:"upload_dir"` // 上传目录
		URLPrefix string `yaml:"url_prefix" json:"url_prefix"` // 访问URL前缀
	} `yaml:"local" json:"local"`

	// OSS配置
	OSS struct {
		Endpoint     string `yaml:"endpoint" json:"endpoint"`           // OSS服务地址
		AccessKey    string `yaml:"access_key" json:"access_key"`       // 访问密钥ID
		AccessSecret string `yaml:"access_secret" json:"access_secret"` // 访问密钥密码
		BucketName   string `yaml:"bucket_name" json:"bucket_name"`     // 存储桶名称
		Region       string `yaml:"region" json:"region"`               // 地域
		RoleArn      string `yaml:"role_arn" json:"role_arn"`           // RAM角色ARN
		URLPrefix    string `yaml:"url_prefix" json:"url_prefix"`       // 访问URL前缀
	} `yaml:"oss" json:"oss"`

	// COS配置
	COS struct {
		Endpoint     string `yaml:"endpoint" json:"endpoint"`
		AccessKey    string `yaml:"access_key" json:"access_key"`
		AccessSecret string `yaml:"access_secret" json:"access_secret"`
		BucketName   string `yaml:"bucket_name" json:"bucket_name"`
		Region       string `yaml:"region" json:"region"`
		URLPrefix    string `yaml:"url_prefix" json:"url_prefix"`
	} `yaml:"cos" json:"cos"`

	// MinIO配置
	MinIO struct {
		Endpoint     string `yaml:"endpoint" json:"endpoint"`
		AccessKey    string `yaml:"access_key" json:"access_key"`
		AccessSecret string `yaml:"access_secret" json:"access_secret"`
		BucketName   string `yaml:"bucket_name" json:"bucket_name"`
		Region       string `yaml:"region" json:"region"`
		URLPrefix    string `yaml:"url_prefix" json:"url_prefix"`
		UseSSL       bool   `yaml:"use_ssl" json:"use_ssl"`
	} `yaml:"minio" json:"minio"`
}

// Validate 验证存储配置
func (c *StorageConfig) Validate() error {
	if c.MaxSize <= 0 {
		return fmt.Errorf("max_size should be positive")
	}

	if len(c.AllowedTypes) == 0 {
		return fmt.Errorf("allowed_types cannot be empty")
	}

	if c.ThumbnailSize <= 0 {
		return fmt.Errorf("thumbnail_size should be positive")
	}

	switch c.Type {
	case "local":
		if c.Local.UploadDir == "" {
			return fmt.Errorf("upload_dir cannot be empty for local storage")
		}
		if c.Local.URLPrefix == "" {
			return fmt.Errorf("url_prefix cannot be empty for local storage")
		}
	case "oss":
		if c.OSS.Endpoint == "" {
			return fmt.Errorf("endpoint cannot be empty for OSS")
		}
		if c.OSS.AccessKey == "" {
			return fmt.Errorf("access_key cannot be empty for OSS")
		}
		if c.OSS.AccessSecret == "" {
			return fmt.Errorf("access_secret cannot be empty for OSS")
		}
		if c.OSS.BucketName == "" {
			return fmt.Errorf("bucket_name cannot be empty for OSS")
		}
		if c.OSS.Region == "" {
			return fmt.Errorf("region cannot be empty for OSS")
		}
		if c.OSS.RoleArn == "" {
			return fmt.Errorf("role_arn cannot be empty for OSS")
		}
	case "cos":
		if c.COS.Endpoint == "" {
			return fmt.Errorf("endpoint cannot be empty for COS")
		}
		if c.COS.AccessKey == "" {
			return fmt.Errorf("access_key cannot be empty for COS")
		}
		if c.COS.AccessSecret == "" {
			return fmt.Errorf("access_secret cannot be empty for COS")
		}
		if c.COS.BucketName == "" {
			return fmt.Errorf("bucket_name cannot be empty for COS")
		}
		if c.COS.Region == "" {
			return fmt.Errorf("region cannot be empty for COS")
		}
	case "minio":
		if c.MinIO.Endpoint == "" {
			return fmt.Errorf("endpoint cannot be empty for MinIO")
		}
		if c.MinIO.AccessKey == "" {
			return fmt.Errorf("access_key cannot be empty for MinIO")
		}
		if c.MinIO.AccessSecret == "" {
			return fmt.Errorf("access_secret cannot be empty for MinIO")
		}
		if c.MinIO.BucketName == "" {
			return fmt.Errorf("bucket_name cannot be empty for MinIO")
		}
	default:
		return fmt.Errorf("unsupported storage type: %s", c.Type)
	}

	return nil
}
