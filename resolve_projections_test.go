package main

import (
	"testing"

	"github.com/bobTheBuilder7/assert"
	"github.com/bobTheBuilder7/gen"
	"github.com/valkdb/postgresparser"
)

func TestResolveProjections_SingleBigintColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int64", Tag: `json:"id"`}})
	assert.Equal(t, scans, []string{"&i.Id"})
}

func TestResolveProjections_SingleTextColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.name FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Name", Type: "string", Tag: `json:"name"`}})
	assert.Equal(t, scans, []string{"&i.Name"})
}

func TestResolveProjections_TextTrailerColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT trailers.url FROM trailers WHERE trailers.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Url", Type: "string", Tag: `json:"url"`}})
	assert.Equal(t, scans, []string{"&i.Url"})
}

func TestResolveProjections_SmallintColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT cities.country_id FROM cities WHERE cities.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "CountryId", Type: "int16", Tag: `json:"country_id"`}})
	assert.Equal(t, scans, []string{"&i.CountryId"})
}

func TestResolveProjections_IntegerColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT trailers.id FROM trailers WHERE trailers.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int32", Tag: `json:"id"`}})
	assert.Equal(t, scans, []string{"&i.Id"})
}

func TestResolveProjections_DateColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.when_released FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`}})
	assert.Equal(t, scans, []string{"&i.WhenReleased"})
}

func TestResolveProjections_MultipleColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name, movies.when_released FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.WhenReleased"})
}

func TestResolveProjections_AliasedTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, m.name, m.box_office FROM movies m WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.BoxOffice"})
}

func TestResolveProjections_ColumnAlias(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id as movie_id, movies.name as movie_name FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "MovieId", Type: "int64", Tag: `json:"movie_id"`},
		{Name: "MovieName", Type: "string", Tag: `json:"movie_name"`},
	})
	assert.Equal(t, scans, []string{"&i.MovieId", "&i.MovieName"})
}

func TestResolveProjections_AllIntSizes(t *testing.T) {
	t.Parallel()
	// cities: state_id (smallint NULL → pgtype.Int2), country_id (smallint NN → int16), id (serial NN → int32)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT cities.state_id, cities.country_id, cities.id FROM cities WHERE cities.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "StateId", Type: "pgtype.Int2", Tag: `json:"state_id"`},
		{Name: "CountryId", Type: "int16", Tag: `json:"country_id"`},
		{Name: "Id", Type: "int32", Tag: `json:"id"`},
	})
	assert.Equal(t, scans, []string{"&i.StateId", "&i.CountryId", "&i.Id"})
}

func TestResolveProjections_MixedAliasAndNoAlias(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id as movie_id, movies.name, movies.box_office FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "MovieId", Type: "int64", Tag: `json:"movie_id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.MovieId", "&i.Name", "&i.BoxOffice"})
}

func TestResolveProjections_StarSelectReturnsError(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT * FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_StarSelectWithAliasReturnsError(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT * FROM movies m WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_TableDotStarReturnsError(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.* FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_AliasDotStarReturnsError(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.* FROM movies m WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveReturning_InsertSingleColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING id;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int64", Tag: `json:"id"`}})
	assert.Equal(t, scans, []string{"&i.Id"})
}

func TestResolveReturning_InsertMultipleColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING id, name, box_office;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.BoxOffice"})
}

func TestResolveReturning_InsertNullableColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING id, when_released;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.WhenReleased"})
}

func TestResolveReturning_UpdateReturning(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1 WHERE movies.id = $2 RETURNING id, name;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveReturning_DeleteReturning(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM movies WHERE movies.id = $1 RETURNING id, name, when_released;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.WhenReleased"})
}

// Test returning columns that are NOT in the INSERT column list
func TestResolveReturning_InsertReturnsColumnsNotInInsertList(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING id, when_released, box_office;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.WhenReleased", "&i.BoxOffice"})
}

// Test returning only nullable columns
func TestResolveReturning_OnlyNullableColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING when_released, box_office;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.WhenReleased", "&i.BoxOffice"})
}

// Test returning all 4 columns from movies
func TestResolveReturning_InsertReturnsAllColumns(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) RETURNING id, name, when_released, box_office;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.WhenReleased", "&i.BoxOffice"})
}

