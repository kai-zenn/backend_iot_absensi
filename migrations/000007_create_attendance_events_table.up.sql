CREATE TABLE attendance_events (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id             VARCHAR(64) NOT NULL,
    device_id            UUID NOT NULL REFERENCES devices(id),
    card_uid             VARCHAR(40) NOT NULL,
    occurred_at          TIMESTAMPTZ NOT NULL,
    received_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    source               VARCHAR(10) NOT NULL,     -- LIVE | SYNC
    result               VARCHAR(20) NOT NULL DEFAULT 'PROCESSING',
                         -- CHECKED_IN | ALREADY_CHECKED_IN | UNKNOWN_CARD | CARD_BLOCKED |
                         -- STUDENT_INACTIVE | REJECTED_CLOCK
                         -- 'PROCESSING' bersifat TRANSIENT: baris ini di-UPDATE sebelum
                         -- COMMIT di transaksi yang sama (lihat 4a); reader lain tidak
                         -- pernah melihat nilai ini (MVCC snapshot isolation Postgres).
    student_id           UUID REFERENCES students(id),
    attendance_id        UUID REFERENCES attendance(id),
    clock_skew_seconds   INT,                       -- hanya diisi untuk source=LIVE
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- INDEX PALING KRITIS: kunci idempotency. GORM: Clauses(clause.OnConflict{
-- Columns: []clause.Column{{Name:"device_id"},{Name:"event_id"}}, DoNothing:true})
CREATE UNIQUE INDEX uidx_attendance_events_device_event
    ON attendance_events (device_id, event_id);

CREATE INDEX idx_attendance_events_card_uid ON attendance_events (card_uid);
CREATE INDEX idx_attendance_events_occurred_at ON attendance_events (occurred_at);

CREATE INDEX idx_attendance_events_needs_review ON attendance_events (result, created_at DESC)
    WHERE result IN ('UNKNOWN_CARD', 'CARD_BLOCKED', 'STUDENT_INACTIVE', 'REJECTED_CLOCK');
