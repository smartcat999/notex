package dto

import "notex/model"

// AIProviderResponse 表示AI提供商响应
type AIProviderResponse struct {
	ID          uint   `json:"id"`
	ProviderID  string `json:"providerId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	HasEndpoint bool   `json:"hasEndpoint"`
	IsEnabled   bool   `json:"isEnabled"`
}

// AIModelResponse 表示AI模型响应
type AIModelResponse struct {
	ID          uint   `json:"id"`
	Provider    string `json:"provider"`
	ModelID     string `json:"modelId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsPaid      bool   `json:"isPaid"`
	IsEnabled   bool   `json:"isEnabled"`
}

// AIUserSettingRequest 表示用户AI设置请求
type AIUserSettingRequest struct {
	ProviderID    string                 `json:"providerId" binding:"required"`
	APIKey        string                 `json:"apiKey"`
	Endpoint      string                 `json:"endpoint"`
	EnabledModels map[string]bool        `json:"enabledModels"`
	ModelParams   map[string]interface{} `json:"modelParams"`
}

// AIUserSettingResponse 表示用户AI设置响应
type AIUserSettingResponse struct {
	ID            uint                   `json:"id"`
	ProviderID    string                 `json:"providerId"`
	APIKey        string                 `json:"apiKey"`
	Endpoint      string                 `json:"endpoint"`
	EnabledModels map[string]bool        `json:"enabledModels"`
	ModelParams   map[string]interface{} `json:"modelParams"`
}

// AIDefaultSettingRequest 表示默认AI设置请求
type AIDefaultSettingRequest struct {
	DefaultModel      string `json:"defaultModel,omitempty"`
	DefaultImageModel string `json:"defaultImageModel,omitempty"`
}

// AIDefaultSettingResponse 表示默认AI设置响应
type AIDefaultSettingResponse struct {
	ID                uint   `json:"id"`
	DefaultModel      string `json:"defaultModel"`
	DefaultImageModel string `json:"defaultImageModel"`
}

// AIChatRequest 表示AI聊天请求
type AIChatRequest struct {
	Provider string                 `json:"provider" binding:"required"`
	Messages []map[string]string    `json:"messages" binding:"required"`
	Model    string                 `json:"model" binding:"required"`
	Stream   bool                   `json:"stream"`
	Params   map[string]interface{} `json:"params"`
}

// AITestConnectionRequest 表示测试AI连接请求
type AITestConnectionRequest struct {
	Provider string `json:"provider" binding:"required"`
	APIKey   string `json:"apiKey" binding:"required"`
	Endpoint string `json:"endpoint"`
	Model    string `json:"model" binding:"required"`
}

// AIImageGenerationRequest 表示图像生成请求
type AIImageGenerationRequest struct {
	Provider string                 `json:"provider" binding:"required"`
	Model    string                 `json:"model" binding:"required"`
	Prompt   string                 `json:"prompt" binding:"required"`
	N        int                    `json:"n,omitempty"`
	Size     string                 `json:"size,omitempty"`
	APIKey   string                 `json:"apiKey,omitempty"`
	Endpoint string                 `json:"endpoint,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

// AIImageGenerationResponse 表示图像生成响应
type AIImageGenerationResponse struct {
	Images []string `json:"images"`
}

// AIAvailableModelsResponse 表示可用AI模型响应
type AIAvailableModelsResponse struct {
	Providers []AIProviderWithModels `json:"providers"`
}

// AIProviderWithModels 表示带有模型的AI提供商
type AIProviderWithModels struct {
	ID          uint              `json:"id"`
	ProviderID  string            `json:"providerId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	HasEndpoint bool              `json:"hasEndpoint"`
	Models      []AIModelResponse `json:"models"`
}

// ConvertToAIProviderResponse 将模型转换为响应
func ConvertToAIProviderResponse(provider *model.AIProvider) AIProviderResponse {
	return AIProviderResponse{
		ID:          provider.ID,
		ProviderID:  provider.ProviderID,
		Name:        provider.Name,
		Description: provider.Description,
		HasEndpoint: provider.HasEndpoint,
		IsEnabled:   provider.IsEnabled,
	}
}

// ConvertToAIModelResponse 将模型转换为响应
func ConvertToAIModelResponse(model *model.AIModel) AIModelResponse {
	return AIModelResponse{
		ID:          model.ID,
		Provider:    model.Provider,
		ModelID:     model.ModelID,
		Name:        model.Name,
		Description: model.Description,
		Type:        model.Type,
		IsPaid:      model.IsPaid,
		IsEnabled:   model.IsEnabled,
	}
}

// ConvertToAIUserSettingResponse 将模型转换为响应
func ConvertToAIUserSettingResponse(setting *model.AIUserSetting) AIUserSettingResponse {
	enabledModels := make(map[string]bool)
	if setting.EnabledModels != nil {
		for k, v := range setting.EnabledModels {
			if boolVal, ok := v.(bool); ok {
				enabledModels[k] = boolVal
			}
		}
	}

	return AIUserSettingResponse{
		ID:            setting.ID,
		ProviderID:    setting.ProviderID,
		APIKey:        setting.APIKey,
		Endpoint:      setting.Endpoint,
		EnabledModels: enabledModels,
		ModelParams:   setting.ModelParams,
	}
}

// ConvertToAIDefaultSettingResponse 将模型转换为响应
func ConvertToAIDefaultSettingResponse(setting *model.AIDefaultSetting) AIDefaultSettingResponse {
	return AIDefaultSettingResponse{
		ID:                setting.ID,
		DefaultModel:      setting.DefaultModel,
		DefaultImageModel: setting.DefaultImageModel,
	}
}
