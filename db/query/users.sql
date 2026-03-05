-- name: GetUserByID :one
SELECT users.id, users.name
FROM users
WHERE users.id = $1 and users.name = $2;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE users.id = $1;
