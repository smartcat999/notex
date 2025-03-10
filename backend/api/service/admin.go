package service

import (
	"errors"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
)

type AdminService struct {
	userRepo     *repository.UserRepository
	auditLogRepo *repository.AuditLogRepository
}

func NewAdminService() *AdminService {
	return &AdminService{
		userRepo:     repository.NewUserRepository(),
		auditLogRepo: repository.NewAuditLogRepository(),
	}
}

// ListUsers 获取用户列表
func (s *AdminService) ListUsers(req *dto.UserListRequest) (*dto.UserListResponse, error) {
	users, total, err := s.userRepo.FindAll(req.Page, req.PageSize, req.Role, req.Status, req.Search)
	if err != nil {
		return nil, err
	}

	items := make([]dto.UserResponse, len(users))
	for i, user := range users {
		items[i] = dto.UserResponse{
			ID:            user.ID,
			Email:         user.Email,
			Username:      user.Username,
			Role:          user.Role,
			Status:        user.Status,
			EmailVerified: user.EmailVerified,
		}
	}

	return &dto.UserListResponse{
		Total: total,
		Items: items,
	}, nil
}

// GetUser 获取用户详情
func (s *AdminService) GetUser(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
	}, nil
}

// UpdateUser 更新用户信息
func (s *AdminService) UpdateUser(id uint, req *dto.UserUpdateRequest) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 不允许修改超级管理员
	if user.Role == model.RoleAdmin {
		return errors.New("cannot modify admin user")
	}

	// 更新角色
	if req.Role != nil {
		if *req.Role != model.RoleUser && *req.Role != model.RoleAdmin {
			return errors.New("invalid role")
		}
		user.Role = *req.Role
	}

	// 更新状态
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "inactive" && *req.Status != "banned" {
			return errors.New("invalid status")
		}
		user.Status = *req.Status
	}

	return s.userRepo.Update(user)
}

// DeleteUser 删除用户
func (s *AdminService) DeleteUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 不允许删除超级管理员
	if user.Role == model.RoleAdmin {
		return errors.New("cannot delete admin user")
	}

	return s.userRepo.Delete(id)
}

// ListAuditLogs 获取审计日志列表
func (s *AdminService) ListAuditLogs(page, pageSize int, userID uint, action, resource string) ([]*model.AuditLog, int64, error) {
	return s.auditLogRepo.FindAll(page, pageSize, userID, action, resource)
}

// CleanupAuditLogs 清理旧的审计日志
func (s *AdminService) CleanupAuditLogs() error {
	return s.auditLogRepo.DeleteOldLogs()
}
