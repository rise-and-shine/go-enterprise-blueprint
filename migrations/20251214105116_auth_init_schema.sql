-- +goose Up
-- +goose StatementBegin
-- Enable pg_uuidv7 extension for UUID v7 generation
CREATE EXTENSION IF NOT EXISTS pg_uuidv7;

CREATE TABLE auth.admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7 (),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_superadmin BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_active_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE TABLE auth.roles (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE TABLE auth.role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role_id BIGINT NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_role_permissions_role_id ON auth.role_permissions (role_id);

CREATE TABLE auth.actor_roles (
    id BIGSERIAL PRIMARY KEY,
    actor_type TEXT NOT NULL,
    actor_id UUID NOT NULL,
    role_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_actor_roles_actor ON auth.actor_roles (actor_type, actor_id);

CREATE INDEX idx_actor_roles_role_id ON auth.actor_roles (role_id);

CREATE TABLE auth.actor_permissions (
    id BIGSERIAL PRIMARY KEY,
    actor_type TEXT NOT NULL,
    actor_id UUID NOT NULL,
    permission TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_actor_permissions_actor ON auth.actor_permissions (actor_type, actor_id);

CREATE TABLE auth.sessions (
    id BIGSERIAL PRIMARY KEY,
    actor_type TEXT NOT NULL,
    actor_id UUID NOT NULL,
    access_token TEXT NOT NULL,
    access_token_expires_at TIMESTAMPTZ NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    last_used_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

CREATE INDEX idx_sessions_actor ON auth.sessions (actor_type, actor_id);

CREATE INDEX idx_sessions_access_token ON auth.sessions (access_token);

CREATE INDEX idx_sessions_refresh_token ON auth.sessions (refresh_token);

ALTER TABLE auth.role_permissions ADD CONSTRAINT fk_role_permissions_role FOREIGN KEY (role_id) REFERENCES auth.roles (id) ON DELETE CASCADE;

ALTER TABLE auth.actor_roles ADD CONSTRAINT fk_actor_roles_role FOREIGN KEY (role_id) REFERENCES auth.roles (id) ON DELETE CASCADE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- Drop foreign keys first
ALTER TABLE IF EXISTS auth.actor_roles
DROP CONSTRAINT IF EXISTS fk_actor_roles_role;

ALTER TABLE IF EXISTS auth.role_permissions
DROP CONSTRAINT IF EXISTS fk_role_permissions_role;

-- Drop tables in reverse order
DROP TABLE IF EXISTS auth.sessions;

DROP TABLE IF EXISTS auth.actor_permissions;

DROP TABLE IF EXISTS auth.actor_roles;

DROP TABLE IF EXISTS auth.role_permissions;

DROP TABLE IF EXISTS auth.roles;

DROP TABLE IF EXISTS auth.admins;

-- Drop extension
DROP EXTENSION IF EXISTS pg_uuidv7;

-- +goose StatementEnd
