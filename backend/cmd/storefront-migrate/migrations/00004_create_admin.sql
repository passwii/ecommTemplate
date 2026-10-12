-- +goose Up
CREATE TABLE admin_users (
    id TEXT PRIMARY KEY CHECK (id <> ''),
    email TEXT NOT NULL UNIQUE CHECK (email <> ''),
    password_hash TEXT NOT NULL CHECK (password_hash <> ''),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admin_sessions (
    token_hash TEXT PRIMARY KEY CHECK (token_hash <> ''),
    admin_user_id TEXT NOT NULL REFERENCES admin_users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX admin_sessions_user_idx ON admin_sessions (admin_user_id);
CREATE INDEX admin_sessions_expires_idx ON admin_sessions (expires_at);

-- +goose Down
DROP TABLE admin_sessions;
DROP TABLE admin_users;
