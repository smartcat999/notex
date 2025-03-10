package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository() *VerificationRepository {
	return &VerificationRepository{
		db: database.GetDB(),
	}
}

// Create 创建验证码记录
func (r *VerificationRepository) Create(verification *model.Verification) error {
	return r.db.Create(verification).Error
}

// FindByEmailAndType 根据邮箱和类型查找最新的未使用验证码
func (r *VerificationRepository) FindByEmailAndType(email string, vType model.VerificationType) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.Where("email = ? AND type = ? AND used = ? AND expires_at > NOW()",
		email, vType, false).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// FindByEmailAndCode 根据邮箱和验证码查找记录
func (r *VerificationRepository) FindByEmailAndCode(email, code string, vType model.VerificationType) (*model.Verification, error) {
	var verification model.Verification
	err := r.db.Where("email = ? AND code = ? AND type = ? AND used = ? AND expires_at > NOW()",
		email, code, vType, false).
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// MarkAsUsed 标记验证码为已使用
func (r *VerificationRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.Verification{}).Where("id = ?", id).Update("used", true).Error
}

// DeleteExpired 删除过期的验证码
func (r *VerificationRepository) DeleteExpired() error {
	return r.db.Where("expires_at < NOW()").Delete(&model.Verification{}).Error
}
