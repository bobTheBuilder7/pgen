CREATE TABLE movies (
	id bigserial PRIMARY KEY,
	name text not null,
	when_released date,
	box_office integer
);