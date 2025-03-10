package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Translator 翻译器
type Translator struct {
	locales map[string]map[string]interface{}
}

var defaultTranslator *Translator

// Initialize 初始化翻译器
func Initialize(localesDir string) error {
	t := &Translator{
		locales: make(map[string]map[string]interface{}),
	}

	// 遍历语言目录
	entries, err := os.ReadDir(localesDir)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		locale := entry.Name()
		localePath := filepath.Join(localesDir, locale)

		// 读取该语言下的所有YAML文件
		files, err := os.ReadDir(localePath)
		if err != nil {
			return fmt.Errorf("failed to read locale directory %s: %v", locale, err)
		}

		t.locales[locale] = make(map[string]interface{})

		for _, file := range files {
			if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
				// 读取YAML文件
				data, err := os.ReadFile(filepath.Join(localePath, file.Name()))
				if err != nil {
					return fmt.Errorf("failed to read file %s: %v", file.Name(), err)
				}

				// 解析YAML
				var content map[string]interface{}
				if err := yaml.Unmarshal(data, &content); err != nil {
					return fmt.Errorf("failed to parse file %s: %v", file.Name(), err)
				}

				// 将内容添加到翻译映射中
				t.locales[locale] = mergeMap(t.locales[locale], content)
			}
		}
	}

	defaultTranslator = t
	return nil
}

// T 获取翻译文本
func T(locale, key string, args ...interface{}) string {
	if defaultTranslator == nil {
		return key
	}
	return defaultTranslator.T(locale, key, args...)
}

// T 获取翻译文本
func (t *Translator) T(locale, key string, args ...interface{}) string {
	// 如果找不到语言，使用英文
	if _, ok := t.locales[locale]; !ok {
		locale = "en-US"
	}

	// 按照路径获取翻译
	parts := strings.Split(key, ".")
	current := t.locales[locale]

	for _, part := range parts {
		if v, ok := current[part]; ok {
			if m, ok := v.(map[string]interface{}); ok {
				current = m
			} else if s, ok := v.(string); ok {
				if len(args) > 0 {
					return fmt.Sprintf(s, args...)
				}
				return s
			} else {
				return key
			}
		} else {
			return key
		}
	}

	return key
}

// GetLocales 获取所有支持的语言
func GetLocales() []string {
	if defaultTranslator == nil {
		return nil
	}
	locales := make([]string, 0, len(defaultTranslator.locales))
	for locale := range defaultTranslator.locales {
		locales = append(locales, locale)
	}
	return locales
}

// mergeMap 合并两个map
func mergeMap(m1, m2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range m1 {
		result[k] = v
	}

	for k, v := range m2 {
		if mv, ok := v.(map[string]interface{}); ok {
			if mv2, ok := result[k].(map[string]interface{}); ok {
				result[k] = mergeMap(mv2, mv)
			} else {
				result[k] = v
			}
		} else {
			result[k] = v
		}
	}

	return result
}
