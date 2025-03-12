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

	"notex/pkg/storage"
	"notex/pkg/types"
)

// UploadHandler 上传处理器
type UploadHandler struct {
	storage storage.Storage
	config  *types.StorageConfig
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(storage storage.Storage, config *types.StorageConfig) *UploadHandler {
	return &UploadHandler{
		storage: storage,
		config:  config,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	result, err := h.storage.Upload(file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUploadConfig 获取上传配置
func (h *UploadHandler) GetUploadConfig(c *gin.Context) {
	config := h.storage.GetUploadConfig()
	c.JSON(http.StatusOK, config)
}

// isAllowedType 检查文件类型是否允许
func (h *UploadHandler) isAllowedType(contentType string) bool {
	for _, allowed := range h.config.AllowedTypes {
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
	thumbnail := resize.Thumbnail(uint(h.config.ThumbnailSize), uint(h.config.ThumbnailSize), img, resize.Lanczos3)

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

// GetCredentials 获取上传凭证
func (h *UploadHandler) GetCredentials(c *gin.Context) {
	filename := c.Query("filename")
	contentType := c.Query("contentType")

	if filename == "" || contentType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename and contentType are required"})
		return
	}

	// 检查文件类型是否允许
	if !h.isAllowedType(contentType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file type not allowed"})
		return
	}

	credentials, err := h.storage.GetCredentials(filename, contentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, credentials)
}
