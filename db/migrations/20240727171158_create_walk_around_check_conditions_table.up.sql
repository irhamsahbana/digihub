CREATE TABLE IF NOT EXISTS walk_around_check_conditions (
    id CHAR(26) PRIMARY KEY,
    walk_around_check_id CHAR(26) NOT NULL,
    area_id CHAR(26) NOT NULL,
    potency_id CHAR(26) NOT NULL,
    is_interested BOOLEAN DEFAULT FALSE NOT NULL,
    path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (walk_around_check_id) REFERENCES walk_around_checks (id),
    FOREIGN KEY (area_id) REFERENCES areas (id),
    FOREIGN KEY (potency_id) REFERENCES potencies (id)
);