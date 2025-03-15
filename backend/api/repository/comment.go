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
	result := r.db.
		Preload("User").
		Preload("Parent").
		Preload("Parent.User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

// ListByPostID lists comments for a post with pagination
func (r *CommentRepository) ListByPostID(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// Get total count of parent comments
	if err := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND parent_id IS NULL AND status = ?", postID, "active").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get parent comments with preloaded user information
	err := r.db.
		Preload("User").
		Where("post_id = ? AND parent_id IS NULL AND status = ?", postID, "active").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
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

	// 获取评论列表，包括关联的文章和用户信息
	err := r.db.Model(&model.Comment{}).
		Where("user_id = ?", userID).
		Preload("Post").
		Preload("User").
		Preload("Parent").
		Preload("Parent.User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// ListReplies lists all replies for a comment
func (r *CommentRepository) ListReplies(commentID uint) ([]model.Comment, error) {
	var replies []model.Comment

	err := r.db.
		Preload("User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Where("parent_id = ? AND status = ?", commentID, "active").
		Order("created_at ASC").
		Find(&replies).Error

	if err != nil {
		return nil, err
	}

	return replies, nil
}

// ListByPost lists all comments for a post
func (r *CommentRepository) ListByPost(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64

	// Get total count of top-level comments
	if err := r.db.Model(&model.Comment{}).
		Where("post_id = ? AND parent_id IS NULL", postID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get comments with their replies
	err := r.db.
		Preload("User").
		Preload("Parent").
		Preload("Parent.User").
		Preload("ReplyTo").
		Preload("ReplyTo.User").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		Preload("Children.User").
		Preload("Children.ReplyTo").
		Preload("Children.ReplyTo.User").
		Where("post_id = ? AND parent_id IS NULL", postID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error

	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetReplyCount 获取评论的回复数量
func (r *CommentRepository) GetReplyCount(commentID uint) (int, error) {
	var count int64
	err := r.db.Model(&model.Comment{}).Where("parent_id = ? AND status = ?", commentID, "active").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
