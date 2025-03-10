package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository() *AuditLogRepository {
	return &AuditLogRepository{
		db: database.GetDB(),
	}
}

// Create 创建审计日志
func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

// FindAll 查找审计日志
func (r *AuditLogRepository) FindAll(page, pageSize int, userID uint, action, resource string) ([]*model.AuditLog, int64, error) {
	var logs []*model.AuditLog
	var total int64

	query := r.db.Model(&model.AuditLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// DeleteOldLogs 删除旧日志（保留最近30天）
func (r *AuditLogRepository) DeleteOldLogs() error {
	return r.db.Where("created_at < NOW() - INTERVAL '30 days'").Delete(&model.AuditLog{}).Error
}
