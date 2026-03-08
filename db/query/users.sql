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

-- name: DeleteUserByName :execresult
DELETE FROM users WHERE users.name = $1;

-- name: UpdateUserAge :execresult
UPDATE users SET age = $1 WHERE users.id = $2;

-- name: InsertUser :execresult
INSERT INTO users (name, age) VALUES ($1, $2);

-- name: InsertUserReturning :one
INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name;

-- name: GetUserPosts :many
SELECT u.id, u.name, p.id as post_id, p.name as post_name FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1;