// Test nullable vs not-null columns in RETURNING
func TestResolveReturning_NullableVsNotNullColumns(t *testing.T) {
	t.Parallel()
	// name is NOT NULL (string), when_released is NULL (pgtype.Date)
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1 WHERE movies.id = $2 RETURNING name, when_released;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
	})
	assert.Equal(t, scans, []string{"&i.Name", "&i.WhenReleased"})
}

// Test same RETURNING columns across all three DML types for consistency
func TestResolveReturning_ConsistentAcrossInsertUpdateDelete(t *testing.T) {
	t.Parallel()

	returningCols := "RETURNING id, name, when_released, box_office;"

	insertSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO movies (name) VALUES ($1) ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse INSERT: %v", err)
	}
	updateSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1 WHERE movies.id = $2 ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse UPDATE: %v", err)
	}
	deleteSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM movies WHERE movies.id = $1 ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse DELETE: %v", err)
	}

	expectedFields := []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	}
	expectedScans := []string{"&i.Id", "&i.Name", "&i.WhenReleased", "&i.BoxOffice"}

	insertFields, insertScans, err := testSharedCli.resolveReturning(insertSQL)
	assert.Nil(t, err)
	assert.Equal(t, insertFields, expectedFields)
	assert.Equal(t, insertScans, expectedScans)

	updateFields, updateScans, err := testSharedCli.resolveReturning(updateSQL)
	assert.Nil(t, err)
	assert.Equal(t, updateFields, expectedFields)
	assert.Equal(t, updateScans, expectedScans)

	deleteFields, deleteScans, err := testSharedCli.resolveReturning(deleteSQL)
	assert.Nil(t, err)
	assert.Equal(t, deleteFields, expectedFields)
	assert.Equal(t, deleteScans, expectedScans)
}

// Test alternating nullable/not-null columns in RETURNING
func TestResolveReturning_AlternatingNullableColumns(t *testing.T) {
	t.Parallel()
	// box_office(nullable int), id(not-null bigint), when_released(nullable date), name(not-null text)
	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM movies WHERE movies.id = $1 RETURNING box_office, id, when_released, name;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.BoxOffice", "&i.Id", "&i.WhenReleased", "&i.Name"})
}

// Test UPDATE returning the column being set + other columns
func TestResolveReturning_UpdateReturnsSetColumnAndOthers(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE movies SET name = $1, box_office = $2 WHERE movies.id = $3 RETURNING id, name, when_released, box_office;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.WhenReleased", "&i.BoxOffice"})
}

// JOIN tests

