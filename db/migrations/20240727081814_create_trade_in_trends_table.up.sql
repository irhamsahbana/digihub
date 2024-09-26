CREATE TABLE IF NOT EXISTS trade_in_trends (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brand VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    year INT NOT NULL,
    min_purchase INT NOT NULL,
    max_purchase INT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    unique (brand, model, type, year)
);