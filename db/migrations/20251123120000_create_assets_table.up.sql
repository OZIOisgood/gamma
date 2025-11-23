CREATE TYPE asset_status AS ENUM ('processing', 'ready', 'failed');

CREATE TABLE assets (
    id UUID PRIMARY KEY,
    upload_id UUID NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    hls_root TEXT NOT NULL,
    status asset_status NOT NULL DEFAULT 'processing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
