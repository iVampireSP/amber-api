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
CREATE SEQUENCE IF NOT EXISTS libraries_id_seq;
CREATE SEQUENCE IF NOT EXISTS memories_id_seq;
CREATE SEQUENCE IF NOT EXISTS message_blocks_id_seq;
CREATE SEQUENCE IF NOT EXISTS tools_id_seq;

-- ----------------------------
-- Table structure for users
-- ----------------------------
CREATE TABLE users (
  id BIGINT PRIMARY KEY DEFAULT nextval('users_id_seq'),
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  avatar VARCHAR(255),
  status VARCHAR(50) DEFAULT 'active',
  last_login_at TIMESTAMP,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON users (username, status, created_at);

-- ----------------------------
-- Table structure for libraries
-- ----------------------------
-- 修正 libraries 表结构
CREATE TABLE libraries (
  id BIGINT PRIMARY KEY DEFAULT nextval('libraries_id_seq'),
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  "default" BOOLEAN DEFAULT false,  -- 转义关键字
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON libraries (user_id, "default", created_at);

-- ----------------------------
-- Table structure for files
-- ----------------------------
CREATE TABLE files (
  id BIGINT PRIMARY KEY DEFAULT nextval('files_id_seq'),
  url VARCHAR(255),
  url_hash VARCHAR(255),
  file_hash VARCHAR(255),
  mime_type VARCHAR(255),
  path VARCHAR(255),
  size BIGINT,
  expired_at TIMESTAMP,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON files (url_hash, file_hash, mime_type, expired_at, size, created_at);

-- ----------------------------
-- Table structure for tools
-- ----------------------------
CREATE TABLE tools (
  id BIGINT PRIMARY KEY DEFAULT nextval('tools_id_seq'),
  name VARCHAR(255),
  description VARCHAR(255),
  discovery_url VARCHAR(255),
  api_key VARCHAR(255),
  data JSONB,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON tools (user_id, created_at);

-- ----------------------------
-- Table structure for user_tools
-- ----------------------------
CREATE TABLE user_tools (
  id BIGINT PRIMARY KEY DEFAULT nextval('user_tools_id_seq'),
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tool_id BIGINT NOT NULL REFERENCES tools(id) ON DELETE CASCADE,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON user_tools (tool_id, user_id, created_at);

-- ----------------------------
-- Table structure for chats
-- ----------------------------
CREATE TABLE chats (
  id BIGINT PRIMARY KEY DEFAULT nextval('chats_id_seq'),
  name VARCHAR(255),
  prompt TEXT,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  guest_id VARCHAR(255),
  owner VARCHAR(255),
  expired_at TIMESTAMP,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON chats (user_id, expired_at, owner, guest_id, created_at);

-- ----------------------------
-- Table structure for chat_messages
-- ----------------------------
CREATE TABLE chat_messages (
  id BIGINT PRIMARY KEY DEFAULT nextval('chat_messages_id_seq'),
  chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
  content TEXT,
  role VARCHAR(255),
  tool_call TEXT,
  file_id BIGINT REFERENCES files(id) ON DELETE SET NULL,
  hidden BOOLEAN DEFAULT false,
  prompt_tokens BIGINT DEFAULT 0,
  completion_tokens BIGINT DEFAULT 0,
  total_tokens BIGINT DEFAULT 0,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON chat_messages (chat_id, created_at, role, hidden);

-- ----------------------------
-- Table structure for documents
-- ----------------------------
CREATE TABLE documents (
  id BIGINT PRIMARY KEY DEFAULT nextval('documents_id_seq'),
  name TEXT NOT NULL,
  library_id BIGINT NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
  chunked BOOLEAN DEFAULT false,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON documents (library_id, chunked, created_at);

-- ----------------------------
-- Table structure for document_chunks
-- ----------------------------
CREATE TABLE document_chunks (
  id BIGINT PRIMARY KEY DEFAULT nextval('document_chunks_id_seq'),
  content TEXT NOT NULL,
  "order" INT NOT NULL,
  document_id BIGINT NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
  library_id BIGINT NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
  vectorized BOOLEAN DEFAULT false,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON document_chunks (library_id, document_id, "order", vectorized, created_at);

-- ----------------------------
-- Table structure for embeddings
-- ----------------------------
CREATE TABLE embeddings (
  id BIGINT PRIMARY KEY DEFAULT nextval('embeddings_id_seq'),
  text TEXT,
  file_id BIGINT REFERENCES files(id) ON DELETE SET NULL,
  text_md5 VARCHAR(255),
  model VARCHAR(255) NOT NULL,
  vector JSONB NOT NULL,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON embeddings (text_md5, model, file_id, created_at);

-- ----------------------------
-- Table structure for memories
-- ----------------------------
CREATE TABLE memories (
  id BIGINT PRIMARY KEY DEFAULT nextval('memories_id_seq'),
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  content_md5 VARCHAR(255) NOT NULL,
  model VARCHAR(255) NOT NULL,
  metadata JSONB,
  vector JSONB NOT NULL,
  created_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON memories (user_id, content_md5, model, created_at);

-- ----------------------------
-- Table structure for message_blocks
-- ----------------------------
CREATE TABLE message_blocks (
  id BIGINT PRIMARY KEY DEFAULT nextval('message_blocks_id_seq'),
  chat_id BIGINT NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
  hash VARCHAR(255) NOT NULL,
  full_content TEXT NOT NULL,
  messages JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX ON message_blocks (chat_id, hash);

-- +goose Down
DROP TABLE message_blocks CASCADE;
DROP TABLE memories CASCADE;
DROP TABLE embeddings CASCADE;
DROP TABLE document_chunks CASCADE;
DROP TABLE documents CASCADE;
DROP TABLE chat_messages CASCADE;
DROP TABLE chats CASCADE;
DROP TABLE user_tools CASCADE;
DROP TABLE tools CASCADE;
DROP TABLE files CASCADE;
DROP TABLE libraries CASCADE;
DROP TABLE users CASCADE;

DROP SEQUENCE users_id_seq, user_tools_id_seq, chat_messages_id_seq, chats_id_seq,
document_chunks_id_seq, documents_id_seq, embeddings_id_seq, files_id_seq, libraries_id_seq, memories_id_seq, message_blocks_id_seq,
tools_id_seq;