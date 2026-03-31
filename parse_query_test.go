package main

import (
	"context"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func parseQueries(t *testing.T, sql string) ([]Query, error) {
	t.Helper()
	return parseFileToQueries(context.Background(), strings.NewReader(sql))
}

func TestParseQuery_SingleQuery(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, `
-- name: GetUser :one
SELECT users.id FROM users WHERE users.id = $1;
`)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 1)
	assert.Equal(t, queries[0].name, "GetUser")
	assert.Equal(t, queries[0].t, "one")
}

func TestParseQuery_MultipleQueries(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, `
-- name: GetUser :one
SELECT users.id FROM users WHERE users.id = $1;

-- name: ListUsers :many
SELECT users.id FROM users;
`)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 2)
	assert.Equal(t, queries[0].name, "GetUser")
	assert.Equal(t, queries[1].name, "ListUsers")
}

func TestParseQuery_AllQueryTypes(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, `
-- name: GetUser :one
SELECT users.id FROM users WHERE users.id = $1;

-- name: ListUsers :many
SELECT users.id FROM users;

-- name: DeleteUser :exec
DELETE FROM users WHERE users.id = $1;

-- name: CreateUser :execresult
INSERT INTO users (name) VALUES ($1);
`)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 4)
	assert.Equal(t, queries[0].t, "one")
	assert.Equal(t, queries[1].t, "many")
	assert.Equal(t, queries[2].t, "exec")
	assert.Equal(t, queries[3].t, "execresult")
}

func TestParseQuery_SQLBodyCaptured(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, `
-- name: GetUser :one
SELECT users.id FROM users WHERE users.id = $1;
`)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 1)
	assert.MatchesRegexp(t, queries[0].sql, `SELECT users\.id FROM users WHERE users\.id = \$1`)
}

func TestParseQuery_MultilineSQL(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, `
-- name: GetUser :one
SELECT users.id
FROM users
WHERE users.id = $1;
`)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 1)
	assert.MatchesRegexp(t, queries[0].sql, `SELECT users\.id`)
	assert.MatchesRegexp(t, queries[0].sql, `FROM users`)
}

func TestParseQuery_InvalidHeaderReturnsError(t *testing.T) {
	t.Parallel()
	_, err := parseQueries(t, `
-- name: GetUser
SELECT users.id FROM users WHERE users.id = $1;
`)
	assert.NotNil(t, err)
}

func TestParseQuery_EmptyInput(t *testing.T) {
	t.Parallel()
	queries, err := parseQueries(t, ``)
	assert.Nil(t, err)
	assert.Equal(t, len(queries), 0)
}
