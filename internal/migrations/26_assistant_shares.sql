-- +goose up
CREATE TABLE assistant_favorites
(
    id                     SERIAL PRIMARY KEY,
    assistant_id           BIGINT UNSIGNED NOT NULL,
    user_id                VARCHAR(255)    NOT NULL,
    deleted                BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at             TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

-- 这里打算统一索引命名，以 idx 开头
CREATE INDEX idx_assistant_shares_id ON assistant_favorites (id);
CREATE INDEX idx_assistant_favorites_assistant_id ON assistant_favorites (assistant_id);
CREATE INDEX idx_assistant_favorites_user_id ON assistant_favorites (user_id);
CREATE INDEX idx_assistant_favorites_deleted ON assistant_favorites (deleted);

-- 外键
ALTER TABLE assistant_favorites
    ADD CONSTRAINT fk_assistant_favorites_library_id FOREIGN KEY (assistant_id) REFERENCES assistants (id);

-- add public to assistants table
ALTER TABLE assistants
    ADD COLUMN `public` BOOLEAN NOT NULL DEFAULT FALSE AFTER temperature;

-- +goose Down
ALTER TABLE assistant_favorites
    DROP FOREIGN KEY fk_assistant_favorites_library_id;
ALTER TABLE assistants
    DROP COLUMN `public`;
DROP TABLE assistant_favorites;
