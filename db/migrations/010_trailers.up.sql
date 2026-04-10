CREATE TABLE trailers (
	id SERIAL PRIMARY KEY,
	movie_id bigint references movies (id),
	url text not null,
	when_released date
);