package handler

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"

	"notex/config"
	"notex/pkg/response"
	"notex/pkg/storage"
)

type UploadHandler struct {
	storage storage.Storage
	cfg     *config.FileStorageConfig
}

func NewUploadHandler(storage storage.Storage, cfg *config.FileStorageConfig) *UploadHandler {
	return &UploadHandler{
		storage: storage,
		cfg:     cfg,
	}
}

// UploadResponse 文件上传响应
type UploadResponse struct {
	URL          string `json:"url"`           // 文件URL
	ThumbnailURL string `json:"thumbnail_url"` // 缩略图URL（仅图片）
	Filename     string `json:"filename"`      // 文件名
	Size         int64  `json:"size"`          // 文件大小
	Type         string `json:"type"`          // 文件类型
}

// Upload 处理文件上传
func (h *UploadHandler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	// 使用存储接口上传文件
	result, err := h.storage.Upload(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	// 直接返回上传结果
	c.JSON(http.StatusOK, result)
}

// GetUploadConfig 获取上传配置
func (h *UploadHandler) GetUploadConfig(c *gin.Context) {
	config := h.storage.GetUploadConfig()
	response.Success(c, "Upload config retrieved", config)
}

// isAllowedType 检查文件类型是否允许
func (h *UploadHandler) isAllowedType(contentType string) bool {
	for _, allowed := range h.cfg.AllowedTypes {
		if contentType == allowed {
			return true
		}
	}
	return false
}

// isImage 检查是否为图片
func (h *UploadHandler) isImage(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// saveFile 保存文件
func (h *UploadHandler) saveFile(file multipart.File, filepath string) error {
	dst, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return err
	}

	return nil
}

// createThumbnail 创建缩略图
func (h *UploadHandler) createThumbnail(filePath string) (string, error) {
	// 打开原图
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 解码图片
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// 生成缩略图
	thumbnail := resize.Thumbnail(uint(h.cfg.ThumbnailSize), uint(h.cfg.ThumbnailSize), img, resize.Lanczos3)

	// 生成缩略图文件名
	dir := filepath.Dir(filePath)
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)
	baseWithoutExt := strings.TrimSuffix(base, ext)
	thumbnailPath := filepath.Join(dir, baseWithoutExt+"_thumb"+ext)

	// 保存缩略图
	out, err := os.Create(thumbnailPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 根据文件扩展名选择编码器
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(out, thumbnail, nil)
	case ".png":
		err = png.Encode(out, thumbnail)
	default:
		err = jpeg.Encode(out, thumbnail, nil)
	}

	if err != nil {
		return "", err
	}

	return thumbnailPath, nil
}
