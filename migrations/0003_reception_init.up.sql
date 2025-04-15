CREATE TABLE receptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status TEXT NOT NULL
);

CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);

CREATE UNIQUE INDEX unique_active_reception_per_pvz
ON receptions(pvz_id)
WHERE status = 'active';

CREATE INDEX idx_receptions_pvz_id_created_at ON receptions(pvz_id, created_at);
