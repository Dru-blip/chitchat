-- +goose Up
-- +goose statementbegin
ALTER TABLE messages DROP constraint messages_conversation_id_sequence_id_key;
-- +goose statementend
