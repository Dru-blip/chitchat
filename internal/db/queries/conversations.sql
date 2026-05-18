-- name: CreateConversation :one
WITH new_conversation AS(
    INSERT INTO conversations(initiator_id, type)
    VALUES ($1, $2)
    RETURNING *
),
conversation_participants AS(
    INSERT INTO conversation_participants(conversation_id, user_id)
    SELECT conversation_id,user_id
    FROM UNNEST(ARRAY[new_conversation.id,new_conversation.id], ARRAY[new_conversation.initiator_id,$3])
    AS t(conversation_id, user_id)
)
