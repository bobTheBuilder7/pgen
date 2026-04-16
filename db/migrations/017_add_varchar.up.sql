ALTER TABLE movies
    ADD COLUMN slug        varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN description varchar(500);
