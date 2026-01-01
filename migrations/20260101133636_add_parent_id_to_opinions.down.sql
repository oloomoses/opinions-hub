DROP INDEX IF EXISTS idx_opinions_parent_id;
ALTER TABLE opinions DROP COLUMN IF EXISTS parent_id;