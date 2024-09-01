CREATE TABLE IF NOT EXISTS wac_follow_up_logs (
    id CHAR(26) PRIMARY KEY,
    walk_around_check_id CHAR(26) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (walk_around_check_id) REFERENCES walk_around_checks(id)
);