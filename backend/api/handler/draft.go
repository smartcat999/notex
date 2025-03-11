package handler

import (
	"errors"
	"net/http"
	"notex/api/dto"
	"notex/api/service"
	"notex/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DraftHandler struct {
	draftService *service.DraftService
}

func NewDraftHandler(draftService *service.DraftService) *DraftHandler {
	return &DraftHandler{
		draftService: draftService,
	}
}

// ListDrafts 获取草稿列表
func (h *DraftHandler) ListDrafts(c *gin.Context) {
	var query dto.DraftListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误", err)
		return
	}

	// 从上下文获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	drafts, total, err := h.draftService.ListDrafts(userID.(uint), query.Page, query.PageSize, query.Search)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取草稿列表失败", err)
		return
	}

	// 确保返回空数组而不是null
	if drafts == nil {
		drafts = []dto.DraftResponse{}
	}

	// 直接返回items和total，不需要嵌套在data字段中
	c.JSON(http.StatusOK, gin.H{
		"items": drafts,
		"total": total,
	})
}

// GetDraft 获取草稿详情
func (h *DraftHandler) GetDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	draft, err := h.draftService.GetDraft(uint(id), userID.(uint))
	if err != nil {
		if errors.Is(err, service.ErrDraftNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get draft"})
		return
	}

	// 直接返回草稿数据
	c.JSON(http.StatusOK, draft)
}

// CreateDraft 创建草稿
func (h *DraftHandler) CreateDraft(c *gin.Context) {
	var req dto.CreateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	draft, err := h.draftService.CreateDraft(userID.(uint), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建草稿失败"})
		return
	}

	// 直接返回草稿数据，不包装在 response 结构中
	c.JSON(http.StatusOK, draft)
}

// UpdateDraft 更新草稿
func (h *DraftHandler) UpdateDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	var req dto.UpdateDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameters"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	draft, err := h.draftService.UpdateDraft(uint(id), userID.(uint), req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if errors.Is(err, service.ErrDraftNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update draft"})
		return
	}

	// 直接返回草稿数据
	c.JSON(http.StatusOK, draft)
}

// DeleteDraft 删除草稿
func (h *DraftHandler) DeleteDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的草稿ID", err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未登录", nil)
		return
	}

	err = h.draftService.DeleteDraft(uint(id), userID.(uint))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "删除草稿失败", err)
		return
	}

	response.Success(c, "删除草稿成功", nil)
}

// PublishDraft 发布草稿
func (h *DraftHandler) PublishDraft(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid draft ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	post, err := h.draftService.PublishDraft(uint(id), userID.(uint))
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		if errors.Is(err, service.ErrDraftNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Draft not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish draft"})
		return
	}

	// 直接返回发布后的文章数据
	c.JSON(http.StatusOK, post)
}
