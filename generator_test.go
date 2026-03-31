package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func generateQuery(t *testing.T, c *cli, name, queryType, sql string, std bool) error {
	t.Helper()
	return c.generateCode(t.Context(), []Query{{name: name, t: queryType, sql: sql}}, io.Discard, std)
}

func generateQueryOutput(t *testing.T, c *cli, name, queryType, sql string, std bool) (string, error) {
	t.Helper()

	buf := new(bytes.Buffer)
	err := c.generateCode(t.Context(), []Query{{name: name, t: queryType, sql: sql}}, buf, std)

	return buf.String(), err
}

// --- std mode method names ---
func TestGenerateCode_StdModeSelectOneUsesQueryRowContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetUser", "one", `SELECT users.id FROM users WHERE users.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRowContext`)
}

func TestGenerateCode_StdModeSelectManyUsesQueryContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "ListUsers", "many", `SELECT users.id FROM users WHERE users.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryContext`)
}

func TestGenerateCode_StdModeExecUsesExecContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `ExecContext`)
}

func TestGenerateCode_StdModeExecResultReturnsSQLResult(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteUser", "execresult", `DELETE FROM users WHERE users.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `sql\.Result`)
}

func TestGenerateCode_DefaultModeUsesQueryRow(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetUser", "one", `SELECT users.id FROM users WHERE users.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`) // QueryRow but not QueryRowContext
}

func TestGenerateCode_DefaultModeExecResultReturnsPgconnCommandTag(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteUser", "execresult", `DELETE FROM users WHERE users.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `pgconn\.CommandTag`)
}

// --- UPDATE without WHERE ---

func TestGenerateCode_UpdateWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE users SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateUserName", "exec", `UPDATE users SET name = $1 WHERE users.id = $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_UpdateMultipleSetWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE users SET name = $1, email = $2;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateNamedParamWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE users SET name = @name;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereExecResultReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "execresult", `UPDATE users SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "BulkUpdateUsers", "exec", `UPDATE users SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `BulkUpdateUsers`)
}

// --- DELETE without WHERE ---

func TestGenerateCode_DeleteWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteAll", "exec", `DELETE FROM users;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteWithoutWhereExecResultReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteAll", "execresult", `DELETE FROM users;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithoutWhereErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "WipeUsers", "exec", `DELETE FROM users;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WipeUsers`)
}

func TestGenerateCode_DeleteNamedParamWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = @user_id;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteMultipleWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteUser", "exec", `DELETE FROM users WHERE users.id = $1 AND users.name = $2;`, false)
	assert.Nil(t, err)
}

// --- SELECT is unaffected ---

func TestGenerateCode_SelectWithoutWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListUsers", "many", `SELECT users.id, users.name FROM users;`, false)
	assert.Nil(t, err)
}

// --- Unknown query type ---

func TestGenerateCode_UnknownQueryTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "banana", `SELECT users.id FROM users WHERE users.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `banana`)
}

func TestGenerateCode_UnknownQueryTypeErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "banana", `SELECT users.id FROM users WHERE users.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetUser`)
}

func TestGenerateCode_UnknownQueryTypeOnInsertReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateUser", "oops", `INSERT INTO users (name) VALUES ($1);`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `oops`)
}

func TestGenerateCode_UnknownQueryTypeOnUpdateReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateUser", "wrong", `UPDATE users SET name = $1 WHERE users.id = $2;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `wrong`)
}

func TestGenerateCode_UnknownQueryTypeOnDeleteReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteUser", "nope", `DELETE FROM users WHERE users.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `nope`)
}

func TestGenerateCode_EmptyQueryTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "", `SELECT users.id FROM users WHERE users.id = $1;`, false)
	assert.NotNil(t, err)
}

// --- DISTINCT ---

func TestGenerateCode_DistinctWithoutWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListUsers", "many", `SELECT DISTINCT users.id, users.name FROM users;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DistinctWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "many", `SELECT DISTINCT users.id, users.name FROM users WHERE users.id = $1;`, false)
	assert.Nil(t, err)
}

// --- LIMIT / OFFSET params ---

func TestGenerateCode_LimitParamSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_LimitAndOffsetParamsSucceed(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_WhereWithLimitParamSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListUsers", "many", `SELECT users.id, users.name FROM users WHERE users.name = $1 LIMIT $2;`, false)
	assert.Nil(t, err)
}

// --- Non-sequential parameters ---

func TestGenerateCode_NonSequentialParamsReturnsError(t *testing.T) {
	t.Parallel()
	// $1 and $31 — missing $2 through $30
	err := generateQuery(t, testSharedCli, "GetFirstNUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $31;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialWhereParamsReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "one", `SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $3;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetFirstNUsers", "many", `SELECT users.id, users.name FROM users LIMIT $1 OFFSET $31;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetFirstNUsers`)
}

func TestGenerateCode_SequentialParamsSucceed(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "one", `SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_NonexistentColumnReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "one", `SELECT users.nonexistent FROM users WHERE users.id = $1;`, false)
	assert.NotNil(t, err)
}

func TestGenerateCode_NonexistentTableReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetUser", "one", `SELECT ghost.id FROM ghost WHERE ghost.id = $1;`, false)
	assert.NotNil(t, err)
}

func TestGenerateCode_ValidQueryTypesSucceed(t *testing.T) {
	t.Parallel()
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
		err := generateQuery(t, testSharedCli, "Q", tc.t, tc.sql, false)
		assert.Nil(t, err)
	}
}
