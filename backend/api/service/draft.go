package service

import (
	"errors"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrDraftNotFound = errors.New("draft not found")
)

type DraftService struct {
	draftRepo *repository.DraftRepository
}

func NewDraftService(draftRepo *repository.DraftRepository) *DraftService {
	return &DraftService{
		draftRepo: draftRepo,
	}
}

// ListDrafts 获取草稿列表
func (s *DraftService) ListDrafts(userID uint, page, pageSize int, search string) ([]dto.DraftResponse, int64, error) {
	conditions := map[string]interface{}{
		"user_id": userID,
	}
	if search != "" {
		conditions["search"] = search
	}

	drafts, total, err := s.draftRepo.List(page, pageSize, conditions)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]dto.DraftResponse, 0)
	for _, draft := range drafts {
		responses = append(responses, convertDraftToResponse(&draft))
	}

	return responses, total, nil
}

// GetDraft 获取草稿详情
func (s *DraftService) GetDraft(id, userID uint) (*model.Draft, error) {
	draft, err := s.draftRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrDraftNotFound) {
			return nil, ErrDraftNotFound
		}
		return nil, err
	}

	// 验证所有权
	if draft.UserID != userID {
		return nil, ErrUnauthorized
	}

	return draft, nil
}

// CreateDraft 创建草稿
func (s *DraftService) CreateDraft(userID uint, req dto.CreateDraftRequest) (*dto.DraftResponse, error) {
	// 开启事务
	tx := s.draftRepo.DB.Begin()

	draft := &model.Draft{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CategoryID: req.CategoryID,
		UserID:     userID,
	}

	// 创建草稿
	if err := tx.Create(draft).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果有标签，添加标签关联
	if len(req.TagIDs) > 0 {
		var tags []model.Tag
		if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(draft).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// 重新加载草稿以获取完整的关联数据
	draft, err := s.draftRepo.FindByID(draft.ID)
	if err != nil {
		return nil, err
	}

	response := convertDraftToResponse(draft)
	return &response, nil
}

// UpdateDraft 更新草稿
func (s *DraftService) UpdateDraft(id, userID uint, req dto.UpdateDraftRequest) (*model.Draft, error) {
	draft, err := s.draftRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrDraftNotFound) {
			return nil, ErrDraftNotFound
		}
		return nil, err
	}

	// 验证所有权
	if draft.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 更新字段
	if req.Title != "" {
		draft.Title = req.Title
	}
	if req.Content != "" {
		draft.Content = req.Content
	}
	if req.Summary != "" {
		draft.Summary = req.Summary
	}
	if req.CategoryID != 0 {
		draft.CategoryID = req.CategoryID
	}

	// 更新基本信息
	if err := s.draftRepo.Update(draft); err != nil {
		return nil, err
	}

	// 更新标签关联
	if len(req.TagIDs) > 0 {
		if err := s.draftRepo.UpdateDraftTags(draft, req.TagIDs); err != nil {
			return nil, err
		}
	}

	// 重新获取完整的草稿信息（包括关联的标签）
	return s.draftRepo.FindByID(id)
}

// DeleteDraft 删除草稿
func (s *DraftService) DeleteDraft(id, userID uint) error {
	draft, err := s.draftRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 验证所有权
	if draft.UserID != userID {
		return ErrUnauthorized
	}

	return s.draftRepo.Delete(draft)
}

// PublishDraft 发布草稿
func (s *DraftService) PublishDraft(id, userID uint) (*model.Post, error) {
	draft, err := s.draftRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// 验证所有权
	if draft.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 发布草稿
	post, err := s.draftRepo.PublishDraft(draft)
	if err != nil {
		return nil, err
	}

	return post, nil
}

// 辅助函数：转换草稿为响应格式
func convertDraftToResponse(draft *model.Draft) dto.DraftResponse {
	response := dto.DraftResponse{
		ID:         draft.ID,
		Title:      draft.Title,
		Content:    draft.Content,
		Summary:    draft.Summary,
		CategoryID: draft.CategoryID,
		CreatedAt:  draft.CreatedAt,
		UpdatedAt:  draft.UpdatedAt,
	}

	// 添加分类信息
	if draft.Category != nil {
		response.Category = draft.Category.Name
	}

	// 添加标签信息
	for _, tag := range draft.Tags {
		response.Tags = append(response.Tags, dto.TagInfo{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return response
}
