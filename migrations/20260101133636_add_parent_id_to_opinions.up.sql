ALTER TABLE opinions
    ADD COLUMN parent_id BIGINT;

CREATE INDEX idx_opinions_parent_id ON opinions(parent_id);