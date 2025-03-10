package email

import (
	"fmt"
	"net/smtp"
	"notex/config"
	"notex/pkg/i18n"
	"notex/pkg/template"
)

type EmailSender struct {
	config config.EmailConfig
}

var defaultSender *EmailSender

// Initialize 初始化邮件发送器
func Initialize(cfg config.EmailConfig) error {
	defaultSender = NewEmailSender(cfg)

	// 初始化模板系统
	if err := template.Initialize(cfg.TemplatesDir); err != nil {
		return fmt.Errorf("failed to initialize email templates: %v", err)
	}

	return nil
}

// NewEmailSender 创建新的邮件发送器
func NewEmailSender(cfg config.EmailConfig) *EmailSender {
	return &EmailSender{
		config: cfg,
	}
}

// SendEmail 发送邮件
func (s *EmailSender) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	msg := []byte(fmt.Sprintf("From: %s <%s>\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", s.config.FromName, s.config.FromEmail, to, subject, body))

	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}

// SendTemplateEmail 使用模板发送邮件
func (s *EmailSender) SendTemplateEmail(to, templateName string, data interface{}, locale string) error {
	// 渲染模板
	body, err := template.RenderTemplate(templateName, data, locale)
	if err != nil {
		return fmt.Errorf("failed to render template: %v", err)
	}

	// 获取模板标题
	subject := i18n.T(locale, templateName+".subject")

	// 发送邮件
	return s.SendEmail(to, subject, body)
}

// SendVerificationEmail 发送验证邮件
func (s *EmailSender) SendVerificationEmail(to, code string, locale string) error {
	data := struct {
		Code      string
		ExpiresIn int
	}{
		Code:      code,
		ExpiresIn: 15, // 15分钟过期
	}

	return s.SendTemplateEmail(to, "verification", data, locale)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *EmailSender) SendPasswordResetEmail(to, code string, locale string) error {
	data := struct {
		Code      string
		ExpiresIn int
	}{
		Code:      code,
		ExpiresIn: 15, // 15分钟过期
	}

	return s.SendTemplateEmail(to, "password_reset", data, locale)
}

// SendEmailChangeNotification 发送邮箱变更通知
func (s *EmailSender) SendEmailChangeNotification(to, newEmail, code string, locale string) error {
	data := struct {
		Code      string
		NewEmail  string
		ExpiresIn int
	}{
		Code:      code,
		NewEmail:  newEmail,
		ExpiresIn: 15, // 15分钟过期
	}

	return s.SendTemplateEmail(to, "email_change", data, locale)
}

// 以下是包级别的便捷函数，使用默认发送器

// SendEmail 使用默认发送器发送邮件
func SendEmail(to, subject, body string) error {
	if defaultSender == nil {
		return fmt.Errorf("email sender not initialized")
	}
	return defaultSender.SendEmail(to, subject, body)
}

// SendVerificationEmail 使用默认发送器发送验证邮件
func SendVerificationEmail(to, code string, locale string) error {
	if defaultSender == nil {
		return fmt.Errorf("email sender not initialized")
	}
	return defaultSender.SendVerificationEmail(to, code, locale)
}

// SendPasswordResetEmail 使用默认发送器发送密码重置邮件
func SendPasswordResetEmail(to, code string, locale string) error {
	if defaultSender == nil {
		return fmt.Errorf("email sender not initialized")
	}
	return defaultSender.SendPasswordResetEmail(to, code, locale)
}

// SendEmailChangeNotification 使用默认发送器发送邮箱变更通知
func SendEmailChangeNotification(to, newEmail, code string, locale string) error {
	if defaultSender == nil {
		return fmt.Errorf("email sender not initialized")
	}
	return defaultSender.SendEmailChangeNotification(to, newEmail, code, locale)
}

// PreviewTemplate 预览邮件模板
func PreviewTemplate(templateName, locale string) (string, error) {
	if defaultSender == nil {
		return "", fmt.Errorf("email sender not initialized")
	}
	return template.PreviewTemplate(templateName, locale)
}

// GetAvailableTemplates 获取所有可用的模板
func GetAvailableTemplates() []string {
	return template.GetAvailableTemplates()
}

// GetSupportedLocales 获取支持的语言
func GetSupportedLocales() []string {
	return i18n.GetLocales()
}
