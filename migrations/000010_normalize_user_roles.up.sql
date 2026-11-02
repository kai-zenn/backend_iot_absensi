UPDATE users
SET role = lower(role);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'users_role_check'
    ) THEN
        ALTER TABLE users
            ADD CONSTRAINT users_role_check
            CHECK (role IN ('admin', 'staff'));
    END IF;
END $$;
