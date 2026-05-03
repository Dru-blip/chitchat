-- +goose Up
-- +goose statementbegin
DROP TABLE IF EXISTS sessions;
DROP TRIGGER IF EXISTS update_sessions_updated_at ON sessions;
-- +goose statementend
