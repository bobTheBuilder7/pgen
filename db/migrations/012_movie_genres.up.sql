CREATE TABLE movie_genres (
	movie_id bigint references movies (id),
	genre_id bigint references genres (id),
	UNIQUE (movie_id, genre_id)
);