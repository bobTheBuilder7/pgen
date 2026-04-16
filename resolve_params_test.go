package main

import (
	"testing"

	"github.com/bobTheBuilder7/assert"
	"github.com/valkdb/postgresparser"
)

func TestResolveParams_BigintParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_TextParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.name = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_TextParamOnTrailer(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT trailers.id FROM trailers WHERE trailers.url = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"url"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_SmallintParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT cities.id FROM cities WHERE cities.country_id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"country_id"})
	assert.Equal(t, types, []string{"int16"})
}

func TestResolveParams_IntegerParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT trailers.url FROM trailers WHERE trailers.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int32"})
}

func TestResolveParams_DateParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.when_released = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released"})
	assert.Equal(t, types, []string{"pgtype.Date"})
}

func TestResolveParams_MultipleParams(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.id = $1 AND movies.box_office = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "box_office"})
	assert.Equal(t, types, []string{"int64", "pgtype.Int4"})
}

func TestResolveParams_AliasedTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT t.id, t.url FROM trailers t WHERE t.movie_id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"movie_id"})
	assert.Equal(t, types, []string{"pgtype.Int8"})
}

func TestResolveParams_ThreeParamsMixedTypes(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id FROM movies m WHERE m.box_office = $1 AND m.name = $2 AND m.when_released = $3;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"box_office", "name", "when_released"})
	assert.Equal(t, types, []string{"pgtype.Int4", "string", "pgtype.Date"})
}

func TestResolveParams_FourParamsMixedIntSizes(t *testing.T) {
	t.Parallel()
	// movies.id (bigint NN→int64), cities.country_id (smallint NN→int16),
	// trailers.id (serial NN→int32), cities.state_id (smallint NULL→pgtype.Int2)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.name FROM movies WHERE movies.id = $1 AND movies.name = $2 AND movies.box_office = $3 AND movies.when_released = $4;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "name", "box_office", "when_released"})
	assert.Equal(t, types, []string{"int64", "string", "pgtype.Int4", "pgtype.Date"})
}

func TestResolveParams_DeleteSingleParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_DeleteMultipleParams(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM movies WHERE movies.id = $1 AND movies.name = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "name"})
	assert.Equal(t, types, []string{"int64", "string"})
}

func TestResolveParams_UpdateSetAndWhere(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1 WHERE movies.id = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "id"})
	assert.Equal(t, types, []string{"string", "int64"})
}

func TestResolveParams_UpdateMultipleSetColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1, box_office = $2, when_released = $3 WHERE movies.id = $4;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "box_office", "when_released", "id"})
	assert.Equal(t, types, []string{"string", "pgtype.Int4", "pgtype.Date", "int64"})
}

func TestResolveParams_UpdateMultipleWhereColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1 WHERE movies.id = $2 AND movies.box_office = $3;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "id", "box_office"})
	assert.Equal(t, types, []string{"string", "int64", "pgtype.Int4"})
}

func TestResolveParams_InsertSimple(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, box_office, when_released) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "box_office", "when_released"})
	assert.Equal(t, types, []string{"string", "pgtype.Int4", "pgtype.Date"})
}

func TestResolveParams_InsertWithNullableColumns(t *testing.T) {
	t.Parallel()
	// cities: name (text NN), state_id (smallint NULL), country_id (smallint NN)
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO cities (name, state_id) VALUES ($1, $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "state_id"})
	assert.Equal(t, types, []string{"string", "pgtype.Int2"})
}

func TestResolveParams_InsertAllIntSizes(t *testing.T) {
	t.Parallel()
	// movies.id (bigserial), movies.name (text NN), movies.box_office (int NULL)
	// Use cities for smallint variety: country_id (smallint NN), state_id (smallint NULL)
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO cities (name, state_id, country_id) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "state_id", "country_id"})
	assert.Equal(t, types, []string{"string", "pgtype.Int2", "int16"})
}

func TestResolveParams_InsertSingleColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_InsertNullableDate(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, when_released) VALUES ($1, $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "when_released"})
	assert.Equal(t, types, []string{"string", "pgtype.Date"})
}

