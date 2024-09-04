-- +goose Up
-- add tool_calls and tool_call_id field after role
ALTER TABLE chat_messages ADD COLUMN tool_call varchar(255) DEFAULT NULL AFTER role;
ALTER TABLE chat_messages ADD COLUMN file_id bigint unsigned DEFAULT NULL AFTER tool_call;


-- +goose Down
ALTER TABLE chat_messages DROP COLUMN tool_call;
ALTER TABLE chat_messages DROP COLUMN file_id;