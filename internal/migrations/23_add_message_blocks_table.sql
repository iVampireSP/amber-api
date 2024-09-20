-- +goose Up
CREATE TABLE message_blocks
(
    id           SERIAL PRIMARY KEY,
    chat_id      BIGINT UNSIGNED NOT NULL,
    hash         VARCHAR(255)    NOT NULL,
    full_content LONGTEXT        NOT NULL,
    messages     JSON            NOT NULL,
    created_at   datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   datetime        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE INDEX message_blocks_hash_idx ON message_blocks (hash);

-- +goose Down
DROP TABLE message_blocks;
