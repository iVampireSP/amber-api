-- +goose Up
ALTER TABLE chats MODIFY COLUMN assistant_id bigint unsigned DEFAULT NULL;
-- add assistant id to chat_messages_table after chat_id
ALTER TABLE chat_messages ADD COLUMN assistant_id bigint unsigned DEFAULT NULL AFTER chat_id;
ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_assistant_id_foreign FOREIGN KEY (assistant_id) REFERENCES assistants (id);
ALTER TABLE document_chunks DROP COLUMN chunked;
ALTER TABLE document_chunks ADD COLUMN vectorized boolean DEFAULT false AFTER library_id;
CREATE INDEX document_chunks_vectorized_index ON document_chunks (vectorized);

-- +goose Down

-- ALTER TABLE chats MODIFY COLUMN assistant_id bigint unsigned NOT NULL ; # 这个不能被回滚
ALTER TABLE chat_messages DROP FOREIGN KEY chat_messages_assistant_id_foreign;
ALTER TABLE chat_messages DROP COLUMN assistant_id;
ALTER TABLE document_chunks DROP COLUMN vectorized;
ALTER TABLE document_chunks ADD COLUMN chunked boolean DEFAULT false AFTER library_id;
CREATE INDEX document_chunks_chunked_index ON document_chunks (chunked);

