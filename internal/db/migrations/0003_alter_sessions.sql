-- +goose Up
-- +goose statementbegin

DROP TABLE IF EXISTS sessions;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;

CREATE TABLE IF NOT EXISTS sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

-- +goose statementend


-- +goose Down
-- +goose statementbegin
DROP TABLE IF EXISTS sessions;

-- +goose statementend
