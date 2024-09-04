-- +goose Up
ALTER TABLE assistants ADD COLUMN disable_memory BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE assistants ADD COLUMN enable_memory_for_assistant_share BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE assistants DROP COLUMN disable_memory;
ALTER TABLE assistants DROP COLUMN enable_memory_for_assistant_share;
