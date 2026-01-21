CREATE TYPE asset_status AS ENUM ('processing', 'ready', 'failed', 'deleted');

CREATE TABLE assets (
    id UUID PRIMARY KEY,
    upload_id UUID NOT NULL REFERENCES uploads(id) ON DELETE CASCADE,
    realm_id UUID NOT NULL REFERENCES realms(id),
    hls_root TEXT NOT NULL,
    thumbnail_root TEXT,
    status asset_status NOT NULL DEFAULT 'processing',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
