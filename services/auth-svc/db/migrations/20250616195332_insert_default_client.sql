-- +goose Up
-- +goose StatementBegin
INSERT INTO clients (id, secret_hash, roles)
VALUES ('cart-svc', '$2a$10$/98/2JtQb2xICYSPulbh4edAakqw0wFoRYnSt5m3059kogF5nKCMy', ARRAY['SERVICE_ORDER']);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM clients WHERE id = 'cart-service';
-- +goose StatementEnd
