-- name: CreateCountry :one
INSERT INTO countries (name) VALUES ($1) RETURNING countries.id, countries.name;

-- name: GetCountryByID :one
SELECT countries.id, countries.name FROM countries WHERE countries.id = $1;

-- name: GetCountryByName :one
SELECT countries.id, countries.name FROM countries WHERE countries.name = $1;

-- name: ListCountries :many
SELECT countries.id, countries.name FROM countries;

-- name: UpdateCountry :one
UPDATE countries SET name = $1 WHERE countries.id = $2 RETURNING countries.id, countries.name;

-- name: DeleteCountry :exec
DELETE FROM countries WHERE countries.id = $1;