func TestResolveProjections_InnerJoinColumnsFromBothTables(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, m.name, t.id as trailer_id, t.url FROM movies m JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "TrailerId", Type: "int32", Tag: `json:"trailer_id"`},
		{Name: "Url", Type: "string", Tag: `json:"url"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.TrailerId", "&i.Url"})
}

func TestResolveProjections_LeftJoinForcesNullableOnJoinedTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, m.name, t.id as trailer_id, t.url FROM movies m LEFT JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},                    // movies.id — NOT NULL, base table
		{Name: "Name", Type: "string", Tag: `json:"name"`},               // movies.name — NOT NULL, base table
		{Name: "TrailerId", Type: "pgtype.Int4", Tag: `json:"trailer_id"`}, // trailers.id — NOT NULL in schema but LEFT JOIN makes it nullable
		{Name: "Url", Type: "pgtype.Text", Tag: `json:"url"`},            // trailers.url — NOT NULL in schema but LEFT JOIN makes it nullable
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.TrailerId", "&i.Url"})
}

func TestResolveProjections_LeftJoinNullableColumnStaysNullable(t *testing.T) {
	t.Parallel()

	// trailers.when_released is already nullable in schema, LEFT JOIN should still produce pgtype.Date
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, t.when_released FROM movies m LEFT JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.WhenReleased"})
}

func TestResolveProjections_InnerJoinDoesNotForceNullable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, t.url, t.id as trailer_id FROM movies m INNER JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Url", Type: "string", Tag: `json:"url"`},         // NOT NULL, INNER JOIN doesn't force nullable
		{Name: "TrailerId", Type: "int32", Tag: `json:"trailer_id"`}, // NOT NULL, INNER JOIN doesn't force nullable
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Url", "&i.TrailerId"})
}

func TestResolveProjections_RightJoinForcesNullableOnBaseTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, m.name, t.url FROM movies m RIGHT JOIN trailers t ON t.movie_id = m.id WHERE t.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "pgtype.Int8", Tag: `json:"id"`},     // movies.id — NOT NULL but RIGHT JOIN makes base table nullable
		{Name: "Name", Type: "pgtype.Text", Tag: `json:"name"`}, // movies.name — NOT NULL but RIGHT JOIN makes base table nullable
		{Name: "Url", Type: "string", Tag: `json:"url"`},        // trailers.url — NOT NULL, joined table in RIGHT JOIN keeps types
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Url"})
}

func TestResolveProjections_JoinWithMixedNullability(t *testing.T) {
	t.Parallel()

	// INNER JOIN: nullable columns stay nullable, not-null stay not-null
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT m.id, m.box_office, t.url, t.when_released FROM movies m JOIN trailers t ON t.movie_id = m.id WHERE m.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},                    // movies.id — NOT NULL
		{Name: "BoxOffice", Type: "pgtype.Int4", Tag: `json:"box_office"`}, // movies.box_office — nullable in schema
		{Name: "Url", Type: "string", Tag: `json:"url"`},                 // trailers.url — NOT NULL
		{Name: "WhenReleased", Type: "pgtype.Date", Tag: `json:"when_released"`}, // trailers.when_released — nullable in schema
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.BoxOffice", "&i.Url", "&i.WhenReleased"})
}

// Subquery tests — projections

func TestResolveProjections_WhereInSubqueryColumns(t *testing.T) {
	t.Parallel()
	// WHERE IN subquery: parser only exposes outer table columns in Columns
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE movies.id IN (SELECT trailers.movie_id FROM trailers WHERE trailers.url = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveProjections_ExistsSubqueryColumns(t *testing.T) {
	t.Parallel()
	// EXISTS subquery: parser only exposes outer table columns
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies WHERE EXISTS (SELECT 1 FROM trailers WHERE trailers.movie_id = movies.id) AND movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveProjections_ScalarSubqueryInSelect(t *testing.T) {
	t.Parallel()
	// Scalar subquery in SELECT: the entire subquery becomes a column expression
	// Our code won't find a table.column pattern, so it falls through to "string"
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, (SELECT COUNT(*) FROM trailers WHERE trailers.movie_id = movies.id) as trailer_count FROM movies WHERE movies.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "TrailerCount", Type: "string", Tag: `json:"trailer_count"`}, // scalar subquery falls through to default string
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.TrailerCount"})
}

func TestResolveProjections_KnownColumnsStillResolveAfterFix(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.id, movies.name FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, scanFields, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scanFields, []string{"&i.Id", "&i.Name"})
}

// --- Aggregation functions: error without alias ---

func TestResolveProjections_CountStarWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COUNT(*) FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

func TestResolveProjections_CountColumnWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COUNT(movies.id) FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

func TestResolveProjections_SumWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT SUM(cities.state_id) FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

// --- COUNT: always int64 ---

func TestResolveProjections_CountStarWithAlias(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COUNT(*) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, scanFields, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "int64", Tag: `json:"total"`}})
	assert.Equal(t, scanFields, []string{"&i.Total"})
}

func TestResolveProjections_CountColumnWithAlias(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COUNT(movies.id) AS cnt FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, scanFields, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Cnt", Type: "int64", Tag: `json:"cnt"`}})
	assert.Equal(t, scanFields, []string{"&i.Cnt"})
}

func TestResolveProjections_CountMixedWithRegularColumn(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.name, COUNT(*) AS movie_count FROM movies GROUP BY movies.name;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, scanFields, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "MovieCount", Type: "int64", Tag: `json:"movie_count"`},
	})
	assert.Equal(t, scanFields, []string{"&i.Name", "&i.MovieCount"})
}

// --- SUM: nullable, same base type as column ---

func TestResolveProjections_SumSmallint(t *testing.T) {
	t.Parallel()
	// cities.state_id SMALLINT (nullable) → pgtype.Int2
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT SUM(cities.state_id) AS total FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "pgtype.Int2", Tag: `json:"total"`}})
}

func TestResolveProjections_SumInteger(t *testing.T) {
	t.Parallel()
	// movies.box_office INTEGER NULL → pgtype.Int4 (SUM of nullable int)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT SUM(movies.box_office) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "pgtype.Int4", Tag: `json:"total"`}})
}

func TestResolveProjections_SumBigint(t *testing.T) {
	t.Parallel()
	// movies.id BIGINT NOT NULL → pgtype.Int8 (SUM forces nullable)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT SUM(movies.id) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "pgtype.Int8", Tag: `json:"total"`}})
}

// --- COALESCE: non-nullable inner type ---

func TestResolveProjections_CoalesceOfSumSmallint(t *testing.T) {
	t.Parallel()
	// COALESCE(SUM(cities.state_id), 0) → int16 (non-nullable smallint)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(SUM(cities.state_id), 0) AS total FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "int16", Tag: `json:"total"`}})
}

func TestResolveProjections_CoalesceOfSumInteger(t *testing.T) {
	t.Parallel()
	// COALESCE(SUM(movies.box_office), 0) → int32
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(SUM(movies.box_office), 0) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "int32", Tag: `json:"total"`}})
}

func TestResolveProjections_CoalesceOfSumBigint(t *testing.T) {
	t.Parallel()
	// COALESCE(SUM(movies.id), 0) → int64
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(SUM(movies.id), 0) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "int64", Tag: `json:"total"`}})
}

// --- AVG: always float64 (nullable) ---

func TestResolveProjections_AvgWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT AVG(cities.state_id) FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

func TestResolveProjections_AvgSmallint(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT AVG(cities.state_id) AS avg_state FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "AvgState", Type: "pgtype.Float8", Tag: `json:"avg_state"`}})
}

func TestResolveProjections_AvgBigint(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT AVG(movies.id) AS avg_id FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "AvgId", Type: "pgtype.Float8", Tag: `json:"avg_id"`}})
}

// --- MIN / MAX: nullable, same base type as column ---

func TestResolveProjections_MinWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MIN(cities.state_id) FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

func TestResolveProjections_MaxWithoutAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MAX(cities.state_id) FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `alias`)
}

func TestResolveProjections_MinSmallint(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MIN(cities.state_id) AS min_state FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MinState", Type: "pgtype.Int2", Tag: `json:"min_state"`}})
}

func TestResolveProjections_MaxSmallint(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MAX(cities.state_id) AS max_state FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MaxState", Type: "pgtype.Int2", Tag: `json:"max_state"`}})
}

func TestResolveProjections_MinInteger(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MIN(movies.box_office) AS min_box FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MinBox", Type: "pgtype.Int4", Tag: `json:"min_box"`}})
}

func TestResolveProjections_MaxBigint(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT MAX(movies.id) AS max_id FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MaxId", Type: "pgtype.Int8", Tag: `json:"max_id"`}})
}

func TestResolveProjections_CoalesceOfMin(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(MIN(cities.state_id), 0) AS min_state FROM cities;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MinState", Type: "int16", Tag: `json:"min_state"`}})
}

func TestResolveProjections_CoalesceOfMax(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(MAX(movies.id), 0) AS max_id FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "MaxId", Type: "int64", Tag: `json:"max_id"`}})
}

func TestResolveProjections_CoalesceOfCount(t *testing.T) {
	t.Parallel()
	// COALESCE(COUNT(*), 0) → int64 (COUNT already non-nullable, stays int64)
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT COALESCE(COUNT(*), 0) AS total FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	fields, _, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Total", Type: "int64", Tag: `json:"total"`}})
}

func TestResolveProjections_FromSubqueryErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT sub.id, sub.name FROM (SELECT movies.id, movies.name FROM movies WHERE movies.box_office > $1) sub WHERE sub.id = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `sub`)
}

// --- Unknown column / table errors ---
func TestResolveProjections_UnknownColumnErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.nonexistent FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `nonexistent`)
}

func TestResolveProjections_UnknownColumnErrorMentionsTable(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT movies.typo_col FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `movies`)
}

func TestResolveProjections_UnknownTableAliasErrors(t *testing.T) {
	t.Parallel()
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT x.id FROM movies;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	_, _, err = testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `x\.id`)
}

func TestResolveProjections_CTEColumnAlias(t *testing.T) {
	t.Parallel()
	// Outer query selects columns prefixed with the CTE alias name
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

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveProjections_CTEMixedWithRealTable(t *testing.T) {
	t.Parallel()
	// Outer query mixes columns from a CTE alias and a real joined table
	parsedSQL, err := postgresparser.ParseSQLStrict(`
WITH recent_trailers AS (
    SELECT trailers.id, trailers.url, trailers.movie_id
    FROM trailers
    WHERE trailers.when_released = $1
)
SELECT recent_trailers.url, movies.name
FROM recent_trailers
JOIN movies ON movies.id = recent_trailers.movie_id;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := testSharedCli.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Url", Type: "string", Tag: `json:"url"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Url", "&i.Name"})
}
