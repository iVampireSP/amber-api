-- +goose Up
-- rename assistant_shares to assistant_tokens
RENAME TABLE assistant_shares TO assistant_tokens;
-- +goose Down
RENAME TABLE assistant_tokens TO assistant_shares;