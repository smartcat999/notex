package storage

import (
	"fmt"
	"mime/multipart"
	"notex/pkg/types"
	"path/filepath"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// OSSStorage 阿里云OSS存储实现
type OSSStorage struct {
	config *types.StorageConfig
	client *oss.Client
	bucket *oss.Bucket
}

// NewOSSStorage 创建OSS存储实例
func NewOSSStorage(config *types.StorageConfig) (Storage, error) {
	// 创建OSS客户端
	client, err := oss.New(
		config.OSS.Endpoint,
		config.OSS.AccessKey,
		config.OSS.AccessSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OSS client: %w", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.OSS.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &OSSStorage{
		config: config,
		client: client,
		bucket: bucket,
	}, nil
}

// Upload 上传文件到OSS（仅作为后备方案，不推荐使用）
func (s *OSSStorage) Upload(file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	return nil, fmt.Errorf("server-side upload is not supported for OSS, please use direct upload with STS credentials")
}

// Delete 删除OSS文件
func (s *OSSStorage) Delete(fileURL string) error {
	// 从URL提取对象键
	objectKey := filepath.Base(fileURL)
	return s.bucket.DeleteObject(objectKey)
}

// GetUploadConfig 获取上传配置
func (s *OSSStorage) GetUploadConfig() interface{} {
	// 使用配置的 CDN 域名作为 URL 前缀
	urlPrefix := s.config.OSS.URLPrefix
	if urlPrefix == "" {
		// 如果未配置 CDN，则使用默认的 OSS 域名（不推荐）
		urlPrefix = fmt.Sprintf("https://%s.%s/", s.config.OSS.BucketName, s.config.OSS.Endpoint)
	}

	return &OSSUploadConfig{
		UploadConfig: UploadConfig{
			Type:         StorageTypeOSS,
			MaxSize:      s.config.MaxSize,
			AllowedTypes: s.config.AllowedTypes,
			DirectUpload: true,
			UploadURL:    fmt.Sprintf("https://%s.%s", s.config.OSS.BucketName, s.config.OSS.Endpoint),
		},
		Region:    s.config.OSS.Region,
		Bucket:    s.config.OSS.BucketName,
		Endpoint:  s.config.OSS.Endpoint,
		Headers:   map[string]string{},
		URLPrefix: urlPrefix, // 使用 CDN 域名
	}
}

// GetType 获取存储类型
func (s *OSSStorage) GetType() StorageType {
	return StorageTypeOSS
}

// getSTSToken 获取STS临时凭证
func (s *OSSStorage) getSTSToken(sessionName string) (*sts.Credentials, error) {
	// 创建STS客户端
	client, err := sts.NewClientWithAccessKey(
		s.config.OSS.Region,
		s.config.OSS.AccessKey,
		s.config.OSS.AccessSecret,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create STS client: %w", err)
	}

	// 构造请求
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"
	request.RoleArn = s.config.OSS.RoleArn
	request.RoleSessionName = sessionName
	request.DurationSeconds = "3600"

	// 发送请求
	response, err := client.AssumeRole(request)
	if err != nil {
		return nil, fmt.Errorf("failed to assume role: %w", err)
	}

	return &response.Credentials, nil
}

// GetCredentials 获取上传凭证
func (s *OSSStorage) GetCredentials(filename string, contentType string) (*UploadCredentials, error) {
	// 生成文件路径
	ext := filepath.Ext(filename)
	objectKey := fmt.Sprintf("%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	// 获取STS临时凭证
	sessionName := "oss-upload-" + uuid.New().String()
	stsToken, err := s.getSTSToken(sessionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get STS token: %w", err)
	}

	// 构造上传凭证
	credentials := &UploadCredentials{
		AccessKeyId:     stsToken.AccessKeyId,
		AccessKeySecret: stsToken.AccessKeySecret,
		SecurityToken:   stsToken.SecurityToken,
		FileKey:         objectKey,
		Expires:         time.Now().Add(time.Hour).Unix(),
		Extra: map[string]string{
			"region":   s.config.OSS.Region,
			"bucket":   s.config.OSS.BucketName,
			"endpoint": s.config.OSS.Endpoint,
		},
	}

	return credentials, nil
}
