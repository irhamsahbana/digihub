-- create table if not exists
CREATE TABLE IF NOT EXISTS users (
    id CHAR(26) PRIMARY KEY,
    role_id CHAR(26) NOT NULL,
    branch_id CHAR(26),
    section_id CHAR(26),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    whatsapp_number VARCHAR(255) NOT NULL DEFAULT '',
    path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (role_id) REFERENCES roles (id),
    FOREIGN KEY (branch_id) REFERENCES branches (id),
    FOREIGN KEY (section_id) REFERENCES potencies (id)
);