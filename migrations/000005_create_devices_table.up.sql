CREATE TABLE devices (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_code       VARCHAR(30) NOT NULL,        -- contoh: "GATE-01"
    name              VARCHAR(100) NOT NULL,
    location          VARCHAR(100),
    role              VARCHAR(20) NOT NULL DEFAULT 'CHECK_IN',  -- CHECK_OUT dicadangkan (lihat 4a)
    status            VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',    -- ACTIVE | DISABLED
    api_key_hash      VARCHAR(64) NOT NULL,        -- SHA-256 hex
    last_seen_at      TIMESTAMPTZ,
    last_sync_at      TIMESTAMPTZ,
    pending_queue     INT NOT NULL DEFAULT 0,      -- nilai TERAKHIR dilaporkan device, bukan real-time
    firmware_version  VARCHAR(20),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uidx_devices_device_code ON devices (device_code);
CREATE UNIQUE INDEX uidx_devices_api_key_hash ON devices (api_key_hash);
