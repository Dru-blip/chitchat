-- name: CreateConversation :one
WITH new_conversation AS(
    INSERT INTO conversations(initiator_id, type)
    VALUES (@initiator_id::uuid, @type::text)
    RETURNING *
),
conversation_participants AS(
    INSERT INTO conversation_participants(conversation_id, user_id)
    SELECT conversation_id, user_id
    FROM unnest(ARRAY[new_conversation.id, new_conversation.id]) AS conversation_id, unnest(ARRAY[@initiator_id::uuid, @participant_id::uuid]) AS user_id
    RETURNING *
)
SELECT new_conversation.*, jsonb_agg(to_jsonb(conversation_participants.*)) AS participants
FROM new_conversation
JOIN conversation_participants ON new_conversation.id = conversation_participants.conversation_id;
