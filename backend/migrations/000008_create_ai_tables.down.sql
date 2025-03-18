-- 删除索引
DROP INDEX IF EXISTS idx_ai_user_settings_user_id;
DROP INDEX IF EXISTS idx_ai_models_provider;

-- 删除表
DROP TABLE IF EXISTS ai_default_settings;
DROP TABLE IF EXISTS ai_user_settings;
DROP TABLE IF EXISTS ai_models;
DROP TABLE IF EXISTS ai_providers; 