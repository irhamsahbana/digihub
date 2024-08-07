DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'wac_status') THEN
        CREATE TYPE wac_status AS ENUM ('created', 'offered', 'wip', 'completed');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS walk_around_checks (
    id CHAR(26) PRIMARY KEY,
    branch_id CHAR(26) NOT NULL,
    section_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    client_id CHAR(26) NOT NULL,
    invoice_number VARCHAR(255),
    revenue DECIMAL(19, 4) DEFAULT 0.0000 NOT NULL,
    status wac_status DEFAULT 'created' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (branch_id) REFERENCES branches (id),
    FOREIGN KEY (section_id) REFERENCES sections (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (client_id) REFERENCES clients (id),
    UNIQUE (invoice_number)
);