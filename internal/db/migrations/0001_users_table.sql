-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255),
    email VARCHAR(360) UNIQUE NOT NULL,
    image TEXT DEFAULT NULL,
    password TEXT DEFAULT NULL,
    ipkey TEXT DEFAULT NULL UNIQUE,


    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +goose statementbegin
CREATE TRIGGER IF NOT EXISTS update_users_updated_At
BEFORE UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at=CURRENT_TIMESTAMP WHERE id=NEW.id;
END;
-- +goose statementend

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TRIGGER IF EXISTS update_users_updated_At;
