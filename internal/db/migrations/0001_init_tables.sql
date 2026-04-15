-- +goose Up
-- +goose statementbegin
CREATE OR REPLACE FUNCTION update_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CLOCK_TIMESTAMP();
RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    name VARCHAR(255),
    email VARCHAR(360) UNIQUE NOT NULL,
    image TEXT DEFAULT NULL,
    password TEXT DEFAULT NULL,
    ipkey TEXT NOT NULL UNIQUE,
    onboarding BOOLEAN DEFAULT TRUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE INDEX idx_users_name ON users(name);


CREATE TABLE IF NOT EXISTS otp_sessions (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    email VARCHAR(360) NOT NULL,
    code VARCHAR(6) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    attempts SMALLINT NOT NULL DEFAULT 0 ,
    pubkey TEXT NOT NULL,
    challenge TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_otp_email ON otp_sessions(email);
CREATE INDEX idx_otp_expires_at ON otp_sessions(expires_at);
CREATE INDEX idx_otp_pubkey ON otp_sessions(pubkey);


CREATE TYPE  client_type AS ENUM ('mobile', 'web', 'desktop');

CREATE TABLE IF NOT EXISTS devices(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    pubkey TEXT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    os TEXT NOT NULL,  -- iOS, Android, Windows, macOS, Linux
    client client_type NOT NULL,
    user_agent TEXT,
    user_id UUID NOT NULL,
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_devices_user_id ON devices(user_id);
CREATE INDEX idx_devices_pubkey ON devices(pubkey);

CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY, 
    user_id UUID NOT NULL,
    device_id UUID NOT NULL,
    
    ip_address INET NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    
    is_revoked BOOLEAN DEFAULT FALSE,

    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_device_id ON sessions(device_id);
CREATE INDEX idx_sessions_active ON sessions(user_id, is_revoked);
CREATE INDEX idx_sessions_expires ON sessions(expires_at) WHERE is_revoked = FALSE;


CREATE OR REPLACE TRIGGER  update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE OR REPLACE TRIGGER  update_otp_sessions_updated_at
BEFORE UPDATE ON otp_sessions
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE OR REPLACE TRIGGER  update_devices_updated_at
BEFORE UPDATE ON devices
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

CREATE OR REPLACE TRIGGER  update_sessions_updated_at
BEFORE UPDATE ON sessions
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();
-- +goose statementend


-- +goose Down
-- +goose statementbegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS otp_sessions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_otp_sessions_updated_at ON otp_sessions;
DROP TRIGGER IF EXISTS update_devices_updated_at ON devices;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
DROP FUNCTION IF EXISTS update_timestamp();
-- +goose statementend
