-- +goose Up
-- drop documents_library_id_foreign
ALTER TABLE documents DROP FOREIGN KEY documents_file_id_foreign;
ALTER TABLE documents DROP COLUMN file_id;
ALTER TABLE documents DROP COLUMN file_hash;

-- +goose Down
ALTER TABLE documents ADD COLUMN file_id BIGINT UNSIGNED DEFAULT NULL AFTER library_id;
ALTER TABLE documents ADD COLUMN file_hash varchar(255) DEFAULT NULL AFTER file_id;
ALTER TABLE documents ADD CONSTRAINT documents_file_id_foreign FOREIGN KEY (file_id) REFERENCES files (id);