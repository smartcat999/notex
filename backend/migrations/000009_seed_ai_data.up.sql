-- 插入AI提供商数据
INSERT INTO ai_providers (provider_id, name, description, has_endpoint, is_enabled)
VALUES 
('openai', 'OpenAI', 'OpenAI提供的AI模型服务', false, true),
('anthropic', 'Anthropic', 'Anthropic提供的Claude系列AI模型', false, true),
('google', 'Google AI', 'Google提供的Gemini系列AI模型', false, true),
('deepseek', 'DeepSeek', 'DeepSeek提供的AI模型服务', false, true),
('custom', '自定义模型', '配置自定义AI模型API端点', true, true)
ON CONFLICT (provider_id) DO NOTHING;

-- 插入AI模型数据
INSERT INTO ai_models (provider, model_id, name, description, is_paid, is_enabled)
VALUES 
-- OpenAI模型
('openai', 'gpt-3.5-turbo', 'GPT-3.5 Turbo', '强大的语言模型，适合大多数任务，响应速度快。', true, true),
('openai', 'gpt-4', 'GPT-4', '最先进的语言模型，具有更强的推理能力和更广泛的知识。', true, true),

-- Anthropic模型
('anthropic', 'claude-3-opus', 'Claude 3 Opus', 'Anthropic的顶级模型，具有强大的推理和创作能力。', true, true),
('anthropic', 'claude-3-sonnet', 'Claude 3 Sonnet', '平衡性能和速度的模型，适合大多数任务。', true, true),

-- Google模型
('google', 'gemini-pro', 'Gemini Pro', 'Google的多模态AI模型，具有强大的理解和生成能力。', true, true),

-- DeepSeek模型
('deepseek', 'deepseek-chat', 'DeepSeek Chat', 'DeepSeek的通用对话模型，擅长自然语言理解和生成。', true, true),
('deepseek', 'deepseek-coder', 'DeepSeek Coder', 'DeepSeek的代码生成模型，专注于编程和开发任务。', true, true),

-- 自定义模型
('custom', 'custom-model', '自定义模型', '配置您自己的AI模型API端点。', false, true)
ON CONFLICT DO NOTHING; 