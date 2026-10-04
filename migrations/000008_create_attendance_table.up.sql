CREATE TABLE attendance (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id            UUID NOT NULL REFERENCES students(id),
    attendance_date       DATE NOT NULL,
    check_in_at           TIMESTAMPTZ NOT NULL,
    check_in_status       VARCHAR(10) NOT NULL,    -- PRESENT | LATE
    check_in_source       VARCHAR(10) NOT NULL,    -- CARD | MANUAL
    check_in_event_id     UUID REFERENCES attendance_events(id),
    check_in_by_user_id   UUID REFERENCES users(id),
    note                  TEXT,
    check_out_at          TIMESTAMPTZ,             -- reserved (fase check-out)
    check_out_status      VARCHAR(10),             -- reserved
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- INDEX KEDUA PALING KRITIS: kunci pencegahan duplicate attendance. GORM:
-- Clauses(clause.OnConflict{Columns: []clause.Column{{Name:"student_id"},
-- {Name:"attendance_date"}}, DoNothing:true})
CREATE UNIQUE INDEX uidx_attendance_student_date ON attendance (student_id, attendance_date);

CREATE INDEX idx_attendance_date ON attendance (attendance_date);
CREATE INDEX idx_attendance_date_status ON attendance (attendance_date, check_in_status);
