package main

import (
	"testing"

	"github.com/bobTheBuilder7/pgen/assert"
	"github.com/valkdb/postgresparser"
)

const usersSchemaSQL = `CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    age SMALLINT,
    status SMALLINT NOT NULL,
    role_id INTEGER NOT NULL,
    login_count INTEGER,
    org_id BIGINT NOT NULL,
    referrer_id BIGINT,
    active BOOLEAN DEFAULT true,
    verified BOOLEAN NOT NULL DEFAULT false
);`

func testCliWithUsersSchema(t *testing.T) *cli {
	t.Helper()
	c := &cli{}
	parsed, err := postgresparser.ParseSQLStrict(usersSchemaSQL)
	if err != nil {
		t.Fatalf("failed to parse users schema: %v", err)
	}
	c.tablesCol.Store("users", parsed.DDLActions[0].ColumnDetails)
	return c
}

func TestResolveParams_BigintParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_TextParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.name = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_VarcharParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.email = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"email"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_SmallintParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.status = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"status"})
	assert.Equal(t, types, []string{"int16"})
}

func TestResolveParams_IntegerParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.role_id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"role_id"})
	assert.Equal(t, types, []string{"int32"})
}

func TestResolveParams_BooleanParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.verified = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"verified"})
	assert.Equal(t, types, []string{"bool"})
}

func TestResolveParams_MultipleParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.org_id = $1 AND users.active = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"org_id", "active"})
	assert.Equal(t, types, []string{"int64", "bool"})
}

func TestResolveParams_AliasedTable(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.email FROM users u WHERE u.referrer_id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"referrer_id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_ThreeParamsMixedTypes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id FROM users u WHERE u.role_id = $1 AND u.name = $2 AND u.active = $3;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"role_id", "name", "active"})
	assert.Equal(t, types, []string{"int32", "string", "bool"})
}

func TestResolveParams_FourParamsAllIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.email FROM users WHERE users.id = $1 AND users.status = $2 AND users.role_id = $3 AND users.verified = $4;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "status", "role_id", "verified"})
	assert.Equal(t, types, []string{"int64", "int16", "int32", "bool"})
}

func TestResolveParams_DeleteSingleParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_DeleteMultipleParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM users WHERE users.id = $1 AND users.name = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "name"})
	assert.Equal(t, types, []string{"int64", "string"})
}
