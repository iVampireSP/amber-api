-- +goose Up
-- create memories table
CREATE TABLE memories (
  id SERIAL PRIMARY KEY,
  user_id bigint NOT NULL,
  content TEXT NOT NULL,
  content_md5 VARCHAR(255) NOT NULL COMMENT 'md5',
  model VARCHAR(255) NOT NULL,
  metadata JSON DEFAULT NULL,
  vector JSON NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- idx
CREATE INDEX memories_user_id_index ON memories (user_id);
CREATE INDEX memories_hash_index ON memories (content_md5);
CREATE INDEX memories_model_index ON memories (model);
CREATE INDEX memories_created_at_index ON memories (created_at);


-- +goose Down
DROP TABLE memories;
