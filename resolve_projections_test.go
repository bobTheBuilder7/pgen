package main

import (
	"testing"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/assert"
	"github.com/valkdb/postgresparser"
)

func TestResolveProjections_SingleBigintColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int64"}})
	assert.Equal(t, scans, []string{"&i.Id"})
}

func TestResolveProjections_SingleTextColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.name FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Name", Type: "string"}})
	assert.Equal(t, scans, []string{"&i.Name"})
}

func TestResolveProjections_VarcharColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.email FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Email", Type: "string"}})
	assert.Equal(t, scans, []string{"&i.Email"})
}

func TestResolveProjections_SmallintColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.status FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Status", Type: "int16"}})
	assert.Equal(t, scans, []string{"&i.Status"})
}

func TestResolveProjections_IntegerColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.role_id FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "RoleId", Type: "int32"}})
	assert.Equal(t, scans, []string{"&i.RoleId"})
}

func TestResolveProjections_BooleanColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.verified FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Verified", Type: "bool"}})
	assert.Equal(t, scans, []string{"&i.Verified"})
}

func TestResolveProjections_MultipleColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name, users.active FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64"},
		{Name: "Name", Type: "string"},
		{Name: "Active", Type: "pgtype.Bool"},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Active"})
}

func TestResolveProjections_AliasedTable(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.email, u.role_id FROM users u WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64"},
		{Name: "Email", Type: "string"},
		{Name: "RoleId", Type: "int32"},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Email", "&i.RoleId"})
}

func TestResolveProjections_ColumnAlias(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id as user_id, users.name as user_name FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "UserId", Type: "int64"},
		{Name: "UserName", Type: "string"},
	})
	assert.Equal(t, scans, []string{"&i.UserId", "&i.UserName"})
}

func TestResolveProjections_AllIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.age, users.role_id, users.org_id FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Age", Type: "pgtype.Int2"},
		{Name: "RoleId", Type: "int32"},
		{Name: "OrgId", Type: "int64"},
	})
	assert.Equal(t, scans, []string{"&i.Age", "&i.RoleId", "&i.OrgId"})
}

func TestResolveProjections_MixedAliasAndNoAlias(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id as user_id, users.email, users.verified FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "UserId", Type: "int64"},
		{Name: "Email", Type: "string"},
		{Name: "Verified", Type: "bool"},
	})
	assert.Equal(t, scans, []string{"&i.UserId", "&i.Email", "&i.Verified"})
}
