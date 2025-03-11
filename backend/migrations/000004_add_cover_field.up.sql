-- 为 posts 表添加 cover 字段
ALTER TABLE posts ADD COLUMN IF NOT EXISTS cover VARCHAR(255);

-- 为 drafts 表添加 cover 字段
ALTER TABLE drafts ADD COLUMN IF NOT EXISTS cover VARCHAR(255); 