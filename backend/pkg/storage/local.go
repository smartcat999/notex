package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"notex/pkg/types"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	config *types.StorageConfig
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(config *types.StorageConfig) Storage {
	return &LocalStorage{
		config: config,
	}
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	// 创建上传目录
	uploadPath := filepath.Join(s.config.Local.UploadDir, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// 生成文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadPath, filename)

	// 保存文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 生成访问URL（修复路径问题）
	relativePath := strings.TrimPrefix(filePath, s.config.Local.UploadDir)
	relativePath = strings.TrimPrefix(relativePath, "/")
	fileURL := filepath.Join(s.config.Local.URLPrefix, relativePath)
	fileURL = strings.ReplaceAll(fileURL, "\\", "/") // 确保 URL 使用正斜杠

	return &UploadResult{
		URL:      fileURL,
		Filename: header.Filename,
		Size:     header.Size,
		Type:     header.Header.Get("Content-Type"),
	}, nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(fileURL string) error {
	// 将URL转换为本地文件路径
	filePath := strings.Replace(fileURL, s.config.Local.URLPrefix, s.config.Local.UploadDir, 1)
	return os.Remove(filePath)
}

// GetUploadConfig 获取上传配置
func (s *LocalStorage) GetUploadConfig() interface{} {
	return &LocalUploadConfig{
		UploadConfig: UploadConfig{
			Type:         StorageTypeLocal,
			MaxSize:      s.config.MaxSize,
			AllowedTypes: s.config.AllowedTypes,
			DirectUpload: false,
			UploadURL:    "/api/upload/file",
		},
	}
}

// GetType 获取存储类型
func (s *LocalStorage) GetType() StorageType {
	return StorageTypeLocal
}

// GetCredentials 获取上传凭证
func (s *LocalStorage) GetCredentials(filename string, contentType string) (*UploadCredentials, error) {
	// 本地存储不需要凭证
	return &UploadCredentials{}, nil
}
