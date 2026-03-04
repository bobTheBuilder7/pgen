CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    user_id bigint not null
);
insert into posts (name, user_id) values ("asdasd", 23)
