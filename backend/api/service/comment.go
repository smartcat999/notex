package service

import (
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type CommentService struct {
	repo *repository.CommentRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		repo: repository.NewCommentRepository(),
	}
}

// CreateComment 创建评论
func (s *CommentService) CreateComment(userID, postID uint, req *dto.CreateCommentRequest) (*dto.CommentResponse, error) {
	comment := &model.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  postID,
		Status:  "active",
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return s.convertToResponse(comment)
}

// DeleteComment 删除评论
func (s *CommentService) DeleteComment(id uint) error {
	return s.repo.Delete(id)
}

// GetComment 获取评论详情
func (s *CommentService) GetComment(id uint) (*dto.CommentResponse, error) {
	comment, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.convertToResponse(comment)
}

// ListComments 获取文章评论列表
func (s *CommentService) ListComments(postID uint, query *dto.CommentListQuery) ([]dto.CommentResponse, int64, error) {
	comments, total, err := s.repo.ListByPostID(postID, query.Page, query.PageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.CommentResponse, 0)
	for _, comment := range comments {
		response, err := s.convertToResponse(&comment)
		if err != nil {
			return nil, 0, err
		}
		responses = append(responses, *response)
	}

	return responses, total, nil
}

// ListUserComments 获取用户的评论列表
func (s *CommentService) ListUserComments(userID uint, page, pageSize int) ([]dto.CommentResponse, int64, error) {
	comments, total, err := s.repo.ListByUser(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	responses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = dto.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
		if comment.Post != nil {
			responses[i].PostTitle = comment.Post.Title
		}
	}

	return responses, total, nil
}

// convertToResponse 将评论模型转换为响应DTO
func (s *CommentService) convertToResponse(comment *model.Comment) (*dto.CommentResponse, error) {
	if comment == nil {
		return nil, nil
	}

	response := &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}

	if comment.User != nil {
		response.User = &dto.UserInfo{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Avatar:   comment.User.Avatar,
		}
	}

	return response, nil
}
