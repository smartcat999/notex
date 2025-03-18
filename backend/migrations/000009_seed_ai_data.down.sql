-- 删除AI模型数据
DELETE FROM ai_models WHERE provider IN ('openai', 'anthropic', 'google', 'deepseek', 'custom');

-- 删除AI提供商数据
DELETE FROM ai_providers WHERE provider_id IN ('openai', 'anthropic', 'google', 'deepseek', 'custom'); 