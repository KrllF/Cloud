-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users (
    ip TEXT PRIMARY KEY,
    token_size BIGINT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Users
-- +goose StatementEnd
