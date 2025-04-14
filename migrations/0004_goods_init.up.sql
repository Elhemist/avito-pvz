CREATE TABLE goods (
    id SERIAL PRIMARY KEY,
    reception_id INTEGER NOT NULL REFERENCES receptions(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL
);

CREATE INDEX idx_goods_reception_id ON goods(reception_id);

CREATE INDEX idx_goods_added_at ON goods(added_at DESC);
