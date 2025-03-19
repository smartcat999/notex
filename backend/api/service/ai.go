package service

import (
	"errors"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"

	"gorm.io/gorm"
)

// AIService 定义AI相关的业务逻辑接口
type AIService interface {
	// 提供商和模型相关
	GetAllProviders() ([]dto.AIProviderResponse, error)
	GetProviderByID(providerID string) (*dto.AIProviderResponse, error)
	GetModelsByProvider(provider string) ([]dto.AIModelResponse, error)
	GetAllModels() ([]dto.AIModelResponse, error)
	GetModelsByType(modelType string) ([]dto.AIModelResponse, error)
	GetAvailableModels() (*dto.AIAvailableModelsResponse, error)

	// 用户设置相关
	GetUserSettings(userID uint) ([]dto.AIUserSettingResponse, error)
	GetUserSettingByProvider(userID uint, providerID string) (*dto.AIUserSettingResponse, error)
	SaveUserSetting(userID uint, req *dto.AIUserSettingRequest) (*dto.AIUserSettingResponse, error)
	DeleteUserSetting(userID uint, providerID string) error

	// 默认设置相关
	GetDefaultSetting(userID uint) (*dto.AIDefaultSettingResponse, error)
	SaveDefaultSetting(userID uint, req *dto.AIDefaultSettingRequest) (*dto.AIDefaultSettingResponse, error)
}

// AIServiceImpl 实现AIService接口
type AIServiceImpl struct {
	repo repository.AIRepository
}

// NewAIService 创建一个新的AIService实例
func NewAIService() AIService {
	return &AIServiceImpl{
		repo: repository.NewAIRepository(),
	}
}

// GetAllProviders 获取所有AI提供商
func (s *AIServiceImpl) GetAllProviders() ([]dto.AIProviderResponse, error) {
	providers, err := s.repo.GetAllProviders()
	if err != nil {
		return nil, err
	}

	result := make([]dto.AIProviderResponse, len(providers))
	for i, provider := range providers {
		result[i] = dto.ConvertToAIProviderResponse(&provider)
	}

	return result, nil
}

// GetProviderByID 根据ID获取AI提供商
func (s *AIServiceImpl) GetProviderByID(providerID string) (*dto.AIProviderResponse, error) {
	provider, err := s.repo.GetProviderByID(providerID)
	if err != nil {
		return nil, err
	}

	result := dto.ConvertToAIProviderResponse(provider)
	return &result, nil
}

// GetModelsByProvider 获取指定提供商的所有模型
func (s *AIServiceImpl) GetModelsByProvider(provider string) ([]dto.AIModelResponse, error) {
	models, err := s.repo.GetModelsByProvider(provider)
	if err != nil {
		return nil, err
	}

	result := make([]dto.AIModelResponse, len(models))
	for i, model := range models {
		result[i] = dto.ConvertToAIModelResponse(&model)
	}

	return result, nil
}

// GetAllModels 获取所有AI模型
func (s *AIServiceImpl) GetAllModels() ([]dto.AIModelResponse, error) {
	models, err := s.repo.GetAllModels()
	if err != nil {
		return nil, err
	}

	result := make([]dto.AIModelResponse, len(models))
	for i, model := range models {
		result[i] = dto.ConvertToAIModelResponse(&model)
	}

	return result, nil
}

// GetModelsByType 获取指定类型的所有模型
func (s *AIServiceImpl) GetModelsByType(modelType string) ([]dto.AIModelResponse, error) {
	models, err := s.repo.GetModelsByType(modelType)
	if err != nil {
		return nil, err
	}

	result := make([]dto.AIModelResponse, len(models))
	for i, model := range models {
		result[i] = dto.ConvertToAIModelResponse(&model)
	}

	return result, nil
}

// GetAvailableModels 获取所有可用的AI模型，按提供商分组
func (s *AIServiceImpl) GetAvailableModels() (*dto.AIAvailableModelsResponse, error) {
	providers, err := s.repo.GetAllProviders()
	if err != nil {
		return nil, err
	}

	result := &dto.AIAvailableModelsResponse{
		Providers: make([]dto.AIProviderWithModels, 0, len(providers)),
	}

	for _, provider := range providers {
		models, err := s.repo.GetModelsByProvider(provider.ProviderID)
		if err != nil {
			return nil, err
		}

		modelResponses := make([]dto.AIModelResponse, len(models))
		for i, model := range models {
			modelResponses[i] = dto.ConvertToAIModelResponse(&model)
		}

		result.Providers = append(result.Providers, dto.AIProviderWithModels{
			ID:          provider.ID,
			ProviderID:  provider.ProviderID,
			Name:        provider.Name,
			Description: provider.Description,
			HasEndpoint: provider.HasEndpoint,
			Models:      modelResponses,
		})
	}

	return result, nil
}

