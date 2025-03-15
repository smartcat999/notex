package service

import (
	"fmt"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type CommentService struct {
	repo            *repository.CommentRepository
	notificationSvc *NotificationService
	postRepo        *repository.PostRepository
}

func NewCommentService() *CommentService {
	return &CommentService{
		repo:            repository.NewCommentRepository(),
		notificationSvc: NewNotificationService(),
		postRepo:        repository.NewPostRepository(),
	}
}

// CreateComment creates a new comment
func (s *CommentService) CreateComment(userID uint, postID uint, req *dto.CommentRequest) (*dto.CommentResponse, error) {
	// 获取文章信息
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, err
	}

	// Create comment model
	comment := &model.Comment{
		Content:   req.Content,
		UserID:    userID,
		PostID:    postID,
		ParentID:  req.ParentID,
		ReplyToID: req.ReplyToID,
		Status:    "active",
	}

	// Save comment
	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	// Get the created comment with relationships
	createdComment, err := s.repo.FindByID(comment.ID)
	if err != nil {
		return nil, err
	}

	// 处理通知
	if req.ParentID != nil {
		// 如果是回复评论，创建回复通知
		parentComment, err := s.repo.FindByID(*req.ParentID)
		if err != nil {
			// 记录错误但不影响评论创建
			// TODO: 添加日志记录
		} else {
			// 如果存在具体回复的评论，优先发送通知给被回复的用户
			if req.ReplyToID != nil {
				replyToComment, err := s.repo.FindByID(*req.ReplyToID)
				if err != nil {
					// 记录错误但不影响评论创建
					// TODO: 添加日志记录
				} else {
					// 发送通知给被回复的用户
					if err := s.notificationSvc.CreateReplyNotification(userID, *req.ReplyToID, comment.ID, req.Content); err != nil {
						// 记录错误但不影响评论创建
						// TODO: 添加日志记录
					}

					// 如果被回复的用户不是父评论作者，且父评论存在，才发送通知给父评论作者
					if replyToComment.UserID != parentComment.UserID {
						if err := s.notificationSvc.CreateReplyNotification(userID, *req.ParentID, comment.ID, req.Content); err != nil {
							// 记录错误但不影响评论创建
							// TODO: 添加日志记录
						}
					}
				}
			} else {
				// 如果没有具体回复的评论，直接发送通知给父评论作者
				if err := s.notificationSvc.CreateReplyNotification(userID, *req.ParentID, comment.ID, req.Content); err != nil {
					// 记录错误但不影响评论创建
					// TODO: 添加日志记录
				}
			}
		}
	} else {
		// 如果是对文章的直接评论，创建评论通知
		if err := s.notificationSvc.CreateCommentNotification(userID, postID, comment.ID, post.Title, req.Content); err != nil {
			// 记录错误但不影响评论创建
			// TODO: 添加日志记录
		}
	}

	// Convert to response
	response, err := s.convertToResponse(createdComment)
	if err != nil {
		return nil, err
	}
	return response, nil
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
		// 只处理顶级评论
		if comment.ParentID == nil {
			// 获取子评论数量
			replyCount, err := s.repo.GetReplyCount(comment.ID)
			if err != nil {
				// 记录错误但不影响评论列表的返回
				// TODO: 添加日志记录
				replyCount = 0
			}

			response, err := s.convertToResponse(&comment)
			if err != nil {
				return nil, 0, err
			}
			response.ReplyCount = replyCount
			responses = append(responses, *response)
		}
	}

	return responses, total, nil
}

// GetCommentReplies 获取评论的回复列表
func (s *CommentService) GetCommentReplies(commentID uint) ([]*dto.CommentBrief, error) {
	replies, err := s.repo.ListReplies(commentID)
	if err != nil {
		return nil, err
	}

	// 转换回复为简要信息
	children := make([]*dto.CommentBrief, len(replies))
	for i, reply := range replies {
		user, err := s.convertToUserInfo(reply.User)
		if err != nil {
			return nil, err
		}

		brief := &dto.CommentBrief{
			ID:        reply.ID,
			Content:   reply.Content,
			UserID:    reply.UserID,
			CreatedAt: reply.CreatedAt,
			User:      user,
		}

		// 如果有被回复的评论，添加被回复评论的信息
		if reply.ReplyTo != nil {
			replyToUser, err := s.convertToUserInfo(reply.ReplyTo.User)
			if err != nil {
				return nil, err
			}
			brief.ReplyTo = &struct {
				ID   uint          `json:"id"`
				User *dto.UserInfo `json:"user"`
			}{
				ID:   reply.ReplyTo.ID,
				User: replyToUser,
			}
		}

		children[i] = brief
	}

	return children, nil
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
		user, err := s.convertToUserInfo(comment.User)
		if err != nil {
			return nil, 0, err
		}
		responses[i] = dto.CommentResponse{
			ID:        comment.ID,
			Content:   comment.Content,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			User:      user,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
		}
		if comment.Post != nil {
			responses[i].PostTitle = comment.Post.Title
		}
	}

	return responses, total, nil
}

// convertToResponse converts a comment model to a comment response
func (s *CommentService) convertToResponse(comment *model.Comment) (*dto.CommentResponse, error) {
	if comment == nil {
		return nil, nil
	}

	// Convert user information
	user, err := s.convertToUserInfo(comment.User)
	if err != nil {
		return nil, err
	}

	response := &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		User:      user,
	}

	// Add parent comment if exists
	if comment.Parent != nil {
		parentUser, err := s.convertToUserInfo(comment.Parent.User)
		if err != nil {
			return nil, err
		}
		response.Parent = &dto.CommentBrief{
			ID:        comment.Parent.ID,
			Content:   comment.Parent.Content,
			CreatedAt: comment.Parent.CreatedAt,
			User:      parentUser,
		}
	}

	// Add reply_to comment if exists
	if comment.ReplyTo != nil {
		replyToUser, err := s.convertToUserInfo(comment.ReplyTo.User)
		if err != nil {
			return nil, err
		}
		response.ReplyTo = &dto.CommentBrief{
			ID:        comment.ReplyTo.ID,
			Content:   comment.ReplyTo.Content,
			CreatedAt: comment.ReplyTo.CreatedAt,
			User:      replyToUser,
		}
	}

	// Add children comments if any
	children := make([]*dto.CommentBrief, 0, len(comment.Children))
	for _, child := range comment.Children {
		childUser, err := s.convertToUserInfo(child.User)
		if err != nil {
			return nil, err
		}
		children = append(children, &dto.CommentBrief{
			ID:        child.ID,
			Content:   child.Content,
			CreatedAt: child.CreatedAt,
			User:      childUser,
		})
	}
	response.Children = children

	return response, nil
}

// convertToUserInfo converts a user model to a user info DTO
func (s *CommentService) convertToUserInfo(user *model.User) (*dto.UserInfo, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}
	return &dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
	}, nil
}