func TestResolveParams_InsertNullableIntSizes(t *testing.T) {
	t.Parallel()
	// trailers: url (text NN), movie_id (bigint NULL), when_released (date NULL)
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO trailers (url, movie_id, when_released) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"url", "movie_id", "when_released"})
	assert.Equal(t, types, []string{"string", "pgtype.Int8", "pgtype.Date"})
}

func TestResolveParams_InsertAllColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (id, name, when_released, box_office) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "name", "when_released", "box_office"})
	assert.Equal(t, types, []string{"int64", "string", "pgtype.Date", "pgtype.Int4"})
}

func TestResolveParams_InsertMixedNullability(t *testing.T) {
	t.Parallel()
	// cities: name (text NN), state_id (smallint NULL), country_id (smallint NN)
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO cities (name, state_id, country_id) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "state_id", "country_id"})
	assert.Equal(t, types, []string{"string", "pgtype.Int2", "int16"})
}

// JOIN param tests

func TestResolveParams_JoinSingleParam(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, t.url FROM movies m JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_JoinParamsFromBothTables(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, t.url FROM movies m JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1 AND t.url = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "url"})
	assert.Equal(t, types, []string{"int64", "string"})
}

func TestResolveParams_LeftJoinParamFromJoinedTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, t.url FROM movies m LEFT JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1 AND t.when_released = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "when_released"})
	// when_released is NULL in schema and LEFT JOIN keeps it nullable
	assert.Equal(t, types, []string{"int64", "pgtype.Date"})
}

// Subquery param tests

func TestResolveParams_ExistsSubqueryParamOnOuterQuery(t *testing.T) {
	t.Parallel()
	// EXISTS subquery with param on the outer WHERE clause — resolves correctly
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE EXISTS (SELECT 1 FROM trailers WHERE trailers.movie_id = movies.id) AND movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_WhereInSubqueryParam(t *testing.T) {
	t.Parallel()
	// WHERE IN subquery: $1 is inside the subquery (trailers.url = $1)
	// Should resolve to trailers.url (string), not the outer movies.id (int64)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.url = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"url"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_NotInSubqueryParam(t *testing.T) {
	t.Parallel()
	// NOT IN subquery: $1 is inside the subquery (trailers.when_released = $1)
	// Should resolve to trailers.when_released (pgtype.Date)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.id NOT IN (SELECT trailers.movie_id FROM trailers WHERE trailers.when_released = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released"})
	assert.Equal(t, types, []string{"pgtype.Date"})
}

func TestResolveParams_MixedOuterAndSubqueryParams(t *testing.T) {
	t.Parallel()
	// $1 is on outer WHERE (movies.name), $2 is in subquery (trailers.url)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.name = $1 AND movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.url = $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "url"})
	assert.Equal(t, types, []string{"string", "string"})
}

func TestResolveParams_SubqueryParamWithNullableColumn(t *testing.T) {
	t.Parallel()
	// $1 is inside subquery referencing a nullable column (trailers.when_released)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.when_released = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released"})
	assert.Equal(t, types, []string{"pgtype.Date"})
}

func TestResolveParams_SubqueryParamAndOuterParamReversed(t *testing.T) {
	t.Parallel()
	// $1 is in subquery (trailers.url), $2 is on outer WHERE (movies.name)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.url = $1) AND movies.name = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"url", "name"})
	assert.Equal(t, types, []string{"string", "string"})
}

// Named parameter tests — convertNamedParams

func TestConvertNamedParams_NoParams(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, sql)
	assert.Nil(t, names)
}

func TestConvertNamedParams_PositionalOnly(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.id = $1;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, sql)
	assert.Nil(t, names)
}

func TestConvertNamedParams_SingleNamed(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.id = $1;`)
	assert.Equal(t, names, []string{"id"})
}

func TestConvertNamedParams_MultipleNamed(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.id = @id AND movies.name = @name;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.id = $1 AND movies.name = $2;`)
	assert.Equal(t, names, []string{"id", "name"})
}

func TestConvertNamedParams_Underscore(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.box_office = @box_office;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.box_office = $1;`)
	assert.Equal(t, names, []string{"box_office"})
}

func TestConvertNamedParams_MixedError(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.id = $1 AND movies.name = @name;`
	_, _, err := convertNamedParams(sql)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `cannot mix`)
}

