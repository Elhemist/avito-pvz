CREATE TABLE pvz (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    city TEXT NOT NULL,
    registration_date TIMESTAMP NOT NULL DEFAULT NOW()
);
