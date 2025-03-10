package service

import (
	"errors"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
	"time"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService() *PostService {
	return &PostService{
		repo: repository.NewPostRepository(),
	}
}

// CreatePost 创建文章
func (s *PostService) CreatePost(req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	post := &model.Post{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Slug:       req.Slug,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		UserID:     req.UserID,
	}

	if req.Status == "published" {
		post.PublishedAt = time.Now()
	}

	if err := s.repo.Create(post); err != nil {
		return nil, err
	}

	return s.convertToResponse(post)
}

// UpdatePost 更新文章
func (s *PostService) UpdatePost(id uint, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Summary != "" {
		post.Summary = req.Summary
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.CategoryID != 0 {
		post.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		if req.Status == "published" && post.Status != "published" {
			post.PublishedAt = time.Now()
		}
		post.Status = req.Status
	}

	if err := s.repo.Update(post); err != nil {
		return nil, err
	}

	return s.convertToResponse(post)
}

// DeletePost 删除文章
func (s *PostService) DeletePost(id uint) error {
	return s.repo.Delete(id)
}

// GetPost 获取文章详情
func (s *PostService) GetPost(id uint) (*dto.PostResponse, error) {
	post, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(post)
}

// ListPosts 获取文章列表
func (s *PostService) ListPosts(query *dto.PostListQuery) ([]dto.PostResponse, int64, error) {
	conditions := make(map[string]interface{})

	if query.Status != "" {
		conditions["status"] = query.Status
	}
	if query.CategoryID > 0 {
		conditions["category_id"] = query.CategoryID
	}
	if query.TagID > 0 {
		conditions["tag_id"] = query.TagID
	}
	if query.Search != "" {
		conditions["search"] = query.Search
	}
	if query.Sort != "" {
		conditions["sort"] = query.Sort
	}
	if query.UserID > 0 {
		conditions["user_id"] = query.UserID
	}

	posts, total, err := s.repo.List(query.Page, query.PageSize, conditions)
	if err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	responses := make([]dto.PostResponse, 0)
	for _, post := range posts {
		response, err := s.convertToResponse(&post)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

// GetRecentPosts 获取最新文章列表
func (s *PostService) GetRecentPosts(limit int) ([]dto.PostResponse, error) {
	conditions := map[string]interface{}{
		"status": "published",
	}
	posts, _, err := s.repo.List(1, limit, conditions)
	if err != nil {
		return make([]dto.PostResponse, 0), err
	}

	responses := make([]dto.PostResponse, 0)
	for _, post := range posts {
		response, err := s.convertToResponse(&post)
		if err != nil {
			return responses, err
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// IncrementViews 增加文章浏览量
func (s *PostService) IncrementViews(id uint) error {
	return s.repo.IncrementViews(id)
}

// GetArchives 获取文章归档列表
func (s *PostService) GetArchives() ([]dto.ArchiveResponse, error) {
	archives, err := s.repo.GetArchives()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ArchiveResponse, 0)
	for _, archive := range archives {
		date := archive["date"].(string)
		count := int64(archive["count"].(int64))

		responses = append(responses, dto.ArchiveResponse{
			Date:  date,
			Count: count,
		})
	}

	return responses, nil
}

// GetPostsByArchive 获取指定归档日期的文章列表
func (s *PostService) GetPostsByArchive(yearMonth string) ([]dto.PostResponse, error) {
	posts, err := s.repo.GetPostsByArchive(yearMonth)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.PostResponse, 0)
	for _, post := range posts {
		response, err := s.convertToResponse(&post)
		if err != nil {
			return nil, err
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

// convertToResponse 将文章模型转换为响应DTO
func (s *PostService) convertToResponse(post *model.Post) (*dto.PostResponse, error) {
	if post == nil {
		return nil, errors.New("post is nil")
	}

	// 获取评论数
	commentCount, err := s.repo.GetCommentCount(post.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.PostResponse{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		Summary:      post.Summary,
		Slug:         post.Slug,
		CategoryID:   post.CategoryID,
		Status:       post.Status,
		Views:        post.Views,
		CommentCount: commentCount,
		PublishedAt:  post.PublishedAt,
		CreatedAt:    post.CreatedAt,
		UpdatedAt:    post.UpdatedAt,
	}

	if post.Category.ID != 0 {
		response.Category = post.Category.Name
	}

	tags := make([]dto.TagInfo, 0)
	for _, tag := range post.Tags {
		tags = append(tags, dto.TagInfo{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}
	response.Tags = tags

	return response, nil
}
