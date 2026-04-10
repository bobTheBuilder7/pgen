CREATE TABLE directors (
	id serial PRIMARY KEY,
	name text not null,
	description text not null default '',
	dob date,
	city_born int references cities (id)
);