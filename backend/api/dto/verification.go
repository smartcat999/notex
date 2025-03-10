package dto

// SendVerificationRequest 发送验证码请求
type SendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// UpdateEmailRequest 更新邮箱请求
type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