func TestConvertNamedParams_InsertNamed(t *testing.T) {
	t.Parallel()
	sql := `INSERT INTO movies (name, box_office, when_released) VALUES (@name, @box_office, @when_released);`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO movies (name, box_office, when_released) VALUES ($1, $2, $3);`)
	assert.Equal(t, names, []string{"name", "box_office", "when_released"})
}

func TestConvertNamedParams_UpdateNamed(t *testing.T) {
	t.Parallel()
	sql := `UPDATE movies SET name = @name WHERE movies.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `UPDATE movies SET name = $1 WHERE movies.id = $2;`)
	assert.Equal(t, names, []string{"name", "id"})
}

func TestConvertNamedParams_DeleteNamed(t *testing.T) {
	t.Parallel()
	sql := `DELETE FROM movies WHERE movies.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `DELETE FROM movies WHERE movies.id = $1;`)
	assert.Equal(t, names, []string{"id"})
}

func TestConvertNamedParams_ThreeParams(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.box_office = @budget AND movies.name = @title AND movies.when_released = @released;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.box_office = $1 AND movies.name = $2 AND movies.when_released = $3;`)
	assert.Equal(t, names, []string{"budget", "title", "released"})
}

func TestConvertNamedParams_ReturningNamed(t *testing.T) {
	t.Parallel()
	sql := `INSERT INTO movies (name) VALUES (@name) RETURNING id, name;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO movies (name) VALUES ($1) RETURNING id, name;`)
	assert.Equal(t, names, []string{"name"})
}

// End-to-end named param tests — convertNamedParams + ParseSQLStrict + resolveParams

func TestNamedParams_SelectSingleParam(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id, movies.name FROM movies WHERE movies.id = @movie_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	// resolveParams returns column-derived names, override with named params
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
	// namedParams has the user-specified names
	assert.Equal(t, namedParams, []string{"movie_id"})
}

func TestNamedParams_SelectMultipleParams(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.id = @movie_id AND movies.name = @movie_name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"movie_id", "movie_name"})
}

func TestNamedParams_SelectWithAliasedTable(t *testing.T) {
	t.Parallel()
	sql := `SELECT t.id, t.url FROM trailers t WHERE t.movie_id = @ref_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"pgtype.Int8"})
	assert.Equal(t, namedParams, []string{"ref_id"})
}

func TestNamedParams_SelectNullableColumn(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.when_released = @release_date AND movies.box_office = @min_box;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"pgtype.Date", "pgtype.Int4"})
	assert.Equal(t, namedParams, []string{"release_date", "min_box"})
}

func TestNamedParams_InsertParams(t *testing.T) {
	t.Parallel()
	sql := `INSERT INTO movies (name, box_office, when_released) VALUES (@movie_name, @budget, @release_date);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "pgtype.Int4", "pgtype.Date"})
	assert.Equal(t, namedParams, []string{"movie_name", "budget", "release_date"})
}

func TestNamedParams_InsertNullableParams(t *testing.T) {
	t.Parallel()
	sql := `INSERT INTO cities (name, state_id) VALUES (@city_name, @state);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "pgtype.Int2"})
	assert.Equal(t, namedParams, []string{"city_name", "state"})
}

func TestNamedParams_UpdateSetAndWhere(t *testing.T) {
	t.Parallel()
	sql := `UPDATE movies SET name = @new_name WHERE movies.id = @movie_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "int64"})
	assert.Equal(t, namedParams, []string{"new_name", "movie_id"})
}

func TestNamedParams_UpdateMultipleSets(t *testing.T) {
	t.Parallel()
	sql := `UPDATE movies SET name = @new_name, box_office = @new_box, when_released = @release WHERE movies.id = @movie_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "pgtype.Int4", "pgtype.Date", "int64"})
	assert.Equal(t, namedParams, []string{"new_name", "new_box", "release", "movie_id"})
}

func TestNamedParams_DeleteSingleParam(t *testing.T) {
	t.Parallel()
	sql := `DELETE FROM movies WHERE movies.id = @movie_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64"})
	assert.Equal(t, namedParams, []string{"movie_id"})
}

