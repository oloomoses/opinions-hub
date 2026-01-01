ALTER TABLE opinions
    ADD CONSTRAINT fk_opinions_parent
        FOREIGN KEY (parent_id)
        REFERENCES opinions(id)
        ON DELETE SET NULL;