-- name: GetUserByID :users :one
SELECT name
FROM users
WHERE users.id = $1;
