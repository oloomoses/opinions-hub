CREATE TABLE images (
    id BIGSERIAL PRIMARY KEY,
    opinion_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_images_opinion
        FOREIGN KEY (opinion_id)
        REFERENCES opinions(id)
        ON DELETE CASCADE 
);

CREATE INDEX idx_images_opinion_id ON images(opinion_id);