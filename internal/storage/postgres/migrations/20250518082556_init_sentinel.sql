-- +goose Up
CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  data BYTEA NOT NULL,
  expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_sessions_expiry ON sessions USING btree (expiry);

CREATE TABLE token_sessions (
    session_id TEXT PRIMARY KEY,
    refresh_jti TEXT NOT NULL,
    expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_token_sessions_expiry ON token_sessions USING btree (expiry);
CREATE INDEX idx_token_sessions_refresh_jti ON token_sessions USING btree (refresh_jti);

CREATE TABLE auth_users (
    id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    name TEXT NOT NULL,
    username TEXT NOT NULL,
    picture TEXT NOT NULL,
    email TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE auth_accounts (
    id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    provider TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    email TEXT NOT NULL,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    access_token TEXT NOT NULL,
    expiry TIMESTAMPTZ NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_expiry TIMESTAMPTZ NOT NULL,
    token_type TEXT NOT NULL,
    id_token TEXT NOT NULL,
    name TEXT NOT NULL,
    preferred_username TEXT NOT NULL,
    nickname TEXT NOT NULL,
    picture TEXT NOT NULL,
    profile TEXT NOT NULL,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES auth_users(id) ON DELETE CASCADE,
    UNIQUE(provider, provider_id),
    UNIQUE(provider, user_id)
);


-- +goose Down
DROP TABLE sessions;
DROP TABLE token_sessions;
DROP TABLE auth_accounts;
DROP TABLE auth_users;
