CREATE TABLE students (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nis         VARCHAR(20) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    class_id    UUID NOT NULL REFERENCES classes(id),
    status      VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',  -- ACTIVE | INACTIVE
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX uidx_students_nis ON students (nis);
CREATE INDEX idx_students_class_id ON students (class_id);
CREATE INDEX idx_students_name ON students (name);
