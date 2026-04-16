ALTER TABLE movies
    ADD COLUMN checksum bytea NOT NULL DEFAULT '',
    ADD COLUMN thumbnail bytea;
