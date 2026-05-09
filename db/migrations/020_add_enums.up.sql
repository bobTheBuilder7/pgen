CREATE TYPE trailer_type AS ENUM ('teaser', 'trailer', 'clip');
CREATE TYPE movie_status AS ENUM ('draft', 'released', 'archived');

ALTER TABLE movies
    ADD COLUMN status movie_status NOT NULL DEFAULT 'draft';

ALTER TABLE trailers
    ADD COLUMN trailer_type trailer_type;
