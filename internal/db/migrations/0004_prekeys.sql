-- +goose Up
-- +goose statementbegin

DROP TABLE IF EXISTS device_prekeys;

CREATE TABLE device_signed_prekeys (
    device_id UUID NOT NULL,
    key_id INTEGER NOT NULL,  
    public_key TEXT NOT NULL,
    signature TEXT NOT NULL,  
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (device_id, key_id),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);


CREATE INDEX idx_signed_prekeys_current 
ON device_signed_prekeys(device_id, created_at DESC);


CREATE TABLE device_prekeys (
    device_id UUID NOT NULL,
    key_id INTEGER NOT NULL, 
    public_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (device_id, key_id),
    FOREIGN KEY (device_id) REFERENCES devices(id) ON DELETE CASCADE
);


CREATE INDEX idx_prekeys_device 
ON device_prekeys(device_id, key_id);

-- +goose statementend
-- +goose Down
-- +goose statementbegin
DROP TABLE IF EXISTS device_prekeys
-- +goose statementend