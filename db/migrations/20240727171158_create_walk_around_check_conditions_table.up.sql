CREATE TABLE IF NOT EXISTS walk_around_check_conditions (
    id CHAR(26) PRIMARY KEY,
    walk_around_check_id CHAR(26) NOT NULL,
    area_id CHAR(26) NOT NULL,
    potency_id CHAR(26) NOT NULL,

    assigned_branch_id CHAR(26),
    assigned_section_id CHAR(26),
    assigned_user_id CHAR(26),

    is_interested BOOLEAN DEFAULT FALSE NOT NULL,
    path TEXT,
    notes VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (walk_around_check_id) REFERENCES walk_around_checks (id),
    FOREIGN KEY (area_id) REFERENCES areas (id),
    FOREIGN KEY (potency_id) REFERENCES potencies (id),
    FOREIGN KEY (assigned_branch_id) REFERENCES branches (id),
    FOREIGN KEY (assigned_section_id) REFERENCES sections (id),
    FOREIGN KEY (assigned_user_id) REFERENCES users (id)
);