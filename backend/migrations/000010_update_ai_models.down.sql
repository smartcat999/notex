-- 删除添加的图像生成模型
DELETE FROM ai_models WHERE type = 'image';

-- 删除新添加的提供商
DELETE FROM ai_providers WHERE provider_id = 'stabilityai';

-- 移除type字段
ALTER TABLE ai_models DROP COLUMN IF EXISTS type; 