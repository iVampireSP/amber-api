-- +goose Up
ALTER TABLE assistants ADD COLUMN library_id bigint unsigned DEFAULT NULL AFTER user_id ;
-- index
CREATE INDEX assistants_library_id_index ON assistants (library_id);
ALTER TABLE assistants ADD CONSTRAINT assistants_library_id_foreign FOREIGN KEY (library_id) REFERENCES libraries (id);

-- +goose Down
-- drop foreign key
ALTER TABLE assistants DROP FOREIGN KEY assistants_library_id_foreign;
ALTER TABLE assistants DROP COLUMN library_id