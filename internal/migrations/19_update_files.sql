-- +goose Up

-- 大改存储方案，废除 user files, 去除 public 属性
ALTER TABLE chat_messages
    DROP FOREIGN KEY chat_messages_user_file_id_foreign;
ALTER TABLE chat_messages
    DROP COLUMN user_file_id;

DROP TABLE user_files;
-- 移除 files 表中的 public 字段
ALTER TABLE files
    DROP COLUMN public;

-- +goose Down
ALTER TABLE files
    ADD COLUMN public tinyint(1) DEFAULT 0;
CREATE INDEX files_public_index ON files (public);

CREATE TABLE user_files
(
    id         SERIAL PRIMARY KEY AUTO_INCREMENT,
    user_id    VARCHAR(255)    NOT NULL,
    file_id    bigint unsigned NOT NULL,
    created_at datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
CREATE INDEX user_files_user_id_index ON user_files (user_id);
CREATE INDEX user_files_file_id_index ON user_files (file_id);
ALTER TABLE user_files
    ADD CONSTRAINT user_files_file_id_foreign FOREIGN KEY (file_id) REFERENCES files (id);

ALTER TABLE chat_messages ADD COLUMN user_file_id bigint unsigned DEFAULT NULL AFTER hidden;
ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_user_file_id_foreign FOREIGN KEY (user_file_id) REFERENCES user_files (id);
ALTER TABLE chat_messages ADD COLUMN file_id bigint unsigned DEFAULT NULL AFTER tool_call;
