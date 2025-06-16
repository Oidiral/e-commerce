-- +goose Up
-- +goose StatementBegin
DELETE FROM roles WHERE name = 'SERVICE_ORDER';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO roles (name) VALUES
    ('SERVICE_ORDER');
-- +goose StatementEnd
