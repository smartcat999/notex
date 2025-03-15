# AI 工具箱

AI 工具箱是 Notex 应用的一个功能模块，提供多种 AI 模型的聊天和配置功能。

## 功能特点

- 支持多种 AI 提供商：OpenAI、Anthropic、Google AI 和自定义模型
- 聊天界面支持 Markdown 格式
- 可配置不同 AI 模型的参数
- 支持保存和加载 AI 模型配置

## 组件结构

- `AIToolbox.vue` - 主页面，包含聊天和设置标签页
- `AIChat.vue` - AI 聊天组件
- `AISettings.vue` - AI 设置组件
- `AIAvatar.vue` - AI 头像组件

## 使用方法

1. 在设置页面配置 AI 模型的 API 密钥
2. 在聊天页面选择要使用的 AI 模型
3. 开始与 AI 助手对话

## 开发说明

### 添加新的 AI 提供商

1. 在 `aiService.js` 中添加新的 API 调用函数
2. 在 `ai.js` store 中的 `availableModels` 数组中添加新模型
3. 在 `AISettings.vue` 中的 `aiProviders` 数组中添加新提供商配置

### 自定义 AI 头像

可以在 `/public/images/ai/` 目录下添加 AI 提供商的头像图片，命名格式为 `{provider-id}-avatar.png`。

## 注意事项

- API 密钥存储在浏览器的 localStorage 中，请确保在生产环境中实现更安全的存储方式
- 不同 AI 提供商的 API 格式可能会随时间变化，请及时更新 `aiService.js` 中的 API 调用代码 