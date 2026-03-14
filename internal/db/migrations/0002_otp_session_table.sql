-- +goose Up
CREATE TABLE IF NOT EXISTS otp_sessions (
    id VARCHAR(40) PRIMARY KEY NOT NULL,
    email VARCHAR(360) UNIQUE NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at DATETIME NOT NULL,
    attempts INTEGER DEFAULT 0,
    pubkey TEXT UNIQUE NOT NULL,
    challenge TEXT,


    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose statementbegin
CREATE TRIGGER IF NOT EXISTS update_otp_sessions_updated_At
BEFORE UPDATE ON otp_sessions
FOR EACH ROW
BEGIN
    UPDATE otp_sessions SET updated_at=CURRENT_TIMESTAMP WHERE id=NEW.id;
END;
-- +goose statementend

-- +goose Down
DROP TABLE IF EXISTS otp_sessions;
DROP TRIGGER IF EXISTS update_otp_sessions_updated_At;
