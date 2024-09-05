-- +goose Up
-- add size to files
ALTER TABLE files ADD COLUMN size bigint unsigned DEFAULT NULL AFTER path;
CREATE INDEX files_size_index ON files (size);

-- +goose Down
ALTER TABLE files DROP COLUMN size;

