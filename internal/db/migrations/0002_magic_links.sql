-- +goose Up
-- +goose statementbegin

-- Drop existing otp_sessions and its trigger [1]
DROP TABLE IF EXISTS otp_sessions;
DROP TRIGGER IF EXISTS update_otp_sessions_updated_at ON otp_sessions;

DROP TYPE IF EXISTS magic_link_status;
-- Create enum type for status

DROP TYPE IF EXISTS magic_link_status;
CREATE TYPE magic_link_status AS ENUM ('pending', 'used', 'expired', 'revoked');

-- Create magic_link_sessions table
CREATE TABLE IF NOT EXISTS magic_link_sessions(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    token TEXT UNIQUE NOT NULL,
    email VARCHAR(360) NOT NULL,
    pubkey TEXT NOT NULL,
    ip_address INET NOT NULL,
    user_agent TEXT,
    attempts SMALLINT NOT NULL DEFAULT 1,
    status magic_link_status NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_magic_link_token ON magic_link_sessions(token);
CREATE INDEX idx_magic_link_email ON magic_link_sessions(email);
CREATE INDEX idx_magic_link_expires_at ON magic_link_sessions(expires_at);
CREATE INDEX idx_magic_link_pubkey ON magic_link_sessions(pubkey);
CREATE INDEX idx_magic_link_status ON magic_link_sessions(status);

-- Trigger for updated_at
CREATE OR REPLACE TRIGGER update_magic_link_sessions_updated_at
BEFORE UPDATE ON magic_link_sessions
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

-- +goose statementend

-- +goose Down
-- +goose statementbegin

-- Clean up magic_link_sessions without recreating otp_sessions
DROP TRIGGER IF EXISTS update_magic_link_sessions_updated_at ON magic_link_sessions;
DROP TABLE IF EXISTS magic_link_sessions;
DROP TYPE IF EXISTS magic_link_status;

-- +goose statementend
