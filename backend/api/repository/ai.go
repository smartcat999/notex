package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

// AIRepository 定义AI相关的数据库操作接口
type AIRepository interface {
	// 提供商和模型相关
	GetAllProviders() ([]model.AIProvider, error)
	GetProviderByID(providerID string) (*model.AIProvider, error)
	GetModelsByProvider(provider string) ([]model.AIModel, error)
	GetAllModels() ([]model.AIModel, error)
	GetModelByID(modelID string) (*model.AIModel, error)

	// 用户设置相关
	GetUserSettings(userID uint) ([]model.AIUserSetting, error)
	GetUserSettingByProvider(userID uint, providerID string) (*model.AIUserSetting, error)
	SaveUserSetting(setting *model.AIUserSetting) error
	DeleteUserSetting(userID uint, providerID string) error

	// 默认设置相关
	GetDefaultSetting(userID uint) (*model.AIDefaultSetting, error)
	SaveDefaultSetting(setting *model.AIDefaultSetting) error
}

// AIRepositoryImpl 实现AIRepository接口
type AIRepositoryImpl struct {
	db *gorm.DB
}

// NewAIRepository 创建一个新的AIRepository实例
func NewAIRepository() AIRepository {
	return &AIRepositoryImpl{
		db: database.GetDB(),
	}
}

// GetAllProviders 获取所有AI提供商
func (r *AIRepositoryImpl) GetAllProviders() ([]model.AIProvider, error) {
	var providers []model.AIProvider
	err := r.db.Where("is_enabled = ?", true).Find(&providers).Error
	return providers, err
}

// GetProviderByID 根据ID获取AI提供商
func (r *AIRepositoryImpl) GetProviderByID(providerID string) (*model.AIProvider, error) {
	var provider model.AIProvider
	err := r.db.Where("provider_id = ? AND is_enabled = ?", providerID, true).First(&provider).Error
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// GetModelsByProvider 获取指定提供商的所有模型
func (r *AIRepositoryImpl) GetModelsByProvider(provider string) ([]model.AIModel, error) {
	var models []model.AIModel
	err := r.db.Where("provider = ? AND is_enabled = ?", provider, true).Find(&models).Error
	return models, err
}

// GetAllModels 获取所有AI模型
func (r *AIRepositoryImpl) GetAllModels() ([]model.AIModel, error) {
	var models []model.AIModel
	err := r.db.Where("is_enabled = ?", true).Find(&models).Error
	return models, err
}

// GetModelByID 根据ID获取AI模型
func (r *AIRepositoryImpl) GetModelByID(modelID string) (*model.AIModel, error) {
	var model model.AIModel
	err := r.db.Where("model_id = ? AND is_enabled = ?", modelID, true).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// GetUserSettings 获取用户的所有AI设置
func (r *AIRepositoryImpl) GetUserSettings(userID uint) ([]model.AIUserSetting, error) {
	var settings []model.AIUserSetting
	err := r.db.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

// GetUserSettingByProvider 获取用户特定提供商的AI设置
func (r *AIRepositoryImpl) GetUserSettingByProvider(userID uint, providerID string) (*model.AIUserSetting, error) {
	var setting model.AIUserSetting
	err := r.db.Where("user_id = ? AND provider_id = ?", userID, providerID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SaveUserSetting 保存用户的AI设置
func (r *AIRepositoryImpl) SaveUserSetting(setting *model.AIUserSetting) error {
	// 检查是否已存在
	var existingSetting model.AIUserSetting
	result := r.db.Where("user_id = ? AND provider_id = ?", setting.UserID, setting.ProviderID).First(&existingSetting)

	if result.Error == nil {
		// 更新现有记录
		setting.ID = existingSetting.ID
		return r.db.Save(setting).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(setting).Error
	}

	return result.Error
}

// DeleteUserSetting 删除用户的AI设置
func (r *AIRepositoryImpl) DeleteUserSetting(userID uint, providerID string) error {
	return r.db.Where("user_id = ? AND provider_id = ?", userID, providerID).Delete(&model.AIUserSetting{}).Error
}

// GetDefaultSetting 获取用户的默认AI设置
func (r *AIRepositoryImpl) GetDefaultSetting(userID uint) (*model.AIDefaultSetting, error) {
	var setting model.AIDefaultSetting
	err := r.db.Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

// SaveDefaultSetting 保存用户的默认AI设置
func (r *AIRepositoryImpl) SaveDefaultSetting(setting *model.AIDefaultSetting) error {
	// 检查是否已存在
	var existingSetting model.AIDefaultSetting
	result := r.db.Where("user_id = ?", setting.UserID).First(&existingSetting)

	if result.Error == nil {
		// 更新现有记录
		setting.ID = existingSetting.ID
		return r.db.Save(setting).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		// 创建新记录
		return r.db.Create(setting).Error
	}

	return result.Error
}
