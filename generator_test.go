package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func generateQuery(t *testing.T, c *cli, name, queryType, sql string) error {
	t.Helper()
	return c.generateCode([]Query{{name: name, t: queryType, sql: sql}}, io.Discard)
}

func generateQueryOutput(t *testing.T, c *cli, name, queryType, sql string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	err := c.generateCode([]Query{{name: name, t: queryType, sql: sql}}, &buf)
	return buf.String(), err
}

// --- std mode method names ---

func TestGenerateCode_StdModeSelectOneUsesQueryRowContext(t *testing.T) {
	c := testCliWithUsersSchema(t)
	c.std = true
	out, err := generateQueryOutput(t, c, "GetUser", "one", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRowContext`)
}

func TestGenerateCode_StdModeSelectManyUsesQueryContext(t *testing.T) {
	c := testCliWithUsersSchema(t)
	c.std = true
	out, err := generateQueryOutput(t, c, "ListUsers", "many", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryContext`)
}

func TestGenerateCode_StdModeExecUsesExecContext(t *testing.T) {
	c := testCliWithUsersSchema(t)
	c.std = true
	out, err := generateQueryOutput(t, c, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `ExecContext`)
}

func TestGenerateCode_StdModeExecResultReturnsSQLResult(t *testing.T) {
	c := testCliWithUsersSchema(t)
	c.std = true
	out, err := generateQueryOutput(t, c, "DeleteUser", "execresult", `DELETE FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `sql\.Result`)
}

func TestGenerateCode_DefaultModeUsesQueryRow(t *testing.T) {
	c := testCliWithUsersSchema(t)
	out, err := generateQueryOutput(t, c, "GetUser", "one", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`) // QueryRow but not QueryRowContext
}

func TestGenerateCode_DefaultModeExecResultReturnsPgconnCommandTag(t *testing.T) {
	c := testCliWithUsersSchema(t)
	out, err := generateQueryOutput(t, c, "DeleteUser", "execresult", `DELETE FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `pgconn\.CommandTag`)
}

// --- UPDATE without WHERE ---

func TestGenerateCode_UpdateWithoutWhereReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateAll", "exec", `UPDATE users SET name = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateUserName", "exec", `UPDATE users SET name = $1 WHERE users.id = $2;`)
	assert.Nil(t, err)
}

func TestGenerateCode_UpdateMultipleSetWithoutWhereReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateAll", "exec", `UPDATE users SET name = $1, email = $2;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateNamedParamWithoutWhereReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateAll", "exec", `UPDATE users SET name = @name;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereExecResultReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateAll", "execresult", `UPDATE users SET name = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereErrorMentionsQueryName(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "BulkUpdateUsers", "exec", `UPDATE users SET name = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `BulkUpdateUsers`)
}

// --- DELETE without WHERE ---

func TestGenerateCode_DeleteWithoutWhereReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteAll", "exec", `DELETE FROM users;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteWithoutWhereExecResultReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteAll", "execresult", `DELETE FROM users;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithoutWhereErrorMentionsQueryName(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "WipeUsers", "exec", `DELETE FROM users;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WipeUsers`)
}

func TestGenerateCode_DeleteNamedParamWithWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = @user_id;`)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteMultipleWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1 AND users.name = $2;`)
	assert.Nil(t, err)
}

// --- SELECT is unaffected ---

func TestGenerateCode_SelectWithoutWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "ListUsers", "many", `SELECT users.id, users.name FROM users;`)
	assert.Nil(t, err)
}

// --- Unknown query type ---

func TestGenerateCode_UnknownQueryTypeReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "banana", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `banana`)
}

func TestGenerateCode_UnknownQueryTypeErrorMentionsQueryName(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "banana", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetUser`)
}

func TestGenerateCode_UnknownQueryTypeOnInsertReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "CreateUser", "oops", `INSERT INTO users (name) VALUES ($1);`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `oops`)
}

func TestGenerateCode_UnknownQueryTypeOnUpdateReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "UpdateUser", "wrong", `UPDATE users SET name = $1 WHERE users.id = $2;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `wrong`)
}

func TestGenerateCode_UnknownQueryTypeOnDeleteReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "DeleteUser", "nope", `DELETE FROM users WHERE users.id = $1;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `nope`)
}

func TestGenerateCode_EmptyQueryTypeReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "", `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.NotNil(t, err)
}

// --- DISTINCT ---

func TestGenerateCode_DistinctWithoutWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "ListUsers", "many", `SELECT DISTINCT users.id, users.name FROM users;`)
	assert.Nil(t, err)
}

func TestGenerateCode_DistinctWithWhereSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "many", `SELECT DISTINCT users.id, users.name FROM users WHERE users.id = $1;`)
	assert.Nil(t, err)
}

// --- LIMIT / OFFSET params ---

func TestGenerateCode_LimitParamSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "ListUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1;`)
	assert.Nil(t, err)
}

func TestGenerateCode_LimitAndOffsetParamsSucceed(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "ListUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $2;`)
	assert.Nil(t, err)
}

func TestGenerateCode_WhereWithLimitParamSucceeds(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "ListUsers", "many", `SELECT users.id, users.name FROM users WHERE users.name = $1 LIMIT $2;`)
	assert.Nil(t, err)
}

// --- Non-sequential parameters ---

func TestGenerateCode_NonSequentialParamsReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	// $1 and $31 — missing $2 through $30
	err := generateQuery(t, c, "GetFirstNUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $31;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialWhereParamsReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "one", `SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $3;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialErrorMentionsQueryName(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetFirstNUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $31;`)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetFirstNUsers`)
}

func TestGenerateCode_SequentialParamsSucceed(t *testing.T) {
	c := testCliWithUsersSchema(t)
	err := generateQuery(t, c, "GetUser", "one", `SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $2;`)
	assert.Nil(t, err)
}

func TestGenerateCode_ValidQueryTypesSucceed(t *testing.T) {
	validTypes := []struct {
		t   string
		sql string
	}{
		{"one", `SELECT users.id, users.name FROM users WHERE users.id = $1;`},
		{"many", `SELECT users.id, users.name FROM users WHERE users.id = $1;`},
		{"exec", `DELETE FROM users WHERE users.id = $1;`},
		{"execresult", `DELETE FROM users WHERE users.id = $1;`},
	}
	for _, tc := range validTypes {
		c := testCliWithUsersSchema(t)
		err := generateQuery(t, c, "Q", tc.t, tc.sql)
		assert.Nil(t, err)
	}
}
