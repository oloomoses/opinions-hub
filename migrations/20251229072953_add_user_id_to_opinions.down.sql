ALTER TABLE opinions
    DROP CONSTRAINT IF EXISTS fk_opinions_user;

DROP INDEX IF EXISTS idx_opinions_user_id;

ALTER TABLE opinions
    DROP COLUMN IF EXISTS user_id;