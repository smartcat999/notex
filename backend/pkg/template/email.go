package template

import (
	"bytes"
	"fmt"
	"html/template"
	"notex/pkg/i18n"
	"path/filepath"
	"time"
)

// EmailTemplateManager 邮件模板管理器
type EmailTemplateManager struct {
	baseTemplate *template.Template
	templates    map[string]*template.Template
}

// 全局模板管理器实例
var defaultManager *EmailTemplateManager

// Initialize 初始化邮件模板管理器
func Initialize(templatesDir string) error {
	manager := &EmailTemplateManager{
		templates: make(map[string]*template.Template),
	}

	// 加载基础模板
	baseTemplatePath := filepath.Join(templatesDir, "base.html")
	var err error
	manager.baseTemplate, err = template.ParseFiles(baseTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to load base template: %v", err)
	}

	// 加载其他模板
	templates := []string{
		"verification.html",
		"password_reset.html",
		"email_change.html",
	}

	for _, tmpl := range templates {
		path := filepath.Join(templatesDir, tmpl)
		t, err := template.Must(manager.baseTemplate.Clone()).ParseFiles(path)
		if err != nil {
			return fmt.Errorf("failed to load template %s: %v", tmpl, err)
		}
		manager.templates[tmpl[:len(tmpl)-5]] = t // 去掉 .html 后缀
	}

	defaultManager = manager
	return nil
}

// BaseTemplateData 基础模板数据
type BaseTemplateData struct {
	Title   string
	Content template.HTML
	Year    int
	Locale  string
	T       func(key string, args ...interface{}) string
}

// RenderTemplate 渲染指定模板
func RenderTemplate(templateName string, data interface{}, locale string) (string, error) {
	if defaultManager == nil {
		return "", fmt.Errorf("template manager not initialized")
	}

	// 先渲染内容模板
	tmpl, ok := defaultManager.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var contentBuf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&contentBuf, templateName, data); err != nil {
		return "", fmt.Errorf("failed to render content template: %v", err)
	}

	// 准备基础模板数据
	baseData := BaseTemplateData{
		Title:   i18n.T(locale, templateName+".title"),
		Content: template.HTML(contentBuf.String()),
		Year:    time.Now().Year(),
		Locale:  locale,
		T: func(key string, args ...interface{}) string {
			return i18n.T(locale, key, args...)
		},
	}

	// 渲染基础模板
	var buf bytes.Buffer
	if err := defaultManager.baseTemplate.Execute(&buf, baseData); err != nil {
		return "", fmt.Errorf("failed to render base template: %v", err)
	}

	return buf.String(), nil
}

// PreviewTemplate 预览指定模板
func PreviewTemplate(templateName string, locale string) (string, error) {
	// 准备预览数据
	previewData := map[string]interface{}{
		"Code":      "123456",
		"ExpiresIn": 15,
		"NewEmail":  "new@example.com",
	}

	return RenderTemplate(templateName, previewData, locale)
}

// GetTemplateTitle 获取模板标题（导出函数）
func GetTemplateTitle(templateName string, locale string) string {
	return i18n.T(locale, templateName+".title")
}

// GetAvailableTemplates 获取所有可用的模板
func GetAvailableTemplates() []string {
	if defaultManager == nil {
		return nil
	}

	templates := make([]string, 0, len(defaultManager.templates))
	for name := range defaultManager.templates {
		templates = append(templates, name)
	}
	return templates
}
