-- name: GetUser :one
SELECT name FROM users
WHERE id = $1;

-- name: GetAuthors :many
SELECT name FROM users
WHERE id = $1;