func TestNamedParams_DeleteMultipleParams(t *testing.T) {
	t.Parallel()
	sql := `DELETE FROM movies WHERE movies.id = @movie_id AND movies.name = @movie_name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"movie_id", "movie_name"})
}

func TestNamedParams_AllIntSizes(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.name FROM movies WHERE movies.id = @big AND movies.box_office = @medium AND movies.when_released = @date;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "pgtype.Int4", "pgtype.Date"})
	assert.Equal(t, namedParams, []string{"big", "medium", "date"})
}

func TestNamedParams_JoinParams(t *testing.T) {
	t.Parallel()
	sql := `SELECT m.id, t.url FROM movies m JOIN trailers t ON t.movie_id = m.id WHERE m.id = @movie_id AND t.url = @trailer_url;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"movie_id", "trailer_url"})
}

func TestNamedParams_LeftJoinForcesNullable(t *testing.T) {
	t.Parallel()
	sql := `SELECT m.id, t.url FROM movies m LEFT JOIN trailers t ON t.movie_id = m.id WHERE m.id = @movie_id AND t.when_released = @release;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "pgtype.Date"})
	assert.Equal(t, namedParams, []string{"movie_id", "release"})
}

func TestNamedParams_InsertWithReturning(t *testing.T) {
	t.Parallel()
	sql := `INSERT INTO movies (name) VALUES (@movie_name) RETURNING id, name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string"})
	assert.Equal(t, namedParams, []string{"movie_name"})
}

func TestNamedParams_SubqueryParam(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id, movies.name FROM movies WHERE movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.url = @trailer_url);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string"})
	assert.Equal(t, namedParams, []string{"trailer_url"})
}

// --- Duplicate param tests ---

func TestConvertNamedParams_DuplicateNameMapsToSamePosition(t *testing.T) {
	t.Parallel()
	// @val used twice → both should map to $1 (same slot)
	sql := `SELECT movies.id FROM movies WHERE movies.name = @val AND movies.box_office = @val;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.name = $1 AND movies.box_office = $1;`)
	assert.Equal(t, names, []string{"val"})
}

func TestConvertNamedParams_DuplicateNameInUpdate(t *testing.T) {
	t.Parallel()
	// @id used in two WHERE conditions → same $2 slot
	sql := `UPDATE movies SET name = @name WHERE movies.id = @id AND movies.box_office = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `UPDATE movies SET name = $1 WHERE movies.id = $2 AND movies.box_office = $2;`)
	assert.Equal(t, names, []string{"name", "id"})
}

func TestConvertNamedParams_DuplicateNameInInsert(t *testing.T) {
	t.Parallel()
	// All distinct names — no dedup expected
	sql := `INSERT INTO movies (name, box_office) VALUES (@name, @box_office);`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO movies (name, box_office) VALUES ($1, $2);`)
	assert.Equal(t, names, []string{"name", "box_office"})
}

func TestConvertNamedParams_ThreeDistinctNames(t *testing.T) {
	t.Parallel()
	sql := `SELECT movies.id FROM movies WHERE movies.name = @a AND movies.box_office = @b AND movies.when_released = @c;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT movies.id FROM movies WHERE movies.name = $1 AND movies.box_office = $2 AND movies.when_released = $3;`)
	assert.Equal(t, names, []string{"a", "b", "c"})
}

func TestResolveParams_DuplicatePositionalParam(t *testing.T) {
	t.Parallel()
	// $1 used for two different columns — should generate only one function param
	const sql = `SELECT movies.id FROM movies WHERE movies.name = $1 AND movies.box_office = $1;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestNamedParams_DuplicateNameSelectOneFunctionParam(t *testing.T) {
	t.Parallel()
	// @val used twice in WHERE → one function param, SQL uses $1 twice
	const sql = `SELECT movies.id FROM movies WHERE movies.name = @val AND movies.box_office = @val;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"val"})
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	// namedParams overrides resolved names
	names = namedParams
	assert.Equal(t, names, []string{"val"})
	assert.Equal(t, types, []string{"string"})
}

