-- name: CreateDirector :one
INSERT INTO directors (name, description, dob, city_born) VALUES ($1, $2, $3, $4) RETURNING directors.id, directors.name, directors.description, directors.dob, directors.city_born;

-- name: GetDirectorByID :one
SELECT directors.id, directors.name, directors.description, directors.dob, directors.city_born FROM directors WHERE directors.id = $1;

-- name: ListDirectors :many
SELECT directors.id, directors.name, directors.description, directors.dob, directors.city_born FROM directors;

-- name: UpdateDirector :one
UPDATE directors SET name = $1, description = $2, dob = $3, city_born = $4 WHERE directors.id = $5 RETURNING directors.id, directors.name, directors.description, directors.dob, directors.city_born;

-- name: DeleteDirector :exec
DELETE FROM directors WHERE directors.id = $1;

-- name: GetDirectorWithLocation :one
SELECT directors.id, directors.name, directors.description, directors.dob, cities.name AS city_name, states.name AS state_name, countries.name AS country_name
FROM directors
JOIN cities ON cities.id = directors.city_born
LEFT JOIN states ON states.id = cities.state_id
JOIN countries ON countries.id = cities.country_id
WHERE directors.id = $1;

-- name: ListDirectorsWithLocation :many
SELECT directors.id, directors.name, directors.dob, cities.name AS city_name, states.name AS state_name, countries.name AS country_name
FROM directors
JOIN cities ON cities.id = directors.city_born
LEFT JOIN states ON states.id = cities.state_id
JOIN countries ON countries.id = cities.country_id;

-- name: ListDirectorsByCountry :many
SELECT directors.id, directors.name, directors.dob, cities.name AS city_name, countries.name AS country_name
FROM directors
JOIN cities ON cities.id = directors.city_born
JOIN countries ON countries.id = cities.country_id
WHERE countries.id = $1;

-- name: ListDirectorsByMovie :many
SELECT directors.id, directors.name, directors.description, directors.dob
FROM directors
JOIN movie_directors ON movie_directors.director_id = directors.id
WHERE movie_directors.movie_id = $1;

-- name: ListDirectorsWithMovieCount :many
SELECT directors.id, directors.name, directors.dob
FROM directors
WHERE directors.id IN (SELECT movie_directors.director_id FROM movie_directors);

-- name: ListDirectorsWithNoMovies :many
SELECT directors.id, directors.name, directors.description, directors.dob, directors.city_born
FROM directors
WHERE directors.id NOT IN (SELECT movie_directors.director_id FROM movie_directors);

-- name: ListDirectorsBornAfter :many
SELECT directors.id, directors.name, directors.dob, directors.city_born
FROM directors
WHERE directors.dob > $1;

-- name: ListDirectorsByMovieReleasedAfter :many
SELECT directors.id, directors.name, directors.dob
FROM directors
JOIN movie_directors ON movie_directors.director_id = directors.id
JOIN movies ON movies.id = movie_directors.movie_id
WHERE movies.when_released > $1;

-- name: ListDirectorsOfHighGrossingMovies :many
SELECT directors.id, directors.name, directors.dob
FROM directors
JOIN movie_directors ON movie_directors.director_id = directors.id
JOIN movies ON movies.id = movie_directors.movie_id
WHERE movies.box_office > $1;

-- name: ListDirectorsBornInCountryWithMoviesAfter :many
SELECT directors.id, directors.name, directors.dob, countries.name AS country_name
FROM directors
JOIN cities ON cities.id = directors.city_born
JOIN countries ON countries.id = cities.country_id
JOIN movie_directors ON movie_directors.director_id = directors.id
JOIN movies ON movies.id = movie_directors.movie_id
WHERE countries.id = $1 AND movies.when_released > $2;

-- name: ListDirectorsWithAnyMovie :many
WITH directors_with_movies AS (
    SELECT movie_directors.director_id
    FROM movie_directors
)
SELECT directors.id, directors.name, directors.description, directors.dob, directors.city_born
FROM directors
JOIN directors_with_movies ON directors_with_movies.director_id = directors.id;

-- name: ListDirectorsBornInCity :many
WITH city_directors AS (
    SELECT directors.id, directors.name, directors.dob, directors.city_born
    FROM directors
    WHERE directors.city_born = $1
)
SELECT city_directors.id, city_directors.name, city_directors.dob, cities.name AS city_name, countries.name AS country_name
FROM city_directors
JOIN cities ON cities.id = city_directors.city_born
JOIN countries ON countries.id = cities.country_id;

-- name: ListDirectorsOfRecentMovies :many
WITH recent_movies AS (
    SELECT movies.id
    FROM movies
    WHERE movies.when_released > $1
)
SELECT directors.id, directors.name, directors.dob
FROM directors
JOIN movie_directors ON movie_directors.director_id = directors.id
JOIN recent_movies ON recent_movies.id = movie_directors.movie_id;

-- name: ListDirectorsAndTheirOldestMovie :many
WITH earliest_releases AS (
    SELECT movie_directors.director_id, movies.id, movies.name, movies.when_released
    FROM movie_directors
    JOIN movies ON movies.id = movie_directors.movie_id
    WHERE movies.when_released = (
        SELECT movies.when_released
        FROM movie_directors md2
        JOIN movies ON movies.id = md2.movie_id
        WHERE md2.director_id = movie_directors.director_id
        ORDER BY movies.when_released ASC
        LIMIT 1
    )
)
SELECT directors.id, directors.name, earliest_releases.name AS movie_name, earliest_releases.when_released
FROM directors
JOIN earliest_releases ON earliest_releases.director_id = directors.id;

-- name: ListDirectorsFromCountryWithHighGrossingMovies :many
WITH directors_in_country AS (
    SELECT directors.id, directors.name, directors.dob
    FROM directors
    JOIN cities ON cities.id = directors.city_born
    JOIN countries ON countries.id = cities.country_id
    WHERE countries.id = $1
),
high_grossing_movies AS (
    SELECT movie_directors.director_id
    FROM movie_directors
    JOIN movies ON movies.id = movie_directors.movie_id
    WHERE movies.box_office > $2
)
SELECT directors_in_country.id, directors_in_country.name, directors_in_country.dob
FROM directors_in_country
JOIN high_grossing_movies ON high_grossing_movies.director_id = directors_in_country.id;
