package service

import (
	"errors"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
	"notex/pkg/auth"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	postRepo    *repository.PostRepository
	commentRepo *repository.CommentRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo:    repository.NewUserRepository(),
		postRepo:    repository.NewPostRepository(),
		commentRepo: repository.NewCommentRepository(),
	}
}

// Register 用户注册
func (s *AuthService) Register(req *dto.RegisterRequest) error {
	// 检查用户名是否已存在
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return errors.New("email already exists")
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     "user", // 默认角色
		Status:   "active",
	}

	if err := user.SetPassword(req.Password); err != nil {
		return err
	}

	return s.userRepo.Create(user)
}

// Login 用户登录
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !user.IsActive() {
		return nil, errors.New("account is not active")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid username or password")
	}

	// 生成访问令牌
	claims := &dto.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	token, err := auth.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24小时
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	// 解析刷新令牌
	userID, err := auth.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive() {
		return nil, errors.New("account is not active")
	}

	// 生成新的访问令牌
	claims := &dto.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	token, err := auth.GenerateToken(claims)
	if err != nil {
		return nil, err
	}

	// 生成新的刷新令牌
	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        token,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24小时
		RefreshToken: refreshToken,
	}, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if !user.CheckPassword(req.OldPassword) {
		return errors.New("invalid old password")
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

// GetUserProfile 获取用户信息
func (s *AuthService) GetUserProfile(userID uint) (*dto.UserProfile, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户的文章数量
	postCount, err := s.postRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户的评论数量
	commentCount, err := s.commentRepo.CountByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 获取用户的文章总浏览量
	viewCount, err := s.postRepo.GetTotalViewsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserProfile{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		Status:       user.Status,
		Bio:          user.Bio,
		Avatar:       user.Avatar,
		PostCount:    postCount,
		CommentCount: commentCount,
		ViewCount:    viewCount,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}, nil
}

// IsUsernameExistsExcept 检查用户名是否已被使用（排除指定用户）
func (s *AuthService) IsUsernameExistsExcept(username string, excludeUserID uint) (bool, error) {
	return s.userRepo.IsUsernameExistsExcept(username, excludeUserID)
}

// UpdateUserProfile 更新用户信息
func (s *AuthService) UpdateUserProfile(user *model.User) error {
	return s.userRepo.Update(user)
}

// GetUserByID 通过ID获取用户
func (s *AuthService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}
