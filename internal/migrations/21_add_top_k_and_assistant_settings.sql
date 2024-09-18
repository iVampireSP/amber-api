-- +goose Up

-- add temperature
ALTER TABLE assistants ADD COLUMN temperature float DEFAULT 1 AFTER library_id;

-- +goose Down
ALTER TABLE assistants DROP COLUMN temperature;
