ALTER TABLE assets DROP COLUMN deleted_at;
-- Removing value from ENUM is not directly supported in Postgres without recreating the type.
-- For down migration, we might just leave the enum value or do the complex dance of recreating it.
-- For simplicity in this context, I will just drop the column.
