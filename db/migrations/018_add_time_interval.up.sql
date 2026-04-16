ALTER TABLE movies
    ADD COLUMN show_time   time     NOT NULL DEFAULT '00:00:00',
    ADD COLUMN show_timetz timetz,
    ADD COLUMN duration    interval NOT NULL DEFAULT '0',
    ADD COLUMN break_time  interval;
