-- +goose Up
RENAME TABLE assistant_api_keys TO assistant_keys;

-- rename index
ALTER TABLE assistant_keys RENAME INDEX assistant_api_keys_secret_index TO assistant_keys_secret_index;
ALTER TABLE assistant_keys RENAME INDEX assistant_api_keys_assistant_id_index TO assistant_keys_assistant_id_index;
-- rename foreign key
ALTER TABLE assistant_keys ADD CONSTRAINT assistant_keys_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE assistant_keys DROP FOREIGN KEY assistant_api_keys_assistant_id_foreign;

-- +goose Down
RENAME TABLE assistant_keys TO assistant_api_keys;
ALTER TABLE assistant_api_keys RENAME INDEX assistant_keys_secret_index TO assistant_api_keys_secret_index;
ALTER TABLE assistant_api_keys RENAME INDEX assistant_keys_assistant_id_index TO assistant_api_keys_assistant_id_index;
ALTER TABLE assistant_api_keys ADD CONSTRAINT assistant_api_keys_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE assistant_api_keys DROP FOREIGN KEY assistant_keys_assistant_id_foreign;