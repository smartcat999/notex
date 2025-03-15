package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		db: database.GetDB(),
	}
}

// Create 创建通知
func (r *NotificationRepository) Create(notification *model.Notification) error {
	return r.db.Create(notification).Error
}

// MarkAsRead 将通知标记为已读
func (r *NotificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&model.Notification{}).Where("id = ?", id).Update("read", true).Error
}

// MarkAllAsRead 将用户的所有通知标记为已读
func (r *NotificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Update("read", true).Error
}

// Delete 删除通知
func (r *NotificationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Notification{}, id).Error
}

// FindByID 根据ID查找通知
func (r *NotificationRepository) FindByID(id uint) (*model.Notification, error) {
	var notification model.Notification
	err := r.db.Preload("Actor").Preload("Post").Preload("Comment").First(&notification, id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

// ListByUser 获取用户的通知列表
func (r *NotificationRepository) ListByUser(userID uint, page, pageSize int, unread *bool, notificationType string) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)

	// 添加未读过滤
	if unread != nil {
		query = query.Where("read = ?", !*unread)
	}

	// 添加类型过滤
	if notificationType != "" {
		query = query.Where("type = ?", notificationType)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Preload("Actor").Preload("Post").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notifications).Error

	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

// GetUnreadCount 获取用户未读通知数量
func (r *NotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&count).Error
	return count, err
}
