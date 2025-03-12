package repository

import (
	"notex/model"
	"notex/pkg/database"

	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository() *PostRepository {
	return &PostRepository{
		DB: database.GetDB(),
	}
}

// Create 创建文章
func (r *PostRepository) Create(post *model.Post) error {
	return r.DB.Create(post).Error
}

// Update 更新文章
func (r *PostRepository) Update(post *model.Post) error {
	return r.DB.Save(post).Error
}

// Delete 删除文章
func (r *PostRepository) Delete(id uint) error {
	return r.DB.Delete(&model.Post{}, id).Error
}

// FindByID 根据ID查找文章
func (r *PostRepository) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.DB.Preload("Category").Preload("Tags").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// List 获取文章列表
func (r *PostRepository) List(page, pageSize int, conditions map[string]interface{}) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.DB.Model(&model.Post{})

	// 应用查询条件
	for key, value := range conditions {
		if value != nil {
			switch key {
			case "status":
				query = query.Where("status = ?", value)
			case "category_id":
				query = query.Where("category_id = ?", value)
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "search":
				searchTerm := "%" + value.(string) + "%"
				query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?) OR LOWER(summary) LIKE LOWER(?)",
					searchTerm, searchTerm, searchTerm)
			case "tag_id":
				query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
					Where("post_tags.tag_id = ?", value)
			case "sort":
				switch value.(string) {
				case "newest":
					query = query.Order("published_at DESC, created_at DESC")
				case "most_viewed":
					query = query.Order("views DESC, published_at DESC")
				case "most_commented":
					query = query.Joins("LEFT JOIN comments ON posts.id = comments.post_id").
						Group("posts.id").
						Order("COUNT(comments.id) DESC, published_at DESC")
				}
			}
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Preload("Category").
		Preload("Tags").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// CountByUserID 获取用户的文章数量
func (r *PostRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Post{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// GetTotalViewsByUserID 获取用户的文章总浏览量
func (r *PostRepository) GetTotalViewsByUserID(userID uint) (int64, error) {
	var totalViews int64
	err := r.DB.Model(&model.Post{}).Where("user_id = ?", userID).Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error
	return totalViews, err
}

// GetCommentCount 获取文章的评论数量
func (r *PostRepository) GetCommentCount(postID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

// IncrementViews 增加文章浏览量
func (r *PostRepository) IncrementViews(postID uint) error {
	return r.DB.Model(&model.Post{}).Where("id = ?", postID).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

// GetArchives 获取文章归档列表
func (r *PostRepository) GetArchives() ([]map[string]interface{}, error) {
	var archives []map[string]interface{}

	// 使用 PostgreSQL 的 to_char 函数按年月分组统计文章数量
	err := r.DB.Model(&model.Post{}).
		Select("to_char(published_at, 'YYYY-MM') as date, COUNT(*) as count").
		Where("status = ? AND published_at IS NOT NULL", "published").
		Group("to_char(published_at, 'YYYY-MM')").
		Order("date DESC").
		Find(&archives).Error

	if err != nil {
		return nil, err
	}

	return archives, nil
}

// GetPostsByArchive 获取指定归档日期的文章列表
func (r *PostRepository) GetPostsByArchive(yearMonth string) ([]model.Post, error) {
	var posts []model.Post

	err := r.DB.Model(&model.Post{}).
		Where("status = ? AND to_char(published_at, 'YYYY-MM') = ?", "published", yearMonth).
		Preload("Category").
		Preload("Tags").
		Order("published_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// ListByUserID 获取用户的文章列表
func (r *PostRepository) ListByUserID(userID uint, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.DB.Model(&model.Post{}).Where("user_id = ?", userID)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Preload("Category").
		Preload("Tags").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("published_at DESC, created_at DESC").
		Find(&posts).Error

	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
