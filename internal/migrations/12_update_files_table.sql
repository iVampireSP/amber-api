-- +goose Up
-- add size to files
ALTER TABLE files ADD COLUMN size bigint unsigned DEFAULT NULL AFTER path;
CREATE INDEX files_size_index ON files (size);

-- add public to files
ALTER TABLE files ADD COLUMN public tinyint unsigned DEFAULT 0 AFTER path;
CREATE INDEX files_public_index ON files (public);

-- +goose Down
ALTER TABLE files DROP COLUMN size;
ALTER TABLE files DROP COLUMN public;

