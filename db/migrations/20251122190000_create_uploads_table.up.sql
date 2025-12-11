CREATE TYPE upload_status AS ENUM ('pending', 'uploaded', 'processing', 'ready', 'failed', 'deleted');

CREATE TABLE uploads (
    id UUID PRIMARY KEY,
    realm_id UUID NOT NULL REFERENCES realms(id),
    title TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    status upload_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
