package main

import (
	"context"
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func prepare(t *testing.T, sql string) error {
	t.Helper()
	return testSharedCli.prepareQuery(context.Background(), sql)
}

// --- SELECT ---

func TestPrepareQuery_SelectSingleColumn(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users WHERE users.id = $1`))
}

func TestPrepareQuery_SelectMultipleColumns(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, users.name, users.email FROM users WHERE users.id = $1`))
}

func TestPrepareQuery_SelectAllColumns(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, users.name, users.email, users.age, users.status, users.role_id, users.login_count, users.org_id, users.referrer_id, users.active, users.verified FROM users WHERE users.id = $1`))
}

func TestPrepareQuery_SelectWithMultipleParams(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users WHERE users.name = $1 AND users.age = $2`))
}

func TestPrepareQuery_SelectWithLike(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, users.name FROM users WHERE users.name LIKE $1`))
}

func TestPrepareQuery_SelectWithOrderBy(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, users.name FROM users ORDER BY users.name`))
}

func TestPrepareQuery_SelectWithLimitOffset(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users ORDER BY users.id LIMIT $1 OFFSET $2`))
}

func TestPrepareQuery_SelectWithNullCheck(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users WHERE users.referrer_id IS NULL`))
}

func TestPrepareQuery_SelectWithBooleanFilter(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users WHERE users.active = $1`))
}

// --- aggregations ---

func TestPrepareQuery_SelectCount(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT COUNT(*) FROM users`))
}

func TestPrepareQuery_SelectCountWithFilter(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT COUNT(*) FROM users WHERE users.active = $1`))
}

func TestPrepareQuery_SelectSum(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT SUM(users.age) FROM users`))
}

func TestPrepareQuery_SelectAvg(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT AVG(users.age) FROM users WHERE users.active = $1`))
}

func TestPrepareQuery_SelectMinMax(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT MIN(users.age), MAX(users.age) FROM users`))
}

func TestPrepareQuery_SelectGroupBy(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.status, COUNT(*) FROM users GROUP BY users.status`))
}

func TestPrepareQuery_SelectHaving(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.status, COUNT(*) FROM users GROUP BY users.status HAVING COUNT(*) > $1`))
}

// --- JOINs ---

func TestPrepareQuery_InnerJoin(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, posts.title FROM users JOIN posts ON posts.user_id = users.id WHERE users.id = $1`))
}

func TestPrepareQuery_LeftJoin(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, posts.title FROM users LEFT JOIN posts ON posts.user_id = users.id WHERE users.id = $1`))
}

func TestPrepareQuery_MultipleJoins(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, posts.title FROM users JOIN posts ON posts.user_id = users.id WHERE users.active = $1 AND posts.published = $2`))
}

// --- subqueries ---

func TestPrepareQuery_SubqueryInWhere(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.published = $1)`))
}

func TestPrepareQuery_SubqueryInSelect(t *testing.T) {
	assert.Nil(t, prepare(t, `SELECT users.id, (SELECT COUNT(*) FROM posts WHERE posts.user_id = users.id) FROM users WHERE users.id = $1`))
}

// --- CTE ---

func TestPrepareQuery_CTE(t *testing.T) {
	assert.Nil(t, prepare(t, `WITH active_users AS (SELECT users.id FROM users WHERE users.active = $1) SELECT id FROM active_users`))
}

// --- INSERT ---

func TestPrepareQuery_Insert(t *testing.T) {
	assert.Nil(t, prepare(t, `INSERT INTO users (name, email, status, role_id, org_id) VALUES ($1, $2, $3, $4, $5)`))
}

func TestPrepareQuery_InsertReturning(t *testing.T) {
	assert.Nil(t, prepare(t, `INSERT INTO users (name, email, status, role_id, org_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`))
}

// --- UPDATE ---

func TestPrepareQuery_Update(t *testing.T) {
	assert.Nil(t, prepare(t, `UPDATE users SET name = $1 WHERE users.id = $2`))
}

func TestPrepareQuery_UpdateMultipleColumns(t *testing.T) {
	assert.Nil(t, prepare(t, `UPDATE users SET name = $1, email = $2, age = $3 WHERE users.id = $4`))
}

func TestPrepareQuery_UpdateReturning(t *testing.T) {
	assert.Nil(t, prepare(t, `UPDATE users SET name = $1 WHERE users.id = $2 RETURNING id, name`))
}

// --- DELETE ---

func TestPrepareQuery_Delete(t *testing.T) {
	assert.Nil(t, prepare(t, `DELETE FROM users WHERE users.id = $1`))
}

func TestPrepareQuery_DeleteReturning(t *testing.T) {
	assert.Nil(t, prepare(t, `DELETE FROM users WHERE users.id = $1 RETURNING id`))
}

// --- invalid queries ---

func TestPrepareQuery_NonexistentColumn(t *testing.T) {
	assert.NotNil(t, prepare(t, `SELECT users.nonexistent FROM users WHERE users.id = $1`))
}

func TestPrepareQuery_NonexistentTable(t *testing.T) {
	assert.NotNil(t, prepare(t, `SELECT ghost.id FROM ghost WHERE ghost.id = $1`))
}

func TestPrepareQuery_NonexistentColumnInWhere(t *testing.T) {
	assert.NotNil(t, prepare(t, `SELECT users.id FROM users WHERE users.ghost = $1`))
}

func TestPrepareQuery_NonexistentColumnInJoin(t *testing.T) {
	assert.NotNil(t, prepare(t, `SELECT users.id FROM users JOIN posts ON posts.ghost = users.id`))
}

func TestPrepareQuery_NonexistentColumnInInsert(t *testing.T) {
	assert.NotNil(t, prepare(t, `INSERT INTO users (ghost) VALUES ($1)`))
}

func TestPrepareQuery_NonexistentColumnInUpdate(t *testing.T) {
	assert.NotNil(t, prepare(t, `UPDATE users SET ghost = $1 WHERE users.id = $2`))
}

func TestPrepareQuery_TypeMismatch(t *testing.T) {
	// age is SMALLINT, comparing to a text literal causes a type error
	assert.NotNil(t, prepare(t, `SELECT users.id FROM users WHERE users.age = 'not_a_number'`))
}
