-- +goose up
ALTER TABLE assistants ADD COLUMN total_token_usage BIGINT DEFAULT 0 AFTER user_id;

-- +goose down
ALTER TABLE assistants DROP COLUMN total_token_usage;