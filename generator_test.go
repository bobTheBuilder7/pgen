package main

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func generateQuery(t *testing.T, c *cli, name, queryType, sql string, std bool) error {
	t.Helper()
	rq, err := c.resolveQuery(query{name: name, t: queryType, sql: sql})
	if err != nil {
		return err
	}
	return c.generateCode(t.Context(), []resolvedQuery{rq}, io.Discard, std)
}

func generateQueryOutput(t *testing.T, c *cli, name, queryType, sql string, std bool) (string, error) {
	t.Helper()
	rq, err := c.resolveQuery(query{name: name, t: queryType, sql: sql})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = c.generateCode(t.Context(), []resolvedQuery{rq}, buf, std)
	return buf.String(), err
}

// --- std mode method names ---
func TestGenerateCode_StdModeSelectOneUsesQueryRowContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one", `SELECT movies.id FROM movies WHERE movies.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRowContext`)
}

func TestGenerateCode_StdModeSelectManyUsesQueryContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "ListMovies", "many", `SELECT movies.id FROM movies WHERE movies.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryContext`)
}

func TestGenerateCode_StdModeExecUsesExecContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteMovie", "exec", `DELETE FROM movies WHERE movies.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `ExecContext`)
}

func TestGenerateCode_StdModeExecResultReturnsSQLResult(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteMovie", "execresult", `DELETE FROM movies WHERE movies.id = $1;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `sql\.Result`)
}

func TestGenerateCode_DefaultModeUsesQueryRow(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one", `SELECT movies.id FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`) // QueryRow but not QueryRowContext
}

func TestGenerateCode_DefaultModeExecResultReturnsPgconnCommandTag(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteMovie", "execresult", `DELETE FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `pgconn\.CommandTag`)
}

// --- UPDATE without WHERE ---

func TestGenerateCode_UpdateWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE movies SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateMovieName", "exec", `UPDATE movies SET name = $1 WHERE movies.id = $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_UpdateMultipleSetWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE movies SET name = $1, box_office = $2;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateNamedParamWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "exec", `UPDATE movies SET name = @name;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereExecResultReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateAll", "execresult", `UPDATE movies SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_UpdateWithoutWhereErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "BulkUpdateMovies", "exec", `UPDATE movies SET name = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `BulkUpdateMovies`)
}

// --- DELETE without WHERE ---

func TestGenerateCode_DeleteWithoutWhereReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteAll", "exec", `DELETE FROM movies;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteMovie", "exec", `DELETE FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteWithoutWhereExecResultReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteAll", "execresult", `DELETE FROM movies;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WHERE`)
}

func TestGenerateCode_DeleteWithoutWhereErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "WipeMovies", "exec", `DELETE FROM movies;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `WipeMovies`)
}

func TestGenerateCode_DeleteNamedParamWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteMovie", "exec", `DELETE FROM movies WHERE movies.id = @movie_id;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DeleteMultipleWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteMovie", "exec", `DELETE FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
}

// --- SELECT is unaffected ---

func TestGenerateCode_SelectWithoutWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListMovies", "many", `SELECT movies.id, movies.name FROM movies;`, false)
	assert.Nil(t, err)
}

// --- Unknown query type ---

func TestGenerateCode_UnknownQueryTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "banana", `SELECT movies.id FROM movies WHERE movies.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `banana`)
}

func TestGenerateCode_UnknownQueryTypeErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "banana", `SELECT movies.id FROM movies WHERE movies.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetMovie`)
}

func TestGenerateCode_UnknownQueryTypeOnInsertReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateMovie", "oops", `INSERT INTO movies (name) VALUES ($1);`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `oops`)
}

func TestGenerateCode_UnknownQueryTypeOnUpdateReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "UpdateMovie", "wrong", `UPDATE movies SET name = $1 WHERE movies.id = $2;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `wrong`)
}

func TestGenerateCode_UnknownQueryTypeOnDeleteReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "DeleteMovie", "nope", `DELETE FROM movies WHERE movies.id = $1;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `nope`)
}

func TestGenerateCode_EmptyQueryTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "", `SELECT movies.id FROM movies WHERE movies.id = $1;`, false)
	assert.NotNil(t, err)
}

// --- DISTINCT ---

func TestGenerateCode_DistinctWithoutWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListMovies", "many", `SELECT DISTINCT movies.id, movies.name FROM movies;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_DistinctWithWhereSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "many", `SELECT DISTINCT movies.id, movies.name FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
}

// --- LIMIT / OFFSET params ---

func TestGenerateCode_LimitParamSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListMovies", "many", `SELECT movies.id, movies.name FROM movies LIMIT $1;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_LimitAndOffsetParamsSucceed(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListMovies", "many", `SELECT movies.id, movies.name FROM movies LIMIT $1 OFFSET $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_WhereWithLimitParamSucceeds(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "ListMovies", "many", `SELECT movies.id, movies.name FROM movies WHERE movies.name = $1 LIMIT $2;`, false)
	assert.Nil(t, err)
}

// --- Non-sequential parameters ---

func TestGenerateCode_NonSequentialParamsReturnsError(t *testing.T) {
	t.Parallel()
	// $1 and $31 — missing $2 through $30
	err := generateQuery(t, testSharedCli, "GetFirstNMovies", "many", `SELECT movies.id, movies.name FROM movies LIMIT $1 OFFSET $31;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialWhereParamsReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "one", `SELECT movies.id, movies.name FROM movies WHERE movies.id = $1 AND movies.name = $3;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `\$2`)
}

func TestGenerateCode_NonSequentialErrorMentionsQueryName(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetFirstNMovies", "many", `SELECT movies.id, movies.name FROM movies LIMIT $1 OFFSET $31;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `GetFirstNMovies`)
}

