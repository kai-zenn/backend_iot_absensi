CREATE TABLE attendance_rules (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    effective_from  DATE NOT NULL,
    on_time_until   TIME NOT NULL,        -- batas PRESENT, inklusif (lihat 4a)
    created_by      UUID REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uidx_attendance_rules_effective_from ON attendance_rules (effective_from);

-- Migrasi awal WAJIB menyisipkan satu baris seed (effective_from = '1970-01-01',
-- on_time_until = '06:30:00') agar query "cari rule berlaku" (4a) SELALU menemukan hasil.
