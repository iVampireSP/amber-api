-- +goose Up
-- 将 assistant_keys 表中的 created_at 字段改为 timestamp(3)
ALTER TABLE assistant_keys
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE assistant_keys
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- assistant_tools
ALTER TABLE assistant_tools
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE assistant_tools
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);


-- assistants
ALTER TABLE assistants
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE assistants
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- chat_messages
ALTER TABLE chat_messages
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE chat_messages
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- chats
ALTER TABLE chats
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE chats
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- document_chunks
ALTER TABLE document_chunks
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE document_chunks
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- documents
ALTER TABLE documents
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE documents
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- embeddings
ALTER TABLE embeddings
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE embeddings
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- files
ALTER TABLE files
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE files
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- libraries
ALTER TABLE libraries
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE libraries
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- memories
ALTER TABLE memories
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE memories
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- tools
ALTER TABLE tools
    MODIFY COLUMN created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE tools
    MODIFY COLUMN updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3);

-- add index
CREATE INDEX assistant_keys_created_at_index ON assistant_keys (created_at);
CREATE INDEX assistant_tools_created_at_index ON assistant_tools (created_at);
CREATE INDEX assistants_created_at_index ON assistants (created_at);
-- CREATE INDEX chat_messages_created_at_index ON chat_messages (created_at);
CREATE INDEX chats_created_at_index ON chats (created_at);
CREATE INDEX document_chunks_created_at_index ON document_chunks (created_at);
CREATE INDEX documents_created_at_index ON documents (created_at);
CREATE INDEX embeddings_created_at_index ON embeddings (created_at);
CREATE INDEX files_created_at_index ON files (created_at);
CREATE INDEX libraries_created_at_index ON libraries (created_at);
CREATE INDEX tools_created_at_index ON tools (created_at);

-- +goose Down

-- Reverting the timestamp(3) changes back to the original timestamp

-- assistant_keys
ALTER TABLE assistant_keys
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE assistant_keys
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- assistant_tools
ALTER TABLE assistant_tools
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE assistant_tools
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- assistants
ALTER TABLE assistants
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE assistants
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- chat_messages
ALTER TABLE chat_messages
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE chat_messages
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- chats
ALTER TABLE chats
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE chats
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- document_chunks
ALTER TABLE document_chunks
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE document_chunks
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- documents
ALTER TABLE documents
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE documents
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- embeddings
ALTER TABLE embeddings
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE embeddings
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- files
ALTER TABLE files
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE files
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- libraries
ALTER TABLE libraries
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE libraries
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- memories
ALTER TABLE memories
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE memories
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- tools
ALTER TABLE tools
    MODIFY COLUMN created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE tools
    MODIFY COLUMN updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;


-- drop index
DROP INDEX assistant_keys_created_at_index ON assistant_keys;
DROP INDEX assistant_tools_created_at_index ON assistant_tools;
DROP INDEX assistants_created_at_index ON assistants;
-- DROP INDEX chat_messages_created_at_index ON chat_messages;
DROP INDEX chats_created_at_index ON chats;
DROP INDEX document_chunks_created_at_index ON document_chunks;
DROP INDEX documents_created_at_index ON documents;
DROP INDEX embeddings_created_at_index ON embeddings;
DROP INDEX files_created_at_index ON files;
DROP INDEX libraries_created_at_index ON libraries;
DROP INDEX tools_created_at_index ON tools;

