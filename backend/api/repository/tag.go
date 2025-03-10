package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository() *TagRepository {
	return &TagRepository{
		db: database.GetDB(),
	}
}

// Create 创建标签
func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

// Update 更新标签
func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

// Delete 删除标签
func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

// FindByID 根据ID查找标签
func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// List 获取标签列表
func (r *TagRepository) List(page, pageSize int, search string) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64

	query := r.db.Model(&model.Tag{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&tags).Error

	if err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

// GetPostCount 获取标签下的文章数量
func (r *TagRepository) GetPostCount(tagID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).
		Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Where("post_tags.tag_id = ?", tagID).
		Count(&count).Error
	return count, err
}

// ListByPostCount 获取热门标签（按文章数量排序）
func (r *TagRepository) ListByPostCount(limit int) ([]model.Tag, error) {
	tags := make([]model.Tag, 0)
	err := r.db.Model(&model.Tag{}).
		Select("tags.*, COUNT(post_tags.post_id) as post_count").
		Joins("LEFT JOIN post_tags ON post_tags.tag_id = tags.id").
		Group("tags.id").
		Order("post_count DESC").
		Limit(limit).
		Find(&tags).Error
	return tags, err
}
