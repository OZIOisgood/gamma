CREATE TYPE upload_status AS ENUM ('pending', 'uploaded', 'processing', 'ready', 'failed');

CREATE TABLE uploads (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    s3_key TEXT NOT NULL,
    status upload_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
