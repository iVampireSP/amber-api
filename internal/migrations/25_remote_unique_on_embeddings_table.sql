-- +goose Up
-- remove text_md5 unique index
DROP INDEX text_md5 ON embeddings;

-- +goose Down
CREATE UNIQUE INDEX text_md5 ON embeddings (text_md5);