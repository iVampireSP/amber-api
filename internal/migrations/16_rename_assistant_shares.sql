-- +goose Up
RENAME TABLE assistant_shares TO assistant_api_keys;
ALTER TABLE assistant_api_keys CHANGE COLUMN token secret varchar(255) NOT NULL;
-- rename index
ALTER TABLE assistant_api_keys RENAME INDEX assistant_shares_token_index TO assistant_api_keys_secret_index;
ALTER TABLE assistant_api_keys RENAME INDEX assistant_shares_assistant_id_index TO assistant_api_keys_assistant_id_index;
-- rename foreign key
ALTER TABLE assistant_api_keys ADD CONSTRAINT assistant_api_keys_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE assistant_api_keys DROP FOREIGN KEY assistant_shares_assistant_id_foreign;

-- +goose Down
RENAME TABLE assistant_api_keys TO assistant_shares;
ALTER TABLE assistant_shares CHANGE COLUMN secret token varchar(255) NOT NULL;
ALTER TABLE assistant_shares RENAME INDEX assistant_api_keys_secret_index TO assistant_shares_token_index;
ALTER TABLE assistant_shares RENAME INDEX assistant_api_keys_assistant_id_index TO assistant_shares_assistant_id_index;
ALTER TABLE assistant_shares ADD CONSTRAINT assistant_shares_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE assistant_shares DROP FOREIGN KEY assistant_api_keys_assistant_id_foreign;