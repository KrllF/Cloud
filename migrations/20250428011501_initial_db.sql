-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS Users (
    id SERIAL PRIMARY KEY,
    ip TEXT NOT NULL,
    token_size BIGINT NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Users
-- +goose StatementEnd
