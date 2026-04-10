CREATE TABLE movie_actors (
	movie_id bigint references movies (id),
	actor_id bigint references actors (id),
	UNIQUE (movie_id, actor_id)
);