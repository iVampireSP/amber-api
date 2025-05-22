-- +goose Up

-- 创建所需序列
CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS user_tools_id_seq;
CREATE SEQUENCE IF NOT EXISTS chat_messages_id_seq;
CREATE SEQUENCE IF NOT EXISTS chats_id_seq;
CREATE SEQUENCE IF NOT EXISTS document_chunks_id_seq;
CREATE SEQUENCE IF NOT EXISTS documents_id_seq;
CREATE SEQUENCE IF NOT EXISTS embeddings_id_seq;
CREATE SEQUENCE IF NOT EXISTS files_id_seq;
CREATE SEQUENCE IF NOT EXISTS goose_db_version_id_seq;
CREATE SEQUENCE IF NOT EXISTS libraries_id_seq;
CREATE SEQUENCE IF NOT EXISTS memories_id_seq;
CREATE SEQUENCE IF NOT EXISTS message_blocks_id_seq;
CREATE SEQUENCE IF NOT EXISTS tools_id_seq;

-- ----------------------------
-- Table structure for users (新增用户表)
-- ----------------------------
DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users (
  id bigint NOT NULL DEFAULT nextval('users_id_seq'),
  username varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password_hash varchar(255) NOT NULL,
  avatar varchar(255),
  status varchar(50) NOT NULL DEFAULT 'active',
  last_login_at timestamp,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_status ON users (status);
CREATE INDEX idx_users_created_at ON users (created_at);

-- ----------------------------
-- Table structure for libraries (需要先创建，因为有外键引用)
-- ----------------------------
DROP TABLE IF EXISTS libraries CASCADE;
CREATE TABLE libraries (
  id bigint NOT NULL DEFAULT nextval('libraries_id_seq'),
  name varchar(255) NOT NULL,
  description varchar(255),
  user_id varchar(255) NOT NULL,
  "default" boolean NOT NULL DEFAULT false,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_libraries_user_id ON libraries (user_id);
CREATE INDEX idx_libraries_default ON libraries ("default");
CREATE INDEX idx_libraries_created_at ON libraries (created_at);

-- ----------------------------
-- Table structure for files (需要先创建，因为有外键引用)
-- ----------------------------
DROP TABLE IF EXISTS files CASCADE;
CREATE TABLE files (
  id bigint NOT NULL DEFAULT nextval('files_id_seq'),
  url varchar(255),
  url_hash varchar(255),
  file_hash varchar(255),
  mime_type varchar(255),
  path varchar(255),
  size bigint,
  expired_at timestamp,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_files_url_hash ON files (url_hash);
CREATE INDEX idx_files_file_hash ON files (file_hash);
CREATE INDEX idx_files_mime_type ON files (mime_type);
CREATE INDEX idx_files_expired_at ON files (expired_at);
CREATE INDEX idx_files_size ON files (size);
CREATE INDEX idx_files_created_at ON files (created_at);

-- ----------------------------
-- Table structure for tools (需要先创建，因为有外键引用)
-- ----------------------------
DROP TABLE IF EXISTS tools CASCADE;
CREATE TABLE tools (
  id bigint NOT NULL DEFAULT nextval('tools_id_seq'),
  name varchar(255),
  description varchar(255),
  discovery_url varchar(255),
  api_key varchar(255),
  data jsonb,
  user_id varchar(255) NOT NULL,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_tools_user_id ON tools (user_id);
CREATE INDEX idx_tools_created_at ON tools (created_at);

-- ----------------------------
-- Table structure for user_tools (替代assistant_tools)
-- ----------------------------
DROP TABLE IF EXISTS user_tools CASCADE;
CREATE TABLE user_tools (
  id bigint NOT NULL DEFAULT nextval('user_tools_id_seq'),
  user_id varchar(255) NOT NULL,
  tool_id bigint NOT NULL,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_user_tools_tool_id FOREIGN KEY (tool_id) REFERENCES tools (id)
);
CREATE INDEX idx_user_tools_tool_id ON user_tools (tool_id);
CREATE INDEX idx_user_tools_user_id ON user_tools (user_id);
CREATE INDEX idx_user_tools_created_at ON user_tools (created_at);

-- ----------------------------
-- Table structure for chats (需要先创建，因为有外键引用)
-- ----------------------------
DROP TABLE IF EXISTS chats CASCADE;
CREATE TABLE chats (
  id bigint NOT NULL DEFAULT nextval('chats_id_seq'),
  name varchar(255),
  prompt text,
  user_id varchar(255) NOT NULL,
  guest_id varchar(255),
  owner varchar(255),
  expired_at timestamp,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_chats_user_id ON chats (user_id);
CREATE INDEX idx_chats_expired_at ON chats (expired_at);
CREATE INDEX idx_chats_owner ON chats (owner);
CREATE INDEX idx_chats_guest_id ON chats (guest_id);
CREATE INDEX idx_chats_created_at ON chats (created_at);

-- ----------------------------
-- Table structure for chat_messages
-- ----------------------------
DROP TABLE IF EXISTS chat_messages CASCADE;
CREATE TABLE chat_messages (
  id bigint NOT NULL DEFAULT nextval('chat_messages_id_seq'),
  chat_id bigint NOT NULL,
  content text,
  role varchar(255),
  tool_call text,
  file_id bigint,
  hidden boolean NOT NULL DEFAULT false,
  prompt_tokens bigint NOT NULL DEFAULT 0,
  completion_tokens bigint NOT NULL DEFAULT 0,
  total_tokens bigint NOT NULL DEFAULT 0,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_chat_messages_chat_id FOREIGN KEY (chat_id) REFERENCES chats (id),
  CONSTRAINT fk_chat_messages_file_id FOREIGN KEY (file_id) REFERENCES files (id)
);
CREATE INDEX idx_chat_messages_chat_id ON chat_messages (chat_id);
CREATE INDEX idx_chat_messages_created_at ON chat_messages (created_at);
CREATE INDEX idx_chat_messages_role ON chat_messages (role);
CREATE INDEX idx_chat_messages_hidden ON chat_messages (hidden);

-- ----------------------------
-- Table structure for documents
-- ----------------------------
DROP TABLE IF EXISTS documents CASCADE;
CREATE TABLE documents (
  id bigint NOT NULL DEFAULT nextval('documents_id_seq'),
  name text NOT NULL,
  library_id bigint NOT NULL,
  chunked boolean NOT NULL DEFAULT false,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_documents_library_id FOREIGN KEY (library_id) REFERENCES libraries (id)
);
CREATE INDEX idx_documents_library_id ON documents (library_id);
CREATE INDEX idx_documents_chunked ON documents (chunked);
CREATE INDEX idx_documents_created_at ON documents (created_at);

-- ----------------------------
-- Table structure for document_chunks
-- ----------------------------
DROP TABLE IF EXISTS document_chunks CASCADE;
CREATE TABLE document_chunks (
  id bigint NOT NULL DEFAULT nextval('document_chunks_id_seq'),
  content text NOT NULL,
  "order" int NOT NULL,
  document_id bigint NOT NULL,
  library_id bigint NOT NULL,
  vectorized boolean DEFAULT false,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_document_chunks_library_id FOREIGN KEY (library_id) REFERENCES libraries (id),
  CONSTRAINT fk_document_chunks_document_id FOREIGN KEY (document_id) REFERENCES documents (id)
);
CREATE INDEX idx_document_chunks_library_id ON document_chunks (library_id);
CREATE INDEX idx_document_chunks_document_id ON document_chunks (document_id);
CREATE INDEX idx_document_chunks_order ON document_chunks ("order");
CREATE INDEX idx_document_chunks_vectorized ON document_chunks (vectorized);
CREATE INDEX idx_document_chunks_created_at ON document_chunks (created_at);

-- ----------------------------
-- Table structure for embeddings
-- ----------------------------
DROP TABLE IF EXISTS embeddings CASCADE;
CREATE TABLE embeddings (
  id bigint NOT NULL DEFAULT nextval('embeddings_id_seq'),
  text text,
  file_id bigint,
  text_md5 varchar(255),
  model varchar(255) NOT NULL,
  vector jsonb NOT NULL,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_embeddings_file_id FOREIGN KEY (file_id) REFERENCES files (id)
);
CREATE INDEX idx_embeddings_text_md5 ON embeddings (text_md5);
CREATE INDEX idx_embeddings_model ON embeddings (model);
CREATE INDEX idx_embeddings_file_id ON embeddings (file_id);
CREATE INDEX idx_embeddings_created_at ON embeddings (created_at);

-- ----------------------------
-- Table structure for goose_db_version
-- ----------------------------
DROP TABLE IF EXISTS goose_db_version CASCADE;
CREATE TABLE goose_db_version (
  id bigint NOT NULL DEFAULT nextval('goose_db_version_id_seq'),
  version_id bigint NOT NULL,
  is_applied boolean NOT NULL,
  tstamp timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);

-- ----------------------------
-- Table structure for memories
-- ----------------------------
DROP TABLE IF EXISTS memories CASCADE;
CREATE TABLE memories (
  id bigint NOT NULL DEFAULT nextval('memories_id_seq'),
  user_id varchar(255) NOT NULL,
  content text NOT NULL,
  content_md5 varchar(255) NOT NULL,
  model varchar(255) NOT NULL,
  metadata jsonb,
  vector jsonb NOT NULL,
  created_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
);
CREATE INDEX idx_memories_user_id ON memories (user_id);
CREATE INDEX idx_memories_content_md5 ON memories (content_md5);
CREATE INDEX idx_memories_model ON memories (model);
CREATE INDEX idx_memories_created_at ON memories (created_at);

-- ----------------------------
-- Table structure for message_blocks
-- ----------------------------
DROP TABLE IF EXISTS message_blocks CASCADE;
CREATE TABLE message_blocks (
  id bigint NOT NULL DEFAULT nextval('message_blocks_id_seq'),
  chat_id bigint NOT NULL,
  hash varchar(255) NOT NULL,
  full_content text NOT NULL,
  messages jsonb NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_message_blocks_chat_id FOREIGN KEY (chat_id) REFERENCES chats (id)
);
CREATE INDEX idx_message_blocks_chat_id ON message_blocks (chat_id);
CREATE INDEX idx_message_blocks_hash ON message_blocks (hash);

-- 设置序列的所有权
ALTER SEQUENCE users_id_seq OWNED BY users.id;
ALTER SEQUENCE user_tools_id_seq OWNED BY user_tools.id;
ALTER SEQUENCE chat_messages_id_seq OWNED BY chat_messages.id;
ALTER SEQUENCE chats_id_seq OWNED BY chats.id;
ALTER SEQUENCE document_chunks_id_seq OWNED BY document_chunks.id;
ALTER SEQUENCE documents_id_seq OWNED BY documents.id;
ALTER SEQUENCE embeddings_id_seq OWNED BY embeddings.id;
ALTER SEQUENCE files_id_seq OWNED BY files.id;
ALTER SEQUENCE goose_db_version_id_seq OWNED BY goose_db_version.id;
ALTER SEQUENCE libraries_id_seq OWNED BY libraries.id;
ALTER SEQUENCE memories_id_seq OWNED BY memories.id;
ALTER SEQUENCE message_blocks_id_seq OWNED BY message_blocks.id;
ALTER SEQUENCE tools_id_seq OWNED BY tools.id;

-- +goose Down
-- 删除序列
DROP TABLE IF EXISTS user_tools CASCADE;
DROP TABLE IF EXISTS chat_messages CASCADE;
DROP TABLE IF EXISTS message_blocks CASCADE;
DROP TABLE IF EXISTS chats CASCADE;
DROP TABLE IF EXISTS document_chunks CASCADE;
DROP TABLE IF EXISTS documents CASCADE;
DROP TABLE IF EXISTS memories CASCADE;
DROP TABLE IF EXISTS embeddings CASCADE;
DROP TABLE IF EXISTS files CASCADE;
DROP TABLE IF EXISTS tools CASCADE;
DROP TABLE IF EXISTS libraries CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS goose_db_version CASCADE;

DROP SEQUENCE IF EXISTS user_tools_id_seq;
DROP SEQUENCE IF EXISTS users_id_seq;
DROP SEQUENCE IF EXISTS chat_messages_id_seq;
DROP SEQUENCE IF EXISTS chats_id_seq;
DROP SEQUENCE IF EXISTS document_chunks_id_seq;
DROP SEQUENCE IF EXISTS documents_id_seq;
DROP SEQUENCE IF EXISTS embeddings_id_seq;
DROP SEQUENCE IF EXISTS files_id_seq;
DROP SEQUENCE IF EXISTS goose_db_version_id_seq;
DROP SEQUENCE IF EXISTS libraries_id_seq;
DROP SEQUENCE IF EXISTS memories_id_seq;
DROP SEQUENCE IF EXISTS message_blocks_id_seq;
DROP SEQUENCE IF EXISTS tools_id_seq;
