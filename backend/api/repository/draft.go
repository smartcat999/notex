package repository

import (
	"errors"
	"fmt"
	"notex/model"
	"notex/pkg/database"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrDraftNotFound = errors.New("draft not found")
)

type DraftRepository struct {
	DB *gorm.DB
}

func NewDraftRepository() *DraftRepository {
	return &DraftRepository{
		DB: database.GetDB(),
	}
}

// Create 创建草稿
func (r *DraftRepository) Create(draft *model.Draft) error {
	return r.DB.Create(draft).Error
}

// Update 更新草稿
func (r *DraftRepository) Update(draft *model.Draft) error {
	return r.DB.Save(draft).Error
}

// Delete 删除草稿
func (r *DraftRepository) Delete(draft *model.Draft) error {
	return r.DB.Delete(draft).Error
}

// FindByID 根据ID查找草稿
func (r *DraftRepository) FindByID(id uint) (*model.Draft, error) {
	var draft model.Draft
	err := r.DB.Preload("Category").Preload("Tags").First(&draft, id).Error
	if err != nil {
		return nil, err
	}
	return &draft, nil
}

// List 获取草稿列表
func (r *DraftRepository) List(page, pageSize int, conditions map[string]interface{}) ([]model.Draft, int64, error) {
	var drafts []model.Draft
	var total int64

	query := r.DB.Model(&model.Draft{})

	// 应用查询条件
	for key, value := range conditions {
		if value != nil {
			switch key {
			case "user_id":
				query = query.Where("user_id = ?", value)
			case "search":
				searchTerm := "%" + value.(string) + "%"
				query = query.Where("LOWER(title) LIKE LOWER(?) OR LOWER(content) LIKE LOWER(?) OR LOWER(summary) LIKE LOWER(?)",
					searchTerm, searchTerm, searchTerm)
			case "category_id":
				query = query.Where("category_id = ?", value)
			case "tag_id":
				query = query.Joins("JOIN draft_tags ON drafts.id = draft_tags.draft_id").
					Where("draft_tags.tag_id = ?", value)
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
		Find(&drafts).Error

	if err != nil {
		return nil, 0, err
	}

	return drafts, total, nil
}

// CountByUserID 获取用户的草稿数量
func (r *DraftRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&model.Draft{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// PublishDraft 发布草稿为文章
func (r *DraftRepository) PublishDraft(draft *model.Draft) (*model.Post, error) {
	// 开启事务
	tx := r.DB.Begin()

	// 生成唯一的 slug
	timestamp := time.Now().UnixNano() / 1e6 // 转换为毫秒
	slug := fmt.Sprintf("%s-%d", generateSlug(draft.Title), timestamp)

	// 创建文章
	post := &model.Post{
		Title:       draft.Title,
		Content:     draft.Content,
		Summary:     draft.Summary,
		CategoryID:  draft.CategoryID,
		UserID:      draft.UserID,
		Status:      "published",
		Slug:        slug,
		PublishedAt: time.Now(),
	}

	// 创建文章
	if err := tx.Create(post).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果有标签，添加标签关联
	if len(draft.Tags) > 0 {
		if err := tx.Model(post).Association("Tags").Replace(draft.Tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 删除草稿
	if err := tx.Delete(draft).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新加载文章以获取完整的关联数据
	if err := r.DB.Preload("Category").Preload("Tags").First(post, post.ID).Error; err != nil {
		return nil, err
	}

	return post, nil
}

// generateSlug 生成 URL 友好的 slug
func generateSlug(title string) string {
	// 将标题转换为小写
	slug := strings.ToLower(title)

	// 将中文字符转换为拼音（可选，如果需要的话）

	// 将非字母数字字符替换为连字符
	reg := regexp.MustCompile(`[^\w\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// 将空格替换为连字符
	reg = regexp.MustCompile(`[-\s]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// 移除首尾的连字符
	slug = strings.Trim(slug, "-")

	// 限制长度
	if len(slug) > 100 {
		slug = slug[:100]
	}

	return slug
}

// UpdateDraftTags 更新草稿的标签关联
func (r *DraftRepository) UpdateDraftTags(draft *model.Draft, tagIDs []uint) error {
	// 开启事务
	tx := r.DB.Begin()

	// 清除现有的标签关联
	if err := tx.Model(draft).Association("Tags").Clear(); err != nil {
		tx.Rollback()
		return err
	}

	// 如果有新的标签ID，添加新的关联
	if len(tagIDs) > 0 {
		var tags []model.Tag
		if err := tx.Where("id IN ?", tagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(draft).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	return tx.Commit().Error
}
