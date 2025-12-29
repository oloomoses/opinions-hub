ALTER TABLE opinions
    ADD COLUMN user_id BIGINT;
-- ALTER TABLE opinions
--     ALTER COLUMN user_id SET NOT NULL;

CREATE INDEX idx_opinions_user_id ON opinions(user_id);

ALTER TABLE opinions
    ADD CONSTRAINT fk_opinions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE;

