-- name: InsertPreKeys :one
WITH signed_key AS (
    INSERT INTO device_signed_prekeys (device_id, key_id, public_key, signature)
    VALUES (
            @deviceId::uuid,
            @signedkeyId::int,
            @signedkey::text,
            @signature::text
        )
    ON CONFLICT (device_id, key_id) DO UPDATE
    SET public_key = EXCLUDED.public_key,
        signature = EXCLUDED.signature
    RETURNING device_id,
        key_id,
        public_key
),
prekeys AS (
    INSERT INTO device_prekeys(device_id, key_id, public_key)
    SELECT s.device_id,
        k.key_id,
        k.public_key
    FROM signed_key as s,
        (
            SELECT unnest(@prekeyIds::int []) as key_id,
                unnest(@prekeys::text []) as public_key
        ) as k
    ON CONFLICT (device_id, key_id) DO UPDATE
    SET public_key = EXCLUDED.public_key
    RETURNING device_id
)
SELECT count(*)
FROM prekeys;


-- name: GetKeybundle :one
WITH user_devices AS (
    SELECT id
    FROM devices
    WHERE user_id = $1
),
signed_keys AS (
    SELECT sp.device_id,
        sp.key_id,
        sp.public_key,
        sp.signature
    FROM device_signed_prekeys sp
    WHERE sp.device_id IN (
            SELECT id
            FROM user_devices
        )
    ORDER BY sp.created_at DESC
    LIMIT 1
),
consumed_prekey AS (
    DELETE FROM device_prekeys
    WHERE (device_id, key_id) = (
            SELECT pk.device_id,
                pk.key_id
            FROM device_prekeys pk
                JOIN signed_keys sk ON sk.device_id = pk.device_id
            ORDER BY pk.created_at ASC,
                pk.key_id ASC
            LIMIT 1
        )
    RETURNING device_id,
        key_id,
        public_key
)
SELECT pk.device_id,
    sk.key_id as signed_key_id,
    sk.public_key as signed_pubkey,
    sk.signature as signed_signature,
    pk.key_id as prekey_id,
    pk.public_key as prekey
FROM consumed_prekey pk
    JOIN signed_keys sk ON sk.device_id = pk.device_id;
