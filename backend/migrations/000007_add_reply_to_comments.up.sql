-- 添加 reply_to_id 列到 comments 表
ALTER TABLE comments
ADD COLUMN reply_to_id bigint NULL,
ADD CONSTRAINT fk_comments_reply_to
    FOREIGN KEY (reply_to_id)
    REFERENCES comments(id)
    ON DELETE SET NULL; 