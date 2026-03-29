package main

import (
	"context"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/pgen/assert"
	"github.com/valkdb/postgresparser"
)

func TestParseSchema_SingleCreateTable(t *testing.T) {
	c := &cli{}
	sql := `CREATE TABLE employees (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL);`
	err := c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.Nil(t, err)
	cols, ok := c.tablesCol.Load("employees")
	assert.True(t, ok)
	assert.Equal(t, len(cols), 2)
	assert.Equal(t, cols[0].Name, "id")
	assert.Equal(t, cols[1].Name, "name")
}

func TestParseSchema_MultipleStatementsOnlyStoresCreateTable(t *testing.T) {
	c := &cli{}
	sql := `
CREATE TABLE employees (id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL);
CREATE INDEX idx_employees_name ON employees (name);
INSERT INTO employees (name) VALUES ('seed');
`
	err := c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.Nil(t, err)
	cols, ok := c.tablesCol.Load("employees")
	assert.True(t, ok)
	assert.Equal(t, len(cols), 2)
}

func TestParseSchema_NonDDLOnlyFileStoresNothing(t *testing.T) {
	c := &cli{}
	sql := `INSERT INTO employees (name) VALUES ('seed');`
	err := c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.Nil(t, err)
	_, ok := c.tablesCol.Load("employees")
	assert.False(t, ok)
}

func TestParseSchema_DropColumnRemovesColumn(t *testing.T) {
	c := &cli{}
	// First store a table with two columns
	parsed, err := postgresparser.ParseSQLStrict(`CREATE TABLE employees (id BIGSERIAL PRIMARY KEY, salary INTEGER);`)
	assert.Nil(t, err)
	c.tablesCol.Store("employees", parsed.DDLActions[0].ColumnDetails)

	// Now parse a migration that drops the salary column
	err = c.parseSchema(context.Background(), strings.NewReader(`ALTER TABLE employees DROP COLUMN salary;`))
	assert.Nil(t, err)

	cols, ok := c.tablesCol.Load("employees")
	assert.True(t, ok)
	assert.Equal(t, len(cols), 1)
	assert.Equal(t, cols[0].Name, "id")
}

func TestParseSchema_DropColumnIfExistsRemovesColumn(t *testing.T) {
	c := &cli{}
	parsed, err := postgresparser.ParseSQLStrict(`CREATE TABLE employees (id BIGSERIAL PRIMARY KEY, salary INTEGER);`)
	assert.Nil(t, err)
	c.tablesCol.Store("employees", parsed.DDLActions[0].ColumnDetails)

	err = c.parseSchema(context.Background(), strings.NewReader(`ALTER TABLE employees DROP COLUMN IF EXISTS salary;`))
	assert.Nil(t, err)

	cols, _ := c.tablesCol.Load("employees")
	assert.Equal(t, len(cols), 1)
	assert.Equal(t, cols[0].Name, "id")
}

func TestParseSchema_DuplicateCreateTableReturnsError(t *testing.T) {
	c := &cli{}
	sql := `CREATE TABLE employees (id BIGSERIAL PRIMARY KEY);`
	err := c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.Nil(t, err)
	err = c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `employees`)
}

func TestParseSchema_CreateTableThenDropColumnInSameFile(t *testing.T) {
	c := &cli{}
	sql := `
CREATE TABLE employees (id BIGSERIAL PRIMARY KEY, salary INTEGER);
ALTER TABLE employees DROP COLUMN salary;
`
	err := c.parseSchema(context.Background(), strings.NewReader(sql))
	assert.Nil(t, err)
	cols, ok := c.tablesCol.Load("employees")
	assert.True(t, ok)
	assert.Equal(t, len(cols), 1)
	assert.Equal(t, cols[0].Name, "id")
}

func TestParseSchema_DropColumnOnUnknownTableErrors(t *testing.T) {
	c := &cli{}
	err := c.parseSchema(context.Background(), strings.NewReader(`ALTER TABLE employees DROP COLUMN salary;`))
	assert.NotNil(t, err)
}
