-- 删除用户个人资料字段
ALTER TABLE users
DROP COLUMN bio,
DROP COLUMN avatar; 