-- +goose Up
-- +goose statementbegin

DROP TYPE IF EXISTS conversation_types;

CREATE TYPE conversation_types AS ENUM ('group', 'dm');


CREATE TABLE IF NOT EXISTS conversations(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuidv7(),
    type conversation_types NOT NULL,
    name VARCHAR(255),
    initiator_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (initiator_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS conversation_participants(
    conversation_id UUID NOT NULL,
    user_id UUID NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    last_read TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (conversation_id, user_id),
    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    conversation_id UUID NOT NULL,
    sender_user_id UUID NOT NULL,
    sender_device_id UUID NOT NULL,
    sequence_id INT NOT NULL,
    content_type TEXT DEFAULT 'text',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_device_id) REFERENCES devices(id) ON DELETE CASCADE,


    UNIQUE (conversation_id, sequence_id)
);


CREATE INDEX idx_messages_conversation_created
    ON messages (conversation_id, created_at);


CREATE TABLE message_envelopes (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    message_id UUID NOT NULL,
    recipient_user_id UUID NOT NULL,
    recipient_device_id UUID NOT NULL,
    is_incoming BOOLEAN NOT NULL,
    context TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (recipient_device_id) REFERENCES devices(id) ON DELETE CASCADE,


    UNIQUE (message_id, recipient_device_id)
);

CREATE INDEX idx_envelopes_message
    ON message_envelopes (message_id);

CREATE INDEX idx_envelopes_recipient
    ON message_envelopes (recipient_user_id, recipient_device_id);

CREATE OR REPLACE TRIGGER  update_conversations_updated_at
BEFORE UPDATE ON conversations
FOR EACH ROW
EXECUTE PROCEDURE update_timestamp();

-- +goose statementend


-- +goose Down
-- +goose statementbegin

DROP TRIGGER IF EXISTS update_conversations_updated_at ON conversations;
DROP INDEX IF EXISTS idx_envelopes_message;
DROP INDEX IF EXISTS idx_envelopes_recipient;
DROP INDEX IF EXISTS idx_messages_conversation_created;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS message_envelopes;
DROP TABLE IF EXISTS conversation_participants;
DROP TABLE IF EXISTS conversations;
DROP TYPE IF EXISTS conversation_types;


-- +goose statementend
