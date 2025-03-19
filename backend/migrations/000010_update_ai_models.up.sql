-- 向ai_models表添加type字段用于区分模型类型
ALTER TABLE ai_models ADD COLUMN IF NOT EXISTS type VARCHAR(50) DEFAULT 'text' NOT NULL;

-- 为现有模型设置类型为文本模型
UPDATE ai_models SET type = 'text' WHERE type IS NULL;

-- 插入图像生成模型
INSERT INTO ai_models (provider, model_id, name, description, is_paid, is_enabled, type)
VALUES 
-- OpenAI图像生成模型
('openai', 'dall-e-3', 'DALL-E 3', 'OpenAI的最新图像生成模型，可根据文本描述创建高度真实和创意的图像。', true, true, 'image'),
('openai', 'dall-e-2', 'DALL-E 2', 'OpenAI的图像生成模型，支持创建多种风格的图像。', true, true, 'image'),

-- 稳定扩散模型
('stabilityai', 'stable-diffusion-xl', 'Stable Diffusion XL', '高质量的开源图像生成模型，支持多种输入控制。', false, true, 'image'),

-- Midjourney API (通过自定义端点)
('custom', 'midjourney', 'Midjourney', '著名的图像生成服务，需要通过自定义API端点接入。', true, true, 'image')
ON CONFLICT DO NOTHING;

-- 添加Stability AI提供商
INSERT INTO ai_providers (provider_id, name, description, has_endpoint, is_enabled)
VALUES 
('stabilityai', 'Stability AI', 'Stability AI提供的图像生成模型服务', false, true)
ON CONFLICT (provider_id) DO NOTHING; 