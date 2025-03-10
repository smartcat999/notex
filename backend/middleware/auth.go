package middleware

import (
	"encoding/json"
	"net/http"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
	"notex/pkg/auth"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// getAuditLogRepo 获取审计日志仓库实例
func getAuditLogRepo() *repository.AuditLogRepository {
	return repository.NewAuditLogRepository()
}

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// 记录访问日志
		go logAccess(c, claims)

		c.Next()
	}
}

// RequireRoles 检查用户是否具有指定角色之一
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, r := range roles {
			if userRole == r {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "permission denied"})
		c.Abort()
	}
}

// RequireAdmin 管理员权限中间件
func RequireAdmin() gin.HandlerFunc {
	return RequireRoles(model.RoleAdmin)
}

// RequireEditor 编辑权限中间件
func RequireEditor() gin.HandlerFunc {
	return RequireRoles(model.RoleAdmin, model.RoleEditor)
}

// AuditLog 审计日志中间件
func AuditLog(action, resource string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")

		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 记录审计日志
		status := "success"
		if len(c.Errors) > 0 {
			status = "failed"
		}

		details := map[string]interface{}{
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"query":       c.Request.URL.RawQuery,
			"duration_ms": time.Since(startTime).Milliseconds(),
		}
		detailsJSON, _ := json.Marshal(details)

		log := &model.AuditLog{
			UserID:     userID.(uint),
			Username:   username.(string),
			Action:     action,
			Resource:   resource,
			ResourceID: c.Param("id"),
			Details:    string(detailsJSON),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			Status:     status,
		}

		if len(c.Errors) > 0 {
			log.Error = c.Errors.String()
		}

		go getAuditLogRepo().Create(log)
	}
}

// logAccess 记录访问日志
func logAccess(c *gin.Context, claims *dto.TokenClaims) {
	details := map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
		"query":  c.Request.URL.RawQuery,
	}
	detailsJSON, _ := json.Marshal(details)

	log := &model.AuditLog{
		UserID:    claims.UserID,
		Username:  claims.Username,
		Action:    "access",
		Resource:  "api",
		Details:   string(detailsJSON),
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Status:    "success",
	}

	getAuditLogRepo().Create(log)
}
