CREATE TABLE IF NOT EXISTS walk_around_checks (
    id CHAR(26) PRIMARY KEY,
    branch_id CHAR(26) NOT NULL,
    section_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    client_id CHAR(26) NOT NULL,
    is_offered BOOLEAN DEFAULT FALSE NOT NULL,
    is_need_follow_up BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (branch_id) REFERENCES branches (id),
    FOREIGN KEY (section_id) REFERENCES sections (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (client_id) REFERENCES clients (id)
);