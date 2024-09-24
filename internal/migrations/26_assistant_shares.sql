-- +goose up
CREATE TABLE favorite_assistants
(
    id                     SERIAL PRIMARY KEY,
    assistant_id           BIGINT UNSIGNED DEFAULT NULL,
    user_id                VARCHAR(255)    NOT NULL,
    deleted                BOOLEAN         NOT NULL DEFAULT FALSE,
    created_at             TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at             TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

-- 这里打算统一索引命名，以 idx 开头
CREATE INDEX idx_favorite_assistants_id ON favorite_assistants (id);
CREATE INDEX idx_favorite_assistants_assistant_id ON favorite_assistants (assistant_id);
CREATE INDEX idx_favorite_assistants_user_id ON favorite_assistants (user_id);
CREATE INDEX idx_favorite_assistants_deleted ON favorite_assistants (deleted);

-- 外键
ALTER TABLE favorite_assistants
    ADD CONSTRAINT fk_favorite_assistants_library_id FOREIGN KEY (assistant_id) REFERENCES assistants (id);

-- add public to assistants table
ALTER TABLE assistants
    ADD COLUMN `public` BOOLEAN NOT NULL DEFAULT FALSE AFTER temperature;

-- +goose Down
ALTER TABLE favorite_assistants
    DROP FOREIGN KEY fk_favorite_assistants_library_id;
ALTER TABLE assistants
    DROP COLUMN `public`;
DROP TABLE favorite_assistants;
