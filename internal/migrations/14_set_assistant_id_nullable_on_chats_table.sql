-- +goose Up
ALTER TABLE chats MODIFY COLUMN assistant_id bigint unsigned DEFAULT NULL;
-- add assistant id to chat_messages_table after chat_id
ALTER TABLE chat_messages ADD COLUMN assistant_id bigint unsigned DEFAULT NULL AFTER chat_id;
ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);

-- +goose Down
# ALTER TABLE chats MODIFY COLUMN assistant_id bigint unsigned NOT NULL ; # not nullable
ALTER TABLE chat_messages DROP FOREIGN KEY chat_messages_assistant_id_foreign;
ALTER TABLE chat_messages DROP COLUMN assistant_id;
