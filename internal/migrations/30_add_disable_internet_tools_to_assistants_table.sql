-- +goose Up

ALTER TABLE assistants
    ADD COLUMN disable_internet_search BOOLEAN NOT NULL DEFAULT false;

-- 禁用网页浏览
ALTER TABLE assistants
    ADD COLUMN disable_web_browsing BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE assistants
    DROP COLUMN disable_internet_search;
ALTER TABLE assistants
    DROP COLUMN disable_web_browsing;