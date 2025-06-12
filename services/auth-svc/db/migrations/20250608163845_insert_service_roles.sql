-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (name) VALUES
    ('SERVICE_ORDER');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM user_roles WHERE role_id IN (
    SELECT id FROM roles WHERE name = 'SERVICE_ORDER'
);
-- +goose StatementEnd
