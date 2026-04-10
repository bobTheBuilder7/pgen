CREATE TABLE movie_directors (
	movie_id bigint references movies (id),
	director_id bigint references directors (id),
	UNIQUE (movie_id, director_id)
);