package service

import (
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type TagService struct {
	repo *repository.TagRepository
}

func NewTagService() *TagService {
	return &TagService{
		repo: repository.NewTagRepository(),
	}
}

// CreateTag 创建标签
func (s *TagService) CreateTag(req *dto.CreateTagRequest) (*dto.TagResponse, error) {
	tag := &model.Tag{
		Name: req.Name,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return s.convertToResponse(tag)
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(id uint, req *dto.UpdateTagRequest) (*dto.TagResponse, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	tag.Name = req.Name

	if err := s.repo.Update(tag); err != nil {
		return nil, err
	}

	return s.convertToResponse(tag)
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(id uint) error {
	return s.repo.Delete(id)
}

// GetTag 获取标签详情
func (s *TagService) GetTag(id uint) (*dto.TagResponse, error) {
	tag, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(tag)
}

// ListTags 获取标签列表
func (s *TagService) ListTags(query *dto.TagListQuery) ([]dto.TagResponse, int64, error) {
	tags, total, err := s.repo.List(query.Page, query.PageSize, query.Search)
	if err != nil {
		return make([]dto.TagResponse, 0), 0, err
	}

	responses := make([]dto.TagResponse, 0)
	for _, tag := range tags {
		response, err := s.convertToResponse(&tag)
		if err != nil {
			return responses, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

// GetTopTags 获取热门标签（按文章数量排序）
func (s *TagService) GetTopTags(limit int) ([]dto.TagResponse, error) {
	tags, err := s.repo.ListByPostCount(limit)
	if err != nil {
		return make([]dto.TagResponse, 0), err
	}

	responses := make([]dto.TagResponse, 0)
	for _, tag := range tags {
		response, err := s.convertToResponse(&tag)
		if err != nil {
			return responses, err
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// convertToResponse 将标签模型转换为响应DTO
func (s *TagService) convertToResponse(tag *model.Tag) (*dto.TagResponse, error) {
	postCount, err := s.repo.GetPostCount(tag.ID)
	if err != nil {
		return nil, err
	}

	return &dto.TagResponse{
		ID:        tag.ID,
		Name:      tag.Name,
		PostCount: postCount,
		CreatedAt: tag.CreatedAt,
		UpdatedAt: tag.UpdatedAt,
	}, nil
}
