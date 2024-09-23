-- +goose up
CREATE TABLE assistant_shares
(
    id           SERIAL PRIMARY KEY,

    assistant_id BIGINT UNSIGNED NOT NULL,
    user_id      BIGINT UNSIGNED NOT NULL,
    -- total tokens
    total_tokens BIGINT UNSIGNED NOT NULL,
    created_at   datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 这里打算统一索引命名，以 idx 开头
CREATE INDEX idx_assistant_shares_id ON assistant_shares (id);
CREATE INDEX idx_assistant_shares_assistant_id_created_at ON assistant_shares (assistant_id);
CREATE INDEX idx_assistant_shares_user_id_created_at ON assistant_shares (user_id);

-- 外键
ALTER TABLE assistant_shares
    ADD CONSTRAINT fk_assistant_shares_assistant_id FOREIGN KEY (assistant_id) REFERENCES assistants (id) ON DELETE CASCADE;

-- 添加 assistant_share_id 到 assistants 表中，放到 library_id后面
ALTER TABLE assistants ADD COLUMN assistant_share_id BIGINT UNSIGNED DEFAULT NULL;
ALTER TABLE assistants ADD CONSTRAINT fk_assistants_assistant_share_id FOREIGN KEY (assistant_share_id) REFERENCES assistant_shares (id) ON DELETE SET NULL;

-- +goose Down
DROP TABLE assistant_shares;