func TestNamedParams_DuplicateNameUpdateOneFunctionParam(t *testing.T) {
	t.Parallel()
	// @id used twice in WHERE → two $N slots collapse to one function param after dedup
	const sql = `UPDATE movies SET name = @new_name WHERE movies.id = @id AND movies.box_office = @id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"new_name", "id"})
	assert.Equal(t, converted, `UPDATE movies SET name = $1 WHERE movies.id = $2 AND movies.box_office = $2;`)
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	// Override with named param names
	names = namedParams
	assert.Equal(t, names, []string{"new_name", "id"})
	assert.Equal(t, types, []string{"string", "int64"})
}

// --- LIMIT / OFFSET params ---

func TestResolveParams_LimitOnly(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies LIMIT $1;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"limit"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_LimitAndOffset(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies LIMIT $1 OFFSET $2;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"limit", "offset"})
	assert.Equal(t, types, []string{"int64", "int64"})
}

func TestResolveParams_WhereAndLimit(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies WHERE movies.name = $1 LIMIT $2;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "limit"})
	assert.Equal(t, types, []string{"string", "int64"})
}

func TestResolveParams_WhereAndLimitAndOffset(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies WHERE movies.name = $1 LIMIT $2 OFFSET $3;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "limit", "offset"})
	assert.Equal(t, types, []string{"string", "int64", "int64"})
}

func TestResolveParams_NamedLimitParam(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies WHERE movies.name = @name LIMIT @lim;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	_, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"name", "lim"})
	assert.Equal(t, types, []string{"string", "int64"})
}

func TestResolveParams_NamedLimitAndOffsetParams(t *testing.T) {
	t.Parallel()
	const sql = `SELECT movies.id FROM movies LIMIT @lim OFFSET @off;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	_, types, err := testSharedCli.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"lim", "off"})
	assert.Equal(t, types, []string{"int64", "int64"})
}

// CTE param tests

