package service

import (
	"errors"
	"fmt"
	"math/rand"
	"notex/api/dto"
	"notex/api/repository"
	"notex/model"
	"notex/pkg/email"
	"time"
)

type VerificationService struct {
	repo     *repository.VerificationRepository
	userRepo *repository.UserRepository
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		repo:     repository.NewVerificationRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// generateVerificationCode 生成6位数字验证码
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// SendVerification 发送验证码
func (s *VerificationService) SendVerification(req *dto.SendVerificationRequest, vType model.VerificationType) error {
	// 检查是否有未过期的验证码
	if existing, _ := s.repo.FindByEmailAndType(req.Email, vType); existing != nil {
		if time.Since(existing.CreatedAt) < time.Minute {
			return errors.New("please wait for 1 minute before requesting a new code")
		}
	}

	// 生成验证码
	code := generateVerificationCode()

	// 创建验证记录
	verification := &model.Verification{
		Email:     req.Email,
		Code:      code,
		Type:      vType,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// 保存验证记录
	if err := s.repo.Create(verification); err != nil {
		return err
	}

	// 获取用户的语言偏好（这里可以从用户设置或请求头中获取）
	locale := "zh-CN" // 默认使用中文

	// 发送验证邮件
	switch vType {
	case model.VerificationTypeEmail:
		return email.SendVerificationEmail(req.Email, code, locale)
	case model.VerificationTypePasswordReset:
		return email.SendPasswordResetEmail(req.Email, code, locale)
	default:
		return errors.New("invalid verification type")
	}
}

// VerifyEmail 验证邮箱
func (s *VerificationService) VerifyEmail(req *dto.VerifyEmailRequest) error {
	// 查找验证码记录
	verification, err := s.repo.FindByEmailAndCode(
		req.Email,
		req.Code,
		model.VerificationTypeEmail,
	)
	if err != nil {
		return errors.New("invalid verification code")
	}

	if !verification.IsValid() {
		return errors.New("verification code has expired")
	}

	// 查找用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// 标记验证码为已使用
	if err := s.repo.MarkAsUsed(verification.ID); err != nil {
		return err
	}

	// 标记邮箱为已验证
	return s.userRepo.MarkEmailAsVerified(user.ID)
}

// UpdateEmail 更新邮箱
func (s *VerificationService) UpdateEmail(userID uint, req *dto.UpdateEmailRequest) error {
	// 检查新邮箱是否已被使用
	if _, err := s.userRepo.FindByEmail(req.NewEmail); err == nil {
		return errors.New("email already in use")
	}

	// 更新邮箱
	if err := s.userRepo.UpdateEmail(userID, req.NewEmail); err != nil {
		return err
	}

	// 发送验证码到新邮箱
	return s.SendVerification(&dto.SendVerificationRequest{
		Email: req.NewEmail,
	}, model.VerificationTypeEmail)
}

// ResetPassword 重置密码
func (s *VerificationService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// 查找验证码记录
	verification, err := s.repo.FindByEmailAndCode(
		req.Email,
		req.Code,
		model.VerificationTypePasswordReset,
	)
	if err != nil {
		return errors.New("invalid verification code")
	}

	if !verification.IsValid() {
		return errors.New("verification code has expired")
	}

	// 查找用户
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// 更新密码
	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	// 标记验证码为已使用
	if err := s.repo.MarkAsUsed(verification.ID); err != nil {
		return err
	}

	return s.userRepo.Update(user)
}

// CleanupExpiredCodes 清理过期的验证码
func (s *VerificationService) CleanupExpiredCodes() error {
	return s.repo.DeleteExpired()
}
