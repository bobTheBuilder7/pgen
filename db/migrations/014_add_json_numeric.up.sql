ALTER TABLE movies
    ADD COLUMN metadata     jsonb,
    ADD COLUMN config       json NOT NULL DEFAULT '{}',
    ADD COLUMN rating       numeric,
    ADD COLUMN price        numeric NOT NULL DEFAULT 0;
