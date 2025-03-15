-- 删除 reply_to_id 列
ALTER TABLE comments
DROP CONSTRAINT IF EXISTS fk_comments_reply_to,
DROP COLUMN IF EXISTS reply_to_id; 