-- 添加default_image_model字段到ai_default_settings表
ALTER TABLE ai_default_settings ADD COLUMN default_image_model VARCHAR(100); 