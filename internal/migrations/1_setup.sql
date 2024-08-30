-- 由于切换到了 Gorm, 所以将原来的 sql 全部合并为一个文件
-- +goose Up
CREATE TABLE tools
(
    id            bigint unsigned AUTO_INCREMENT,
    name          varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
    description   varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
    discovery_url varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
    api_key       varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
    data          json                                    DEFAULT NULL,
    user_id       bigint    NOT NULL,
    created_at    timestamp NULL                          DEFAULT NULL,
    updated_at    timestamp NULL                          DEFAULT NULL,
    PRIMARY KEY (id)
);

CREATE INDEX tools_user_id_index ON tools (user_id);

CREATE TABLE assistants
(
    id                     bigint unsigned AUTO_INCREMENT,
    name                   varchar(255)   DEFAULT NULL,
    description            varchar(255)   DEFAULT NULL,
    prompt                 text           DEFAULT NULL,
    disable_default_prompt boolean   NOT NULL,
    user_id                bigint unsigned NOT NULL,
    created_at             timestamp NULL DEFAULT NULL,
    updated_at             timestamp NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
CREATE INDEX assistants_user_id_index ON assistants (user_id);

CREATE TABLE assistant_tools
(
    id           bigint unsigned AUTO_INCREMENT,
    assistant_id bigint unsigned NOT NULL,
    tool_id      bigint unsigned NOT NULL,
    created_at   timestamp NULL DEFAULT NULL,
    updated_at   timestamp NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;



CREATE INDEX assistant_tools_tool_id_index ON assistant_tools (tool_id);
CREATE INDEX assistant_tools_assistant_id_index ON assistant_tools (assistant_id);

CREATE TABLE chats
(
    id           bigint unsigned AUTO_INCREMENT,
    name         varchar(255)   DEFAULT NULL,
    assistant_id bigint unsigned NOT NULL,
    user_id      bigint unsigned DEFAULT NULL,
    created_at   timestamp NULL DEFAULT NULL,
    updated_at   timestamp NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE INDEX chats_assistant_id_index ON chats (assistant_id);
CREATE INDEX chats_user_id_index ON chats (user_id);


CREATE TABLE chat_messages
(
    id                bigint unsigned AUTO_INCREMENT,
    chat_id           bigint unsigned NOT NULL,
    content           text      NOT NULL,
    role              varchar(255)       DEFAULT NULL,
    prompt_tokens     bigint    NOT NULL DEFAULT 0,
    completion_tokens bigint    NOT NULL DEFAULT 0,
    total_tokens      bigint    NOT NULL DEFAULT 0,
    created_at        timestamp NULL     DEFAULT NULL,
    updated_at        timestamp NULL     DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE INDEX chat_messages_chat_id_index ON chat_messages (chat_id);
CREATE INDEX chat_messages_created_at_index ON chat_messages (created_at);
CREATE INDEX chat_messages_role_index ON chat_messages (role);

CREATE TABLE assistant_shares
(
    id           bigint unsigned AUTO_INCREMENT,
    assistant_id bigint unsigned NOT NULL,
    token        varchar(255) NOT NULL,
    created_at   timestamp    NULL DEFAULT NULL,
    updated_at   timestamp    NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE INDEX assistant_shares_assistant_id_index ON assistant_shares (assistant_id);
CREATE INDEX assistant_shares_token_index ON assistant_shares (token);

ALTER TABLE chats
    add column expired_at timestamp NULL DEFAULT NULL AFTER user_id;
ALTER TABLE chats
    add column owner varchar(255) DEFAULT NULL AFTER user_id;

ALTER TABLE chats
    add column guest_id varchar(255) DEFAULT NULL AFTER user_id;
create index chats_expired_at_index on chats (expired_at);
create index chats_owner_index on chats (owner);
create index chats_guest_id_index on chats (guest_id);

CREATE TABLE files
(
    id         bigint unsigned AUTO_INCREMENT,
    url        varchar(255)   DEFAULT NULL,
    url_hash   varchar(255)   DEFAULT NULL,
    file_hash  varchar(255)   DEFAULT NULL,
    mime_type  varchar(255)   DEFAULT NULL,
    path       varchar(255)   DEFAULT NULL,
    expired_at timestamp NULL DEFAULT NULL,
    created_at timestamp NULL DEFAULT NULL,
    updated_at timestamp NULL DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

create index files_url_hash_index on files (url_hash);

create index files_file_hash_index on files (file_hash);

create index files_mime_type_index on files (mime_type);

create index files_expired_at_index on files (expired_at);


ALTER TABLE chat_messages
    ADD COLUMN hidden boolean NOT NULL DEFAULT false AFTER role;
create index chat_messages_hidden_index on chat_messages (hidden);
update chat_messages
set hidden = true
where role LIKE "%_hide";

UPDATE chat_messages
SET role = 'file'
WHERE role = 'image';

