-- +goose up
CREATE TABLE unsettled_tokens (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id VARCHAR(255) NOT NULL,
    count BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- count index
CREATE INDEX unsettled_tokens_count_idx ON unsettled_tokens (count);

-- +goose Down
DROP TABLE unsettled_tokens;