func TestGenerateCode_SequentialParamsSucceed(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "GetMovie", "one", `SELECT movies.id, movies.name FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
}

func TestGenerateCode_ValidQueryTypesSucceed(t *testing.T) {
	t.Parallel()
	validTypes := []struct {
		t   string
		sql string
	}{
		{"one", `SELECT movies.id, movies.name FROM movies WHERE movies.id = $1;`},
		{"many", `SELECT movies.id, movies.name FROM movies WHERE movies.id = $1;`},
		{"exec", `DELETE FROM movies WHERE movies.id = $1;`},
		{"execresult", `DELETE FROM movies WHERE movies.id = $1;`},
	}
	for _, tc := range validTypes {
		err := generateQuery(t, testSharedCli, "Q", tc.t, tc.sql, false)
		assert.Nil(t, err)
	}
}

// --- RETURNING ---

func TestGenerateCode_InsertReturningOneUsesQueryRow(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`)
}

func TestGenerateCode_InsertReturningManyUsesQuery(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "many",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `\.Query\b`)
}

func TestGenerateCode_StdModeInsertReturningOneUsesQueryRowContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRowContext`)
}

func TestGenerateCode_StdModeInsertReturningManyUsesQueryContext(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "many",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, true)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryContext`)
}

func TestGenerateCode_UpdateReturningOneSucceeds(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "UpdateMovieName", "one",
		`UPDATE movies SET name = $1 WHERE movies.id = $2 RETURNING id, name;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`)
}

func TestGenerateCode_DeleteReturningOneSucceeds(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteMovie", "one",
		`DELETE FROM movies WHERE movies.id = $1 RETURNING id, name;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `QueryRow[^C]`)
}

func TestGenerateCode_InsertReturningGeneratesRowStruct(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `CreateMovieRow`)
}

func TestGenerateCode_InsertReturningGeneratesSQLConst(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `createMovieSQL`)
}

func TestGenerateCode_InsertReturningGeneratesScanCall(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `\.Scan\(`)
}

// --- Params struct (2+ params) ---

func TestGenerateCode_TwoParamsGeneratesParamsStruct(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one",
		`SELECT movies.id FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `GetMovieParams`)
}

func TestGenerateCode_TwoParamsUsesArgInSignature(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one",
		`SELECT movies.id FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `arg GetMovieParams`)
}

func TestGenerateCode_TwoParamsUnpacksFields(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one",
		`SELECT movies.id FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `arg\.Id`)
	assert.MatchesRegexp(t, out, `arg\.Name`)
}

func TestGenerateCode_TwoParamsExecQuery(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "DeleteMovie", "exec",
		`DELETE FROM movies WHERE movies.id = $1 AND movies.name = $2;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `DeleteMovieParams`)
}

func TestGenerateCode_OneParamNoParamsStruct(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovie", "one",
		`SELECT movies.id FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
	assert.False(t, strings.Contains(out, "Params"))
}

func TestGenerateCode_ZeroParamsNoParamsStruct(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "ListMovies", "many",
		`SELECT movies.id FROM movies;`, false)
	assert.Nil(t, err)
	assert.False(t, strings.Contains(out, "Params"))
}

// --- RETURNING type mismatch errors ---

func TestGenerateCode_InsertReturningWithExecTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateMovie", "exec",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `RETURNING`)
}

func TestGenerateCode_InsertReturningWithExecResultTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateMovie", "execresult",
		`INSERT INTO movies (name) VALUES ($1) RETURNING id;`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `RETURNING`)
}

func TestGenerateCode_InsertWithoutReturningWithOneTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateMovie", "one",
		`INSERT INTO movies (name) VALUES ($1);`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `one`)
}

func TestGenerateCode_InsertWithoutReturningWithManyTypeReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "CreateMovie", "many",
		`INSERT INTO movies (name) VALUES ($1);`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `many`)
}

// --- EXISTS ---

func TestGenerateCode_SelectExistsGeneratesBoolField(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "TrailerExists", "one",
		`SELECT EXISTS(SELECT 1 FROM trailers WHERE trailers.movie_id = $1) AS exists_trailers;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `ExistsTrailers\s+bool`)
}

func TestGenerateCode_SelectExistsWithoutAliasReturnsError(t *testing.T) {
	t.Parallel()
	err := generateQuery(t, testSharedCli, "TrailerExists", "one",
		`SELECT EXISTS(SELECT 1 FROM trailers WHERE trailers.movie_id = $1);`, false)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

// --- Enum types ---

func TestGenerateCode_EnumNotNullFieldInRow(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovieStatus", "one",
		`SELECT movies.status FROM movies WHERE movies.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `Status\s+MovieStatus`)
}

func TestGenerateCode_EnumNullableFieldInRow(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetTrailerType", "one",
		`SELECT trailers.trailer_type FROM trailers WHERE trailers.id = $1;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `TrailerType\s+NullTrailerType`)
}

func TestGenerateCode_EnumParamInParamsStruct(t *testing.T) {
	t.Parallel()
	out, err := generateQueryOutput(t, testSharedCli, "GetMovieByStatusAndId", "one",
		`SELECT movies.name FROM movies WHERE movies.id = $1 AND movies.status = $2;`, false)
	assert.Nil(t, err)
	assert.MatchesRegexp(t, out, `GetMovieByStatusAndIdParams`)
	assert.MatchesRegexp(t, out, `Status\s+MovieStatus`)
}
