-- +goose Up
-- +goose StatementBegin
INSERT INTO auth.roles (name) VALUES
    ('admin'),
    ('user')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM auth.user_roles WHERE role_id IN (
    SELECT id FROM auth.roles WHERE name IN ('user','admin')
);
DELETE FROM auth.roles WHERE name IN ('user','admin');
-- +goose StatementEnd
