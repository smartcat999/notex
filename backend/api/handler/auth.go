package handler

import (
	"net/http"
	"notex/api/dto"
	"notex/api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service     *service.AuthService
	postService *service.PostService
}

func NewAuthHandler(authService *service.AuthService, postService *service.PostService) *AuthHandler {
	return &AuthHandler{
		service:     authService,
		postService: postService,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Register(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registration successful"})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken 刷新访问令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangePassword(userID.(uint), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}

// GetProfile 获取用户个人信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权访问"})
		return
	}

	user, err := h.service.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateProfile 更新用户个人信息
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未授权访问"})
		return
	}

	var input struct {
		Username string `json:"username" binding:"required,min=3,max=20"`
		Bio      string `json:"bio" binding:"max=500"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求参数错误"})
		return
	}

	// 检查用户名是否已被使用（排除当前用户）
	exists, err := h.service.IsUsernameExistsExcept(input.Username, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "检查用户名失败"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名已被使用"})
		return
	}

	// 获取当前用户信息
	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取用户信息失败"})
		return
	}

	// 更新用户信息
	user.Username = input.Username
	user.Bio = input.Bio
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}

	if err := h.service.UpdateUserProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新用户信息失败"})
		return
	}

	// 获取更新后的完整用户信息
	updatedProfile, err := h.service.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取更新后的用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"user":    updatedProfile,
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// GetPublicProfile 获取用户公开个人主页信息
func (h *AuthHandler) GetPublicProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	user, err := h.service.GetUserProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// GetUserHome 获取用户主页信息（包含用户资料和文章列表）
func (h *AuthHandler) GetUserHome(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "无效的用户ID"})
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	// 获取用户信息
	user, err := h.service.GetUserProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	// 创建查询对象
	query := &dto.PostListQuery{
		UserID:   uint(userID),
		Page:     page,
		PageSize: pageSize,
		Status:   "published", // 强制设置状态为已发布
	}

	// 获取用户文章列表
	posts, total, err := h.postService.ListPosts(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取文章列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
		"posts": gin.H{
			"items": posts,
			"total": total,
		},
	})
}
