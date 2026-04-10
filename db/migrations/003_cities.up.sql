CREATE TABLE cities (
    id serial PRIMARY KEY,
    name text not null,
    state_id smallint references states (id),
    country_id smallint not null references countries (id)
);