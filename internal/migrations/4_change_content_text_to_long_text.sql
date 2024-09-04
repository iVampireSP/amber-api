-- +goose Up
-- change chat message content to long text
ALTER TABLE chat_messages MODIFY COLUMN content longtext;


-- +goose Down
ALTER TABLE chat_messages MODIFY COLUMN content text;