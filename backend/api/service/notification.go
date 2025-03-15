package service

import (
	"fmt"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type NotificationService struct {
	repo        *repository.NotificationRepository
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		repo:        repository.NewNotificationRepository(),
		postRepo:    repository.NewPostRepository(),
		commentRepo: repository.NewCommentRepository(),
	}
}

// CreateCommentNotification 创建评论通知
func (s *NotificationService) CreateCommentNotification(actorID, postID, commentID uint, postTitle, commentContent string) error {
	// 获取文章作者ID
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return err
	}

	// 如果评论者不是文章作者，则创建通知
	if actorID != post.UserID {
		notification := &model.Notification{
			Type:      model.NotificationTypePostComment,
			UserID:    post.UserID,
			ActorID:   actorID,
			PostID:    &postID,
			CommentID: &commentID,
			Content:   fmt.Sprintf("在你的文章《%s》中发表了评论: %s", postTitle, commentContent),
		}
		return s.repo.Create(notification)
	}
	return nil
}

// CreateReplyNotification 创建回复通知
func (s *NotificationService) CreateReplyNotification(actorID, parentCommentID, commentID uint, commentContent string) error {
	// 获取父评论信息
	parentComment, err := s.commentRepo.FindByID(parentCommentID)
	if err != nil {
		return err
	}

	// 如果回复者不是父评论作者，则创建通知
	if actorID != parentComment.UserID {
		notification := &model.Notification{
			Type:      model.NotificationTypeCommentReply,
			UserID:    parentComment.UserID,
			ActorID:   actorID,
			PostID:    &parentComment.PostID,
			CommentID: &commentID,
			Content:   fmt.Sprintf("回复了你的评论: %s", commentContent),
		}
		return s.repo.Create(notification)
	}
	return nil
}

// ListNotifications 获取用户的通知列表
func (s *NotificationService) ListNotifications(userID uint, query *dto.NotificationListQuery) ([]dto.NotificationResponse, int64, error) {
	notifications, total, err := s.repo.ListByUser(userID, query.Page, query.PageSize, query.Unread, query.Type)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = s.convertToResponse(&notification)
	}

	return responses, total, nil
}

// MarkAsRead 将通知标记为已读
func (s *NotificationService) MarkAsRead(id uint) error {
	return s.repo.MarkAsRead(id)
}

// MarkAllAsRead 将用户的所有通知标记为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.repo.MarkAllAsRead(userID)
}

// GetUnreadCount 获取用户未读通知数量
func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

// convertToResponse 将通知模型转换为响应DTO
func (s *NotificationService) convertToResponse(notification *model.Notification) dto.NotificationResponse {
	response := dto.NotificationResponse{
		ID:        notification.ID,
		Type:      notification.Type,
		Content:   notification.Content,
		IsRead:    notification.Read,
		CreatedAt: notification.CreatedAt,
	}

	// 设置可选字段
	if notification.PostID != nil {
		response.PostID = *notification.PostID
	}
	if notification.CommentID != nil {
		response.CommentID = *notification.CommentID
	}

	// 设置发送者信息
	if notification.Actor != nil {
		response.SenderID = notification.Actor.ID
		response.SenderUsername = notification.Actor.Username
		response.SenderAvatar = notification.Actor.Avatar
	}

	// 设置文章标题
	if notification.Post != nil {
		response.PostTitle = notification.Post.Title
	}

	return response
}
