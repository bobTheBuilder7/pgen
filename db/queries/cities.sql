-- name: CreateCity :one
INSERT INTO cities (name, state_id, country_id) VALUES ($1, $2, $3) RETURNING cities.id, cities.name, cities.state_id, cities.country_id;

-- name: GetCityByID :one
SELECT cities.id, cities.name, cities.state_id, cities.country_id FROM cities WHERE cities.id = $1;

-- name: ListCities :many
SELECT cities.id, cities.name, cities.state_id, cities.country_id FROM cities;

-- name: ListCitiesByCountry :many
SELECT cities.id, cities.name, cities.state_id, cities.country_id FROM cities WHERE cities.country_id = $1;

-- name: ListCitiesByState :many
SELECT cities.id, cities.name, cities.state_id, cities.country_id FROM cities WHERE cities.state_id = $1;

-- name: UpdateCity :one
UPDATE cities SET name = $1, state_id = $2, country_id = $3 WHERE cities.id = $4 RETURNING cities.id, cities.name, cities.state_id, cities.country_id;

-- name: DeleteCity :exec
DELETE FROM cities WHERE cities.id = $1;
