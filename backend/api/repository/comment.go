package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		db: database.GetDB(),
	}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// FindByID 根据ID查找评论
func (r *CommentRepository) FindByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListByPostID 获取文章的评论列表
func (r *CommentRepository) ListByPostID(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	query := r.db.Model(&model.Comment{}).Where("post_id = ? AND status = ?", postID, "active")

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Preload("User").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// CountByUserID 获取用户的评论数量
func (r *CommentRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// ListByUser 获取用户的评论列表
func (r *CommentRepository) ListByUser(userID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// 获取总数
	if err := r.db.Model(&model.Comment{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取评论列表，包括关联的文章信息
	err := r.db.Model(&model.Comment{}).
		Where("user_id = ?", userID).
		Preload("Post").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}
