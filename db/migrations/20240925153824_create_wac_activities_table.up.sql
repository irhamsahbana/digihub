CREATE TABLE IF NOT EXISTS wac_activities (
    id CHAR(26) PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    status VARCHAR(255) NOT NULL,
    total_potential_leads INT NOT NULL DEFAULT 0,
    total_leads INT NOT NULL DEFAULT 0,
    total_completed_leads INT NOT NULL DEFAULT 0,
    total_revenue INT NOT NULL DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);