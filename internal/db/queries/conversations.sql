-- name: CreateConversation :one
WITH new_conversation AS(
    INSERT INTO conversations(initiator_id, type)
    VALUES (@initiator_id::uuid, @type::conversation_types)
    RETURNING *
),
conversation_participants AS(
    INSERT INTO conversation_participants(conversation_id, user_id)
    SELECT nc.id,k.user_id
    FROM new_conversation nc,(SELECT unnest(ARRAY[@initiator_id::uuid, (SELECT u.id from users u where u.email=$1 LIMIT 1)]) AS user_id) k
    RETURNING *
)
SELECT nc.*,
    jsonb_agg(
        jsonb_build_object('user_id',cp.user_id,
            'joined_at',cp.joined_at,
            'name',u.name,
            'email',u.email,
            'image',u.image
        )
    ) as participants
FROM new_conversation nc
JOIN conversation_participants cp ON nc.id=cp.conversation_id
LEFT JOIN users u ON cp.user_id = u.id
GROUP BY nc.id,nc.type,nc.name,nc.initiator_id,nc.created_at,nc.updated_at;