func TestResolveParams_CTEFilterParam(t *testing.T) {
	t.Parallel()
	// $1 is inside the CTE's WHERE clause — not in the outer query's ColumnUsage
	parsedSQL, err := postgresparser.ParseSQLStrict(`
WITH recent_movies AS (
    SELECT movies.id, movies.name
    FROM movies
    WHERE movies.when_released = $1
)
SELECT recent_movies.id, recent_movies.name
FROM recent_movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released"})
	assert.Equal(t, types, []string{"pgtype.Date"})
}

func TestResolveParams_CTEMultipleParams(t *testing.T) {
	t.Parallel()
	// $1 and $2 both inside a single CTE's WHERE clause
	parsedSQL, err := postgresparser.ParseSQLStrict(`
WITH filtered AS (
    SELECT movies.id, movies.name
    FROM movies
    WHERE movies.when_released = $1 AND movies.box_office > $2
)
SELECT filtered.id, filtered.name
FROM filtered;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released", "box_office"})
	assert.Equal(t, types, []string{"pgtype.Date", "pgtype.Int4"})
}

func TestResolveParams_CTEFilterAndOuterParam(t *testing.T) {
	t.Parallel()
	// $1 is inside CTE WHERE, $2 is on the outer query
	parsedSQL, err := postgresparser.ParseSQLStrict(`
WITH recent_trailers AS (
    SELECT trailers.id, trailers.url, trailers.movie_id
    FROM trailers
    WHERE trailers.when_released = $1
)
SELECT recent_trailers.url, movies.name
FROM recent_trailers
JOIN movies ON movies.id = recent_trailers.movie_id
WHERE movies.id = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"when_released", "id"})
	assert.Equal(t, types, []string{"pgtype.Date", "int64"})
}

func TestResolveParams_JsonbParam(t *testing.T) {
	t.Parallel()
	// movies.metadata is jsonb NULL → pgtype.JSONB
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.metadata = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"metadata"})
	assert.Equal(t, types, []string{"pgtype.JSONB"})
}

func TestResolveParams_JsonNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.config is json NOT NULL → []byte
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.config = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"config"})
	assert.Equal(t, types, []string{"[]byte"})
}

func TestResolveParams_NumericNullableParam(t *testing.T) {
	t.Parallel()
	// movies.rating is numeric NULL → pgtype.Numeric
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.rating = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"rating"})
	assert.Equal(t, types, []string{"pgtype.Numeric"})
}

func TestResolveParams_NumericNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.price is numeric NOT NULL → pgtype.Numeric
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.price = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"price"})
	assert.Equal(t, types, []string{"pgtype.Numeric"})
}

func TestResolveParams_InsertJsonbAndNumeric(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, metadata, config, rating, price) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "metadata", "config", "rating", "price"})
	assert.Equal(t, types, []string{"string", "pgtype.JSONB", "[]byte", "pgtype.Numeric", "pgtype.Numeric"})
}

func TestResolveParams_ByteaNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.checksum is bytea NOT NULL → []byte
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.checksum = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"checksum"})
	assert.Equal(t, types, []string{"[]byte"})
}

func TestResolveParams_ByteaNullableParam(t *testing.T) {
	t.Parallel()
	// movies.thumbnail is bytea NULL → pgtype.Bytea
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.thumbnail = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"thumbnail"})
	assert.Equal(t, types, []string{"pgtype.Bytea"})
}

func TestResolveParams_InsertBytea(t *testing.T) {
	t.Parallel()
	// checksum is bytea NOT NULL → []byte, thumbnail is bytea NULL → pgtype.Bytea
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, checksum, thumbnail) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "checksum", "thumbnail"})
	assert.Equal(t, types, []string{"string", "[]byte", "pgtype.Bytea"})
}

func TestResolveParams_InsertTextArrayParam(t *testing.T) {
	t.Parallel()
	// movies.tags is text[] NOT NULL → []string
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, tags) VALUES ($1, $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "tags"})
	assert.Equal(t, types, []string{"string", "[]string"})
}

func TestResolveParams_InsertNullableIntArrayParam(t *testing.T) {
	t.Parallel()
	// movies.scores is integer[] NULL → pgtype.Array[int32]
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, scores) VALUES ($1, $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "scores"})
	assert.Equal(t, types, []string{"string", "pgtype.Array[int32]"})
}

func TestResolveParams_WhereTextArrayParam(t *testing.T) {
	t.Parallel()
	// movies.tags is text[] NOT NULL → []string
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.tags = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"tags"})
	assert.Equal(t, types, []string{"[]string"})
}

func TestResolveParams_WhereNullableArrayParam(t *testing.T) {
	t.Parallel()
	// movies.scores is integer[] NULL → pgtype.Array[int32]
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.scores = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"scores"})
	assert.Equal(t, types, []string{"pgtype.Array[int32]"})
}

func TestResolveParams_VarcharNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.slug is varchar(100) NOT NULL → string
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.slug = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"slug"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_VarcharNullableParam(t *testing.T) {
	t.Parallel()
	// movies.description is varchar(500) NULL → pgtype.Text
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.description = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"description"})
	assert.Equal(t, types, []string{"pgtype.Text"})
}

func TestResolveParams_InsertVarchar(t *testing.T) {
	t.Parallel()
	// slug is varchar NOT NULL → string, description is varchar NULL → pgtype.Text
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, slug, description) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "slug", "description"})
	assert.Equal(t, types, []string{"string", "string", "pgtype.Text"})
}

func TestResolveParams_TimeNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.show_time is time NOT NULL → pgtype.Time
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.show_time = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"show_time"})
	assert.Equal(t, types, []string{"pgtype.Time"})
}

func TestResolveParams_TimetzNullableParam(t *testing.T) {
	t.Parallel()
	// movies.show_timetz is timetz NULL → pgtype.Time
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.show_timetz = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"show_timetz"})
	assert.Equal(t, types, []string{"pgtype.Time"})
}

func TestResolveParams_IntervalNotNullParam(t *testing.T) {
	t.Parallel()
	// movies.duration is interval NOT NULL → pgtype.Interval
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.duration = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"duration"})
	assert.Equal(t, types, []string{"pgtype.Interval"})
}

func TestResolveParams_IntervalNullableParam(t *testing.T) {
	t.Parallel()
	// movies.break_time is interval NULL → pgtype.Interval
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.break_time = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"break_time"})
	assert.Equal(t, types, []string{"pgtype.Interval"})
}

func TestResolveParams_InsertTimeInterval(t *testing.T) {
	t.Parallel()
	// show_time NOT NULL → pgtype.Time, show_timetz NULL → pgtype.Time
	// duration NOT NULL → pgtype.Interval, break_time NULL → pgtype.Interval
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name, show_time, show_timetz, duration, break_time) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := testSharedCli.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "show_time", "show_timetz", "duration", "break_time"})
	assert.Equal(t, types, []string{"string", "pgtype.Time", "pgtype.Time", "pgtype.Interval", "pgtype.Interval"})
}
