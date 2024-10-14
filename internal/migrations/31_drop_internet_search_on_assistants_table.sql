-- +goose Up
ALTER TABLE assistants
    DROP COLUMN disable_internet_search;

-- +goose Down

ALTER TABLE assistants
    ADD COLUMN disable_internet_search BOOLEAN NOT NULL DEFAULT false;
