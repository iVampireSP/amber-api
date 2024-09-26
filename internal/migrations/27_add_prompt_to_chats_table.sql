-- +goose up
ALTER TABLE chats ADD COLUMN prompt TEXT DEFAULT NULL AFTER name;

-- +goose down
ALTER TABLE chats DROP COLUMN prompt;