-- +goose Up
CREATE TABLE embeddings (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  text LONGTEXT DEFAULT NULL,
  file_id BIGINT unsigned DEFAULT NULL,
  text_md5 VARCHAR(255) DEFAULT NULL COMMENT 'md5',
  model VARCHAR(255) NOT NULL,
  vector JSON NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  INDEX embeddings_text_md5_index (text_md5),
  INDEX embeddings_model_index (model),
  INDEX embeddings_file_id_index (file_id),
  UNIQUE (text_md5)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

ALTER TABLE embeddings ADD CONSTRAINT embeddings_file_id_foreign FOREIGN KEY (file_id) REFERENCES files (id);


-- +goose Down
DROP TABLE embeddings;