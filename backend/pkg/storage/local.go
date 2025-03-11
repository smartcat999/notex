package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	config *StorageConfig
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(config *StorageConfig) Storage {
	return &LocalStorage{
		config: config,
	}
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(file multipart.File, header *multipart.FileHeader) (*UploadResult, error) {
	// 创建上传目录
	uploadPath := filepath.Join(s.config.LocalConfig.UploadDir, time.Now().Format("2006/01/02"))
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
	relativePath := strings.TrimPrefix(filePath, s.config.LocalConfig.UploadDir)
	relativePath = strings.TrimPrefix(relativePath, "/")
	fileURL := filepath.Join(s.config.URLPrefix, relativePath)
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
	filePath := strings.Replace(fileURL, s.config.URLPrefix, s.config.LocalConfig.UploadDir, 1)
	return os.Remove(filePath)
}

// GetUploadConfig 获取上传配置
func (s *LocalStorage) GetUploadConfig() map[string]interface{} {
	return map[string]interface{}{
		"type":          "local",
		"upload_url":    "/api/upload", // 本地上传使用后端API
		"max_size":      s.config.MaxSize,
		"allowed_types": s.config.AllowedTypes,
	}
}
