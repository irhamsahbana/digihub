CREATE TABLE IF NOT EXISTS clients (
    id CHAR(26) PRIMARY KEY,
    vehicle_type_id CHAR(26) NOT NULL,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(255) NOT NULL,
    vehicle_license_number VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,

    FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types (id),
    UNIQUE (vehicle_license_number)
);