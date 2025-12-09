ALTER TYPE asset_status ADD VALUE 'deleted';
ALTER TABLE assets ADD COLUMN deleted_at TIMESTAMPTZ;
