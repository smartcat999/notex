package service

import (
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		repo: repository.NewCategoryRepository(),
	}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	category := &model.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return s.convertToResponse(category)
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(id uint, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return s.convertToResponse(category)
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(id uint) error {
	return s.repo.Delete(id)
}

// GetCategory 获取分类详情
func (s *CategoryService) GetCategory(id uint) (*dto.CategoryResponse, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(category)
}

// ListCategories 获取分类列表
func (s *CategoryService) ListCategories(query *dto.CategoryListQuery) ([]dto.CategoryResponse, int64, error) {
	categories, total, err := s.repo.List(query.Page, query.PageSize, query.Search)
	if err != nil {
		return make([]dto.CategoryResponse, 0), 0, err
	}

	responses := make([]dto.CategoryResponse, 0)
	for _, category := range categories {
		response, err := s.convertToResponse(&category)
		if err != nil {
			return responses, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

// GetTopCategories 获取热门分类（按文章数量排序）
func (s *CategoryService) GetTopCategories(limit int) ([]dto.CategoryResponse, error) {
	categories, err := s.repo.ListByPostCount(limit)
	if err != nil {
		return make([]dto.CategoryResponse, 0), err
	}

	responses := make([]dto.CategoryResponse, 0)
	for _, category := range categories {
		response, err := s.convertToResponse(&category)
		if err != nil {
			return responses, err
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// convertToResponse 将分类模型转换为响应DTO
func (s *CategoryService) convertToResponse(category *model.Category) (*dto.CategoryResponse, error) {
	postCount, err := s.repo.GetPostCount(category.ID)
	if err != nil {
		return nil, err
	}

	return &dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		PostCount:   postCount,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}, nil
}
