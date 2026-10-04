CREATE TABLE rfid_cards (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    uid         VARCHAR(40) NOT NULL,       -- UID ter-normalisasi (hex uppercase, tanpa separator)
    student_id  UUID NOT NULL REFERENCES students(id),
    status      VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE | BLOCKED | RETIRED
    note        TEXT,
    issued_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- UID unik SECARA GLOBAL, termasuk kartu BLOCKED/RETIRED — supaya tap kartu lama
-- yang sudah diblokir tetap dikenali sebagai CARD_BLOCKED, bukan UNKNOWN_CARD.
CREATE UNIQUE INDEX uidx_rfid_cards_uid ON rfid_cards (uid);

-- Partial unique: hanya SATU kartu ACTIVE per siswa. INDEX PALING PENTING di
-- tabel ini — mencegah dua kartu aktif tertukar untuk siswa yang sama (lihat 4e).
CREATE UNIQUE INDEX uidx_rfid_cards_student_active
    ON rfid_cards (student_id) WHERE status = 'ACTIVE';

CREATE INDEX idx_rfid_cards_student_id ON rfid_cards (student_id);
