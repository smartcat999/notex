package handler

import (
	"net/http"
	"notex/api/dto"
	"notex/api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		service: postService,
	}
}

// ListPosts 获取文章列表
func (h *PostHandler) ListPosts(c *gin.Context) {
	var query dto.PostListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理当前用户的文章过滤
	if query.User == "current" {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
			return
		}
		query.UserID = userID.(uint)
	}

	posts, total, err := h.service.ListPosts(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": posts,
		"total": total,
	})
}

// GetPost 获取文章详情
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	post, err := h.service.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 增加浏览量
	go h.service.IncrementViews(uint(id))

	c.JSON(http.StatusOK, post)
}

// CreatePost 创建文章
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID.(uint)

	post, err := h.service.CreatePost(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// UpdatePost 更新文章
func (h *PostHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.service.UpdatePost(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost 删除文章
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeletePost(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetRecentPosts 获取最新文章列表
func (h *PostHandler) GetRecentPosts(c *gin.Context) {
	limit := 5 // 默认获取5篇
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	posts, err := h.service.GetRecentPosts(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": posts,
	})
}

// ListPublicPosts 获取公开文章列表（匿名访问）
func (h *PostHandler) ListPublicPosts(c *gin.Context) {
	var query dto.PostListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 强制设置状态为已发布
	query.Status = "published"

	posts, total, err := h.service.ListPosts(&query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": posts,
		"total": total,
	})
}

// GetArchives 获取文章归档列表
func (h *PostHandler) GetArchives(c *gin.Context) {
	archives, err := h.service.GetArchives()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": archives,
	})
}

// GetPostsByArchive 获取指定归档日期的文章列表
func (h *PostHandler) GetPostsByArchive(c *gin.Context) {
	yearMonth := c.Param("yearMonth") // 格式：YYYY-MM
	if yearMonth == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year-month format"})
		return
	}

	posts, err := h.service.GetPostsByArchive(yearMonth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": posts,
	})
}
