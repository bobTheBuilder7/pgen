-- name: GetUserByID :one
SELECT users.id, users.name
FROM users
WHERE users.id = $1 and users.name = $2;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE users.id = $1;

-- name: UpdateUserName :exec
UPDATE users SET name = $1
WHERE users.id = $2;

-- name: CreateUser :exec
INSERT INTO users (name, age) VALUES ($1, $2);

-- name: ListUsers :many
SELECT users.id, users.name FROM users;
