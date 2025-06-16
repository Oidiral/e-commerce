-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    secret_hash TEXT NOT NULL,
    roles TEXT[] NOT NULL DEFAULT '{}',
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd
