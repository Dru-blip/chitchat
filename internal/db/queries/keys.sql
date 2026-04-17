-- name: InsertPreKeys :one
WITH signed_key AS (
    INSERT INTO device_signed_prekeys (device_id, key_id, public_key, signature)
    VALUES (
            @deviceId::uuid,
            @signedkeyId::int,
            @signedkey::text,
            @signature::text
        )
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
    SELECT sp.key_id,
        sp.public_key,
        sp.signature
    FROM device_signed_prekeys sp
        LEFT JOIN user_devices ud ON sp.device_id = ud.id
    LIMIT 1
), consumed_prekey AS (
    DELETE FROM device_prekeys
    WHERE id = (
            SELECT pk.device_id
            FROM device_prekeys pk
            WHERE pk.device_id IN (
                    SELECT id
                    FROM user_devices
                )
            LIMIT 1
        )
    RETURNING key_id,
        public_key
)
SELECT ud.id as device_id,
    sk.key_id as signed_key_id,
    sk.public_key as signed_pubkey,
    pk.key_id as prekey_id,
    pk.public_key as prekey
FROM user_devices as ud
    JOIN signed_keys sk ON sk.device_id = ud.id
    JOIN consumed_prekey pk ON pk.device_id = ud.id;