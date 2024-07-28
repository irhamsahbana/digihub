DROP TABLE IF EXISTS areas;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'area_type') THEN
        DROP TYPE area_type;
    END IF;
END $$;
