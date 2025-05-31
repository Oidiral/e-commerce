-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name) VALUES
    ('admin'),
    ('user');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_roles WHERE role_id IN (
    SELECT id FROM roles WHERE name IN ('user', 'admin')
);
DELETE FROM roles WHERE name IN ('user', 'admin');
-- +goose StatementEnd
