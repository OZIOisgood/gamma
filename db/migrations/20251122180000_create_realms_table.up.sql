CREATE TYPE realm_status AS ENUM ('active', 'deleted');

CREATE TABLE realms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status realm_status NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

-- Partial unique index: name must be unique among non-deleted realms
CREATE UNIQUE INDEX realms_name_key ON realms (name) WHERE status != 'deleted';

INSERT INTO realms (name) VALUES ('default');
