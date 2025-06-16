-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name) VALUES
    ('admin'),
    ('model');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_roles WHERE role_id IN (
    SELECT id FROM roles WHERE name IN ('model', 'admin')
);
DELETE FROM roles WHERE name IN ('model', 'admin');
-- +goose StatementEnd
