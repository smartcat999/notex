package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	_ "image/png"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"

	"notex/config"
	"notex/pkg/response"
)

type UploadHandler struct {
	cfg *config.FileStorageConfig
}

func NewUploadHandler(cfg *config.FileStorageConfig) *UploadHandler {
	return &UploadHandler{cfg: cfg}
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
		response.Error(c, http.StatusBadRequest, "Failed to get file", err)
		return
	}
	defer file.Close()

	// 检查文件大小
	if header.Size > h.cfg.MaxSize {
		response.Error(c, http.StatusBadRequest, "File too large", nil)
		return
	}

	// 检查文件类型
	contentType := header.Header.Get("Content-Type")
	if !h.isAllowedType(contentType) {
		response.Error(c, http.StatusBadRequest, "File type not allowed", nil)
		return
	}

	// 创建上传目录
	uploadPath := filepath.Join(h.cfg.UploadDir, time.Now().Format("2006/01/02"))
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create upload directory", err)
		return
	}

	// 生成文件名
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(uploadPath, filename)

	// 保存文件
	if err := h.saveFile(file, filePath); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to save file", err)
		return
	}

	// 生成缩略图（如果是图片）
	var thumbnailURL string
	if h.isImage(contentType) {
		if thumbnail, err := h.createThumbnail(filePath); err == nil {
			thumbnailURL = strings.Replace(thumbnail, h.cfg.UploadDir, h.cfg.URLPrefix, 1)
		}
	}

	// 生成响应
	fileURL := strings.Replace(filePath, h.cfg.UploadDir, h.cfg.URLPrefix, 1)
	resp := UploadResponse{
		URL:          fileURL,
		ThumbnailURL: thumbnailURL,
		Filename:     header.Filename,
		Size:         header.Size,
		Type:         contentType,
	}

	response.Success(c, "File uploaded successfully", resp)
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
