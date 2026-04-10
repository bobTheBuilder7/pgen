-- name: CreateState :one
INSERT INTO states (name) VALUES ($1) RETURNING states.id, states.name;

-- name: GetStateByID :one
SELECT states.id, states.name FROM states WHERE states.id = $1;

-- name: GetStateByName :one
SELECT states.id, states.name FROM states WHERE states.name = $1;

-- name: ListStates :many
SELECT states.id, states.name FROM states;

-- name: UpdateState :one
UPDATE states SET name = $1 WHERE states.id = $2 RETURNING states.id, states.name;

-- name: DeleteState :exec
DELETE FROM states WHERE states.id = $1;
