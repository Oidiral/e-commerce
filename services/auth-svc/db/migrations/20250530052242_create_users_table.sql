-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE IF NOT EXISTS auth.users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    status        SMALLINT    NOT NULL DEFAULT 1,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);

CREATE TABLE IF NOT EXISTS auth.roles (
    id   SMALLSERIAL PRIMARY KEY,
    name TEXT        NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS auth.user_roles (
    user_id UUID     REFERENCES auth.users(id) ON DELETE CASCADE,
    role_id SMALLINT REFERENCES auth.roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION auth.set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER trg_users_updated
BEFORE UPDATE ON auth.users
FOR EACH ROW
EXECUTE FUNCTION auth.set_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_updated ON auth.users;
DROP FUNCTION IF EXISTS auth.set_timestamp();
DROP TABLE IF EXISTS auth.user_roles;
DROP TABLE IF EXISTS auth.roles;
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth CASCADE;
