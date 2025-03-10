package dto

// UserListRequest 用户列表请求
type UserListRequest struct {
	Page     int    `form:"page" binding:"required,min=1"`
	PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
	Role     string `form:"role"`
	Status   string `form:"status"`
	Search   string `form:"search"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Role   *string `json:"role"`
	Status *string `json:"status"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	Status        string `json:"status"`
	EmailVerified bool   `json:"email_verified"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total"`
	Items []UserResponse `json:"items"`
}

// AuditLogResponse 审计日志响应
type AuditLogResponse struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
	ResourceID string `json:"resource_id"`
	Details    string `json:"details"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Status     string `json:"status"`
	Error      string `json:"error"`
	CreatedAt  string `json:"created_at"`
}
