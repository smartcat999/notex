-- 创建AI提供商表
CREATE TABLE IF NOT EXISTS ai_providers (
    id SERIAL PRIMARY KEY,
    provider_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    has_endpoint BOOLEAN DEFAULT FALSE,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ai_providers_provider_id_unique UNIQUE (provider_id)
);

-- 创建AI模型表
CREATE TABLE IF NOT EXISTS ai_models (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    model_id VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    is_paid BOOLEAN DEFAULT FALSE,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建用户AI设置表
CREATE TABLE IF NOT EXISTS ai_user_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    provider_id VARCHAR(50) NOT NULL,
    api_key VARCHAR(500),
    endpoint VARCHAR(500),
    enabled_models JSONB DEFAULT '{}'::JSONB,
    model_params JSONB DEFAULT '{}'::JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建用户默认AI设置表
CREATE TABLE IF NOT EXISTS ai_default_settings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    default_model VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ai_default_settings_user_id_unique UNIQUE (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_ai_user_settings_user_id ON ai_user_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_ai_models_provider ON ai_models(provider); 