// GetUserSettings 获取用户的所有AI设置
func (s *AIServiceImpl) GetUserSettings(userID uint) ([]dto.AIUserSettingResponse, error) {
	settings, err := s.repo.GetUserSettings(userID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.AIUserSettingResponse, len(settings))
	for i, setting := range settings {
		result[i] = dto.ConvertToAIUserSettingResponse(&setting)
	}

	return result, nil
}

// GetUserSettingByProvider 获取用户特定提供商的AI设置
func (s *AIServiceImpl) GetUserSettingByProvider(userID uint, providerID string) (*dto.AIUserSettingResponse, error) {
	setting, err := s.repo.GetUserSettingByProvider(userID, providerID)
	if err != nil {
		return nil, err
	}

	result := dto.ConvertToAIUserSettingResponse(setting)
	return &result, nil
}

// SaveUserSetting 保存用户的AI设置
func (s *AIServiceImpl) SaveUserSetting(userID uint, req *dto.AIUserSettingRequest) (*dto.AIUserSettingResponse, error) {
	// 验证提供商是否存在
	_, err := s.repo.GetProviderByID(req.ProviderID)
	if err != nil {
		return nil, errors.New("提供商不存在")
	}

	// 转换为模型
	enabledModels := make(model.JSONMap)
	for k, v := range req.EnabledModels {
		enabledModels[k] = v
	}

	modelParams := make(model.JSONMap)
	for k, v := range req.ModelParams {
		modelParams[k] = v
	}

	setting := &model.AIUserSetting{
		UserID:        userID,
		ProviderID:    req.ProviderID,
		APIKey:        req.APIKey,
		Endpoint:      req.Endpoint,
		EnabledModels: enabledModels,
		ModelParams:   modelParams,
	}

	// 保存设置
	err = s.repo.SaveUserSetting(setting)
	if err != nil {
		return nil, err
	}

	// 返回响应
	result := dto.ConvertToAIUserSettingResponse(setting)
	return &result, nil
}

// DeleteUserSetting 删除用户的AI设置
func (s *AIServiceImpl) DeleteUserSetting(userID uint, providerID string) error {
	return s.repo.DeleteUserSetting(userID, providerID)
}

// GetDefaultSetting 获取用户的默认AI设置
func (s *AIServiceImpl) GetDefaultSetting(userID uint) (*dto.AIDefaultSettingResponse, error) {
	setting, err := s.repo.GetDefaultSetting(userID)
	if err != nil {
		return nil, err
	}

	result := dto.ConvertToAIDefaultSettingResponse(setting)
	return &result, nil
}

// SaveDefaultSetting 保存用户的默认AI设置
func (s *AIServiceImpl) SaveDefaultSetting(userID uint, req *dto.AIDefaultSettingRequest) (*dto.AIDefaultSettingResponse, error) {
	// 获取当前设置
	currentSetting, err := s.repo.GetDefaultSetting(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 创建或更新设置
	setting := &model.AIDefaultSetting{
		UserID: userID,
	}

	// 如果有现有设置，则保留ID
	if currentSetting != nil {
		setting.ID = currentSetting.ID
		setting.DefaultModel = currentSetting.DefaultModel
		setting.DefaultImageModel = currentSetting.DefaultImageModel
	}

	// 验证并更新默认文本模型
	if req.DefaultModel != "" {
		_, err := s.repo.GetModelByID(req.DefaultModel)
		if err != nil {
			return nil, errors.New("指定的默认文本模型不存在")
		}
		setting.DefaultModel = req.DefaultModel
	}

	// 验证并更新默认图像模型
	if req.DefaultImageModel != "" {
		_, err := s.repo.GetModelByID(req.DefaultImageModel)
		if err != nil {
			return nil, errors.New("指定的默认图像模型不存在")
		}
		setting.DefaultImageModel = req.DefaultImageModel
	}

	// 保存设置
	err = s.repo.SaveDefaultSetting(setting)
	if err != nil {
		return nil, err
	}

	// 返回响应
	result := dto.ConvertToAIDefaultSettingResponse(setting)
	return &result, nil
}
