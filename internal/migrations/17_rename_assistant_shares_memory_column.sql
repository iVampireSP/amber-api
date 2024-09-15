-- +goose Up
-- rename enable_memory_for_assistant_share column on assistants table
ALTER TABLE assistants CHANGE COLUMN enable_memory_for_assistant_share enable_memory_for_assistant_api BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE assistants CHANGE COLUMN enable_memory_for_assistant_api enable_memory_for_assistant_share BOOLEAN NOT NULL DEFAULT FALSE;