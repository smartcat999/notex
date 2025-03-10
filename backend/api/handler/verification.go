package handler

import (
	"net/http"
	"notex/api/dto"
	"notex/api/service"
	"notex/model"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	service *service.VerificationService
}

func NewVerificationHandler() *VerificationHandler {
	return &VerificationHandler{
		service: service.NewVerificationService(),
	}
}

// SendEmailVerification 发送邮箱验证码
func (h *VerificationHandler) SendEmailVerification(c *gin.Context) {
	var req dto.SendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SendVerification(&req, model.VerificationTypeEmail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent"})
}

// SendPasswordReset 发送密码重置验证码
func (h *VerificationHandler) SendPasswordReset(c *gin.Context) {
	var req dto.SendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SendVerification(&req, model.VerificationTypePasswordReset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset code sent"})
}

// VerifyEmail 验证邮箱
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.VerifyEmail(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email verified"})
}

// UpdateEmail 更新邮箱
func (h *VerificationHandler) UpdateEmail(c *gin.Context) {
	var req dto.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.UpdateEmail(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "email update initiated"})
}

// ResetPassword 重置密码
func (h *VerificationHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successful"})
}
