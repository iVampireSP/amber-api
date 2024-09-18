-- +goose Up
-- tool call tokens 可以为一些异步任务提供支持
CREATE TABLE tool_call_tokens
(
    id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    chat_id    BIGINT UNSIGNED NOT NULL,
    token      VARCHAR(255)    NOT NULL,
    expired_at TIMESTAMP(3)    NOT NULL,
    created_at TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at TIMESTAMP(3)    NOT NULL DEFAULT CURRENT_TIMESTAMP(3)
);

-- index
CREATE INDEX tool_call_tokens_chat_id_index ON tool_call_tokens (chat_id);
CREATE INDEX tool_call_tokens_token_index ON tool_call_tokens (token);

-- +goose Down
DROP TABLE tool_call_tokens;