package response

import (
	"github.com/gin-gonic/gin"
)

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`            // 状态码
	Message string      `json:"message"`         // 消息
	Data    interface{} `json:"data,omitempty"`  // 数据
	Error   string      `json:"error,omitempty"` // 错误信息
}

// Success 返回成功响应
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Error 返回错误响应
func Error(c *gin.Context, code int, message string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}

	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Error:   errMsg,
	})
}
