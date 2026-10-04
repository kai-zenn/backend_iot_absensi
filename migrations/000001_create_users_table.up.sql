CREATE TABLE users (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username       VARCHAR(50) NOT NULL,
    password_hash  VARCHAR(60) NOT NULL,   -- bcrypt, panjang tetap 60 char
    name           VARCHAR(100) NOT NULL,
    role           VARCHAR(10) NOT NULL,    -- ADMIN | STAFF
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uidx_users_username ON users (username);
