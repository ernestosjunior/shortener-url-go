-- +goose Up
-- +goose StatementBegin
CREATE TABLE short_urls(
    id SERIAL PRIMARY KEY,
    code VARCHAR(16) NOT NULL UNIQUE,
    url VARCHAR(2048) NOT NULL,
    clicks BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE short_urls;
-- +goose StatementEnd
