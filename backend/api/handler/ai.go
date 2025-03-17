package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AIProvider 表示AI提供商
type AIProvider struct {
	Provider string                 `json:"provider"`
	Messages []map[string]string    `json:"messages"`
	Model    string                 `json:"model"`
	Stream   bool                   `json:"stream"`
	ApiKey   string                 `json:"apiKey"`
	Endpoint string                 `json:"endpoint,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// 获取AI提供商的API端点
func getProviderEndpoint(provider string) string {
	endpoints := map[string]string{
		"openai":    "https://api.openai.com/v1/chat/completions",
		"anthropic": "https://api.anthropic.com/v1/messages",
		"google":    "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
		"deepseek":  "https://api.deepseek.com/v1/chat/completions",
	}
	return endpoints[provider]
}

// 获取AI提供商的请求头
func getProviderHeaders(provider, apiKey string) map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	switch provider {
	case "openai":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	case "anthropic":
		headers["x-api-key"] = apiKey
	case "google":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	case "deepseek":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	case "custom":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", apiKey)
	}

	return headers
}

// 格式化消息以适应不同提供商的格式
func formatMessages(messages []map[string]string, provider string) interface{} {
	switch provider {
	case "anthropic":
		formatted := make([]map[string]interface{}, len(messages))
		for i, msg := range messages {
			formatted[i] = map[string]interface{}{
				"role":    msg["role"],
				"content": msg["content"],
			}
		}
		return formatted
	case "google":
		formatted := make([]map[string]interface{}, len(messages))
		for i, msg := range messages {
			formatted[i] = map[string]interface{}{
				"role": msg["role"],
				"parts": []map[string]string{
					{"text": msg["content"]},
				},
			}
		}
		return formatted
	default:
		return messages
	}
}

// HandleAIChat 处理AI聊天请求
func HandleAIChat(c *gin.Context) {
	var req AIProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 获取API端点
	endpoint := req.Endpoint
	if endpoint == "" {
		endpoint = getProviderEndpoint(req.Provider)
		if endpoint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider or missing endpoint"})
			return
		}
	}

	// 准备请求体
	requestBody := map[string]interface{}{
		"messages": formatMessages(req.Messages, req.Provider),
		"model":    req.Model,
		"stream":   req.Stream,
	}

	// 添加其他参数
	if req.Params != nil {
		for k, v := range req.Params {
			requestBody[k] = v
		}
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	// 创建请求
	proxyReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 设置请求头
	headers := getProviderHeaders(req.Provider, req.ApiKey)
	for k, v := range headers {
		proxyReq.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// 设置响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // 禁用 Nginx 缓冲

	// 如果是流式响应，直接转发
	if req.Stream {
		// 创建一个缓冲区
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				// 写入数据
				if _, err := c.Writer.Write(buf[:n]); err != nil {
					return
				}
				// 刷新响应
				c.Writer.Flush()
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				return
			}
		}
		return
	}

	// 非流式响应，读取并转发
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

// HandleAITest 处理AI连接测试请求
func HandleAITest(c *gin.Context) {
	var req AIProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// 获取API端点
	endpoint := req.Endpoint
	if endpoint == "" {
		endpoint = getProviderEndpoint(req.Provider)
		if endpoint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid provider or missing endpoint"})
			return
		}
	}

	// 准备测试消息
	testMessage := map[string]string{
		"role":    "user",
		"content": "Hello",
	}

	// 准备请求体
	requestBody := map[string]interface{}{
		"messages": []map[string]string{testMessage},
		"model":    req.Model,
		"stream":   false,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
		return
	}

	// 创建请求
	proxyReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// 设置请求头
	headers := getProviderHeaders(req.Provider, req.ApiKey)
	for k, v := range headers {
		proxyReq.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// 返回响应
	c.Data(http.StatusOK, "application/json", body)
}
