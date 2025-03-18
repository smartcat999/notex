package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"notex/api/dto"
	"notex/api/service"
	"notex/middleware"

	"github.com/gin-gonic/gin"
)

// AIHandler 处理AI相关的请求
type AIHandler struct {
	aiService service.AIService
}

// NewAIHandler 创建一个新的AIHandler实例
func NewAIHandler(aiService service.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// getUserIDFromContext 从上下文中获取用户ID
func getUserIDFromContext(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}

	id, ok := userID.(uint)
	if !ok {
		return 0
	}

	return id
}

// RegisterRoutes 注册路由
func (h *AIHandler) RegisterRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai")
	{
		// 公开接口
		ai.GET("/providers", h.GetAllProviders)
		ai.GET("/models", h.GetAllModels)
		ai.GET("/available-models", h.GetAvailableModels)

		// 需要认证的接口
		authenticated := ai.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 用户设置相关
			authenticated.GET("/settings", h.GetUserSettings)
			authenticated.GET("/settings/:providerId", h.GetUserSettingByProvider)
			authenticated.POST("/settings", h.SaveUserSetting)
			authenticated.DELETE("/settings/:providerId", h.DeleteUserSetting)

			// 默认设置相关
			authenticated.GET("/default-setting", h.GetDefaultSetting)
			authenticated.POST("/default-setting", h.SaveDefaultSetting)

			// 聊天相关
			authenticated.POST("/chat", h.HandleAIChat)
			authenticated.POST("/test-connection", h.HandleAITest)
		}
	}
}

// GetAllProviders 获取所有AI提供商
func (h *AIHandler) GetAllProviders(c *gin.Context) {
	providers, err := h.aiService.GetAllProviders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取AI提供商失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"providers": providers})
}

// GetAllModels 获取所有AI模型
func (h *AIHandler) GetAllModels(c *gin.Context) {
	models, err := h.aiService.GetAllModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取AI模型失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"models": models})
}

// GetAvailableModels 获取所有可用的AI模型，按提供商分组
func (h *AIHandler) GetAvailableModels(c *gin.Context) {
	result, err := h.aiService.GetAvailableModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取可用AI模型失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserSettings 获取用户的所有AI设置
func (h *AIHandler) GetUserSettings(c *gin.Context) {
	userID := getUserIDFromContext(c)

	settings, err := h.aiService.GetUserSettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户AI设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"settings": settings})
}

// GetUserSettingByProvider 获取用户特定提供商的AI设置
func (h *AIHandler) GetUserSettingByProvider(c *gin.Context) {
	userID := getUserIDFromContext(c)
	providerID := c.Param("providerId")

	setting, err := h.aiService.GetUserSettingByProvider(userID, providerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到用户AI设置"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// SaveUserSetting 保存用户的AI设置
func (h *AIHandler) SaveUserSetting(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req dto.AIUserSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	setting, err := h.aiService.SaveUserSetting(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// DeleteUserSetting 删除用户的AI设置
func (h *AIHandler) DeleteUserSetting(c *gin.Context) {
	userID := getUserIDFromContext(c)
	providerID := c.Param("providerId")

	err := h.aiService.DeleteUserSetting(userID, providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户AI设置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户AI设置已删除"})
}

// GetDefaultSetting 获取用户的默认AI设置
func (h *AIHandler) GetDefaultSetting(c *gin.Context) {
	userID := getUserIDFromContext(c)

	setting, err := h.aiService.GetDefaultSetting(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到用户默认AI设置"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

// SaveDefaultSetting 保存用户的默认AI设置
func (h *AIHandler) SaveDefaultSetting(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req dto.AIDefaultSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	setting, err := h.aiService.SaveDefaultSetting(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, setting)
}

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
func (h *AIHandler) HandleAIChat(c *gin.Context) {
	userID := getUserIDFromContext(c)

	var req dto.AIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取用户的AI设置
	setting, err := h.aiService.GetUserSettingByProvider(userID, req.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未找到用户AI设置，请先配置API密钥"})
		return
	}

	// 获取API端点
	endpoint := setting.Endpoint
	if endpoint == "" {
		endpoint = getProviderEndpoint(req.Provider)
		if endpoint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提供商或缺少端点"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求体序列化失败"})
		return
	}

	// 创建请求
	proxyReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	// 设置请求头
	headers := getProviderHeaders(req.Provider, setting.APIKey)
	for k, v := range headers {
		proxyReq.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应失败"})
		return
	}

	c.Data(http.StatusOK, "application/json", body)
}

// HandleAITest 处理AI连接测试请求
func (h *AIHandler) HandleAITest(c *gin.Context) {
	var req dto.AITestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 获取API端点
	endpoint := req.Endpoint
	if endpoint == "" {
		endpoint = getProviderEndpoint(req.Provider)
		if endpoint == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提供商或缺少端点"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "请求体序列化失败"})
		return
	}

	// 创建请求
	proxyReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建请求失败"})
		return
	}

	// 设置请求头
	headers := getProviderHeaders(req.Provider, req.APIKey)
	for k, v := range headers {
		proxyReq.Header.Set(k, v)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送请求失败"})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应失败"})
		return
	}

	// 返回响应
	c.Data(http.StatusOK, "application/json", body)
}
