CREATE TABLE movie_countries (
	movie_id bigint references movies (id),
	country_id bigint references countries (id),
	UNIQUE (movie_id, country_id)
);