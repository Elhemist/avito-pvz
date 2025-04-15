CREATE TABLE goods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    type TEXT NOT NULL
);

CREATE INDEX idx_goods_reception_id ON goods(reception_id);

CREATE INDEX idx_goods_added_at ON goods(added_at DESC);
