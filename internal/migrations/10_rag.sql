-- +goose Up
CREATE TABLE libraries (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT NULL,
    user_id VARCHAR(255) NOT NULL,
    `default` boolean NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- index
CREATE INDEX library_user_id_index ON libraries (user_id);
CREATE INDEX library_default_index ON libraries (`default`);

-- create documents table
CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    name LONGTEXT NOT NULL,
    library_id BIGINT UNSIGNED NOT NULL,
    file_id BIGINT UNSIGNED DEFAULT NULL,
    file_hash  varchar(255) DEFAULT NULL,
    chunked boolean NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX documents_library_id_index ON documents (library_id);
CREATE INDEX documents_file_id_index ON documents (file_id);
CREATE INDEX documents_chunked_index ON documents (chunked);
CREATE INDEX documents_file_hash ON documents (file_hash);
ALTER TABLE documents ADD CONSTRAINT documents_library_id_foreign FOREIGN KEY (library_id) REFERENCES libraries (id);
ALTER TABLE documents ADD CONSTRAINT documents_file_id_foreign FOREIGN KEY (file_id) REFERENCES files (id);

-- create chunks table
CREATE TABLE document_chunks (
    id SERIAL PRIMARY KEY,
    content LONGTEXT NOT NULL,
    `order` INT NOT NULL,
    document_id BIGINT UNSIGNED NOT NULL,
    library_id BIGINT UNSIGNED NOT NULL,
    chunked boolean NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX document_chunks_library_id_index ON document_chunks (library_id);
CREATE INDEX document_chunks_document_id_index ON document_chunks (document_id);
CREATE INDEX document_chunks_chunked_index ON document_chunks (chunked);
-- order idx
CREATE INDEX document_chunks_order_index ON document_chunks (`order`);
ALTER TABLE document_chunks ADD CONSTRAINT document_chunks_library_id_foreign FOREIGN KEY (library_id) REFERENCES libraries (id);

-- +goose Down
DROP TABLE document_chunks;
DROP TABLE documents;
DROP TABLE libraries;
