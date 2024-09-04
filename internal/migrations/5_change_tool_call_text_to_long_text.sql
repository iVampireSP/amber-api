-- +goose Up
-- change chat message content to long text
ALTER TABLE chat_messages MODIFY COLUMN tool_call longtext;


-- +goose Down
ALTER TABLE chat_messages MODIFY COLUMN tool_call text;