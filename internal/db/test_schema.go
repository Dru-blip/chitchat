package db

import (
	"context"
	"database/sql"
	"fmt"
)

func SetupTestSchema(db *sql.DB) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS pg_uuidv7`); err != nil {
		return fmt.Errorf("create extension: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE OR REPLACE FUNCTION update_timestamp() 
        RETURNS TRIGGER AS $$ 
        BEGIN 
            NEW.updated_at = CLOCK_TIMESTAMP();
            RETURN NEW;
        END;
        $$ LANGUAGE PLPGSQL
    `); err != nil {
		return fmt.Errorf("create update_timestamp function: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        DO $$ BEGIN
            CREATE TYPE client_type AS ENUM ('mobile', 'web', 'desktop');
        EXCEPTION WHEN duplicate_object THEN NULL;
        END $$;
        
        DO $$ BEGIN
            CREATE TYPE magic_link_status AS ENUM ('pending', 'used', 'expired', 'revoked');
        EXCEPTION WHEN duplicate_object THEN NULL;
        END $$;
    `); err != nil {
		return fmt.Errorf("create enums: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
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
        )
    `); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS magic_link_sessions (
            id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
            token TEXT UNIQUE NOT NULL,
            email VARCHAR(360) NOT NULL,
            pubkey TEXT NOT NULL,
            ip_address INET NOT NULL,
            user_agent TEXT,
            status magic_link_status NOT NULL DEFAULT 'pending',
            expires_at TIMESTAMPTZ NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
        )
    `); err != nil {
		return fmt.Errorf("create magic_link_sessions table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS devices (
            id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
            pubkey TEXT UNIQUE NOT NULL,
            name VARCHAR(255) NOT NULL,
            os TEXT NOT NULL,
            client client_type NOT NULL,
            user_agent TEXT,
            user_id UUID NOT NULL,
            last_seen TIMESTAMPTZ,
            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
            
            FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
        )
    `); err != nil {
		return fmt.Errorf("create devices table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS sessions (
            token TEXT PRIMARY KEY,
            data BYTEA NOT NULL,
            expiry TIMESTAMPTZ NOT NULL
        )
    `); err != nil {
		return fmt.Errorf("create sessions table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS device_signed_prekeys (
            device_id UUID NOT NULL,
            key_id INTEGER NOT NULL,  
            public_key TEXT NOT NULL,
            signature TEXT NOT NULL,  
            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
            
            PRIMARY KEY (device_id, key_id),
            FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE
        )
    `); err != nil {
		return fmt.Errorf("create device_signed_prekeys table: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `
        CREATE TABLE IF NOT EXISTS device_prekeys (
            device_id UUID NOT NULL,
            key_id INTEGER NOT NULL, 
            public_key TEXT NOT NULL,
            created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
            
            PRIMARY KEY (device_id, key_id),
            FOREIGN KEY(device_id) REFERENCES devices(id) ON DELETE CASCADE
        )
    `); err != nil {
		return fmt.Errorf("create device_prekeys table: %w", err)
	}

	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_users_name ON users(name)`,
		`CREATE INDEX IF NOT EXISTS idx_magic_link_token ON magic_link_sessions(token)`,
		`CREATE INDEX IF NOT EXISTS idx_magic_link_email ON magic_link_sessions(email)`,
		`CREATE INDEX IF NOT EXISTS idx_magic_link_expires_at ON magic_link_sessions(expires_at)`,
		`CREATE INDEX IF NOT EXISTS idx_magic_link_pubkey ON magic_link_sessions(pubkey)`,
		`CREATE INDEX IF NOT EXISTS idx_magic_link_status ON magic_link_sessions(status)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_pubkey ON devices(pubkey)`,
		`CREATE INDEX IF NOT EXISTS idx_signed_prekeys_current ON device_signed_prekeys(device_id, created_at DESC)`,
		`CREATE INDEX IF NOT EXISTS idx_prekeys_device ON device_prekeys(device_id, key_id)`,
		`CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry)`,
	}

	for i, idx := range indexes {
		if _, err := tx.ExecContext(ctx, idx); err != nil {
			return fmt.Errorf("create index %d: %w", i, err)
		}
	}

	triggers := []string{
		`CREATE OR REPLACE TRIGGER update_users_updated_at
         BEFORE UPDATE ON users
         FOR EACH ROW EXECUTE PROCEDURE update_timestamp()`,
		`CREATE OR REPLACE TRIGGER update_magic_link_sessions_updated_at
         BEFORE UPDATE ON magic_link_sessions
         FOR EACH ROW EXECUTE PROCEDURE update_timestamp()`,
		`CREATE OR REPLACE TRIGGER update_devices_updated_at
         BEFORE UPDATE ON devices
         FOR EACH ROW EXECUTE PROCEDURE update_timestamp()`,
	}

	for i, trig := range triggers {
		if _, err := tx.ExecContext(ctx, trig); err != nil {
			return fmt.Errorf("create trigger %d: %w", i, err)
		}
	}

	return tx.Commit()
}

func TeardownTestSchema(db *sql.DB) error {
	ctx := context.Background()

	tables := []string{
		"device_prekeys",
		"device_signed_prekeys",
		"sessions",
		"devices",
		"magic_link_sessions",
		"users",
	}

	for _, table := range tables {
		_, _ = db.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
	}

	// Drop types
	_, _ = db.ExecContext(ctx, "DROP TYPE IF EXISTS magic_link_status")
	_, _ = db.ExecContext(ctx, "DROP TYPE IF EXISTS client_type")
	_, _ = db.ExecContext(ctx, "DROP FUNCTION IF EXISTS update_timestamp()")

	return nil
}
