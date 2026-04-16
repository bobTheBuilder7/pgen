ALTER TABLE movies
    ADD COLUMN tags        text[]    NOT NULL DEFAULT '{}',
    ADD COLUMN scores      integer[],
    ADD COLUMN flags       boolean[] NOT NULL DEFAULT '{}',
    ADD COLUMN related_ids bigint[];
