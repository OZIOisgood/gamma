ALTER TYPE upload_status ADD VALUE 'deleted';
ALTER TABLE uploads ADD COLUMN deleted_at TIMESTAMPTZ;
