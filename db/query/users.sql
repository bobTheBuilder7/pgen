-- name: GetUserByID :one
SELECT u.id as user_id, u.name as user_name, 'asdsda' as suren
FROM users u
WHERE u.id = $1;
