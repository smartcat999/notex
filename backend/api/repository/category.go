package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		db: database.GetDB(),
	}
}

// Create 创建分类
func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

// FindByID 根据ID查找分类
func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// List 获取分类列表
func (r *CategoryRepository) List(page, pageSize int, search string) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.Model(&model.Category{})

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&categories).Error

	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetPostCount 获取分类下的文章数量
func (r *CategoryRepository) GetPostCount(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Post{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

// ListByPostCount 获取热门分类（按文章数量排序）
func (r *CategoryRepository) ListByPostCount(limit int) ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Model(&model.Category{}).
		Select("categories.*, COUNT(posts.id) as post_count").
		Joins("LEFT JOIN posts ON posts.category_id = categories.id").
		Group("categories.id").
		Order("post_count DESC").
		Limit(limit).
		Find(&categories).Error
	return categories, err
}
