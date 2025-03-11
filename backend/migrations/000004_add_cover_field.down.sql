-- 从 posts 表中删除 cover 字段
ALTER TABLE posts DROP COLUMN IF EXISTS cover;

-- 从 drafts 表中删除 cover 字段
ALTER TABLE drafts DROP COLUMN IF EXISTS cover; 