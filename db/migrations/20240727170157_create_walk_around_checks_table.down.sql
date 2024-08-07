DROP TABLE IF EXISTS walk_around_checks;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'wac_status') THEN
        DROP TYPE wac_status;
    END IF;
END $$;