CREATE TABLE pvz (
    id SERIAL PRIMARY KEY,
    city TEXT NOT NULL,
    registration_date TIMESTAMP NOT NULL DEFAULT NOW()
);
