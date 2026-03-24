-- +goose Up
-- +goose statementbegin
ALTER TABLE otp_sessions ALTER COLUMN challenge SET NOT NULL;
-- +goose statementend


-- +goose Down
-- +goose statementbegin
ALTER TABLE otp_sessions ALTER COLUMN challenge DROP NOT NULL;
-- +goose statementend

--TODO: should migrate 