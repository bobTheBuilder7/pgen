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
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int64", Tag: `json:"id"`}})
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
	assert.Equal(t, fields, []gen.Field{{Name: "Name", Type: "string", Tag: `json:"name"`}})
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
	assert.Equal(t, fields, []gen.Field{{Name: "Email", Type: "string", Tag: `json:"email"`}})
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
	assert.Equal(t, fields, []gen.Field{{Name: "Status", Type: "int16", Tag: `json:"status"`}})
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
	assert.Equal(t, fields, []gen.Field{{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`}})
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
	assert.Equal(t, fields, []gen.Field{{Name: "Verified", Type: "bool", Tag: `json:"verified"`}})
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
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
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
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Email", Type: "string", Tag: `json:"email"`},
		{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`},
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
		{Name: "UserId", Type: "int64", Tag: `json:"user_id"`},
		{Name: "UserName", Type: "string", Tag: `json:"user_name"`},
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
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`},
		{Name: "OrgId", Type: "int64", Tag: `json:"org_id"`},
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
		{Name: "UserId", Type: "int64", Tag: `json:"user_id"`},
		{Name: "Email", Type: "string", Tag: `json:"email"`},
		{Name: "Verified", Type: "bool", Tag: `json:"verified"`},
	})
	assert.Equal(t, scans, []string{"&i.UserId", "&i.Email", "&i.Verified"})
}

func TestResolveProjections_StarSelectReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT * FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_StarSelectWithAliasReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT * FROM users u WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_TableDotStarReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.* FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveProjections_AliasDotStarReturnsError(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.* FROM users u WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `SELECT \*`)
}

func TestResolveReturning_InsertSingleColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name) VALUES ($1) RETURNING id;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{{Name: "Id", Type: "int64", Tag: `json:"id"`}})
	assert.Equal(t, scans, []string{"&i.Id"})
}

func TestResolveReturning_InsertMultipleColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Email", Type: "string", Tag: `json:"email"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Email"})
}

func TestResolveReturning_InsertNullableColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, age;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Age"})
}

func TestResolveReturning_UpdateReturning(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET name = $1 WHERE users.id = $2 RETURNING id, name;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveReturning_DeleteReturning(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM users WHERE users.id = $1 RETURNING id, name, active;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Active"})
}

// Test returning columns that are NOT in the INSERT column list
func TestResolveReturning_InsertReturnsColumnsNotInInsertList(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, age, login_count, active, verified;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "LoginCount", Type: "pgtype.Int4", Tag: `json:"login_count"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
		{Name: "Verified", Type: "bool", Tag: `json:"verified"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Age", "&i.LoginCount", "&i.Active", "&i.Verified"})
}

// Test returning only nullable columns
func TestResolveReturning_OnlyNullableColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING age, login_count, referrer_id, active;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "LoginCount", Type: "pgtype.Int4", Tag: `json:"login_count"`},
		{Name: "ReferrerId", Type: "pgtype.Int8", Tag: `json:"referrer_id"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
	})
	assert.Equal(t, scans, []string{"&i.Age", "&i.LoginCount", "&i.ReferrerId", "&i.Active"})
}

// Test returning all 11 columns from INSERT
func TestResolveReturning_InsertReturnsAllColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, status) VALUES ($1, $2, $3) RETURNING id, name, email, age, status, role_id, login_count, org_id, referrer_id, active, verified;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Email", Type: "string", Tag: `json:"email"`},
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "Status", Type: "int16", Tag: `json:"status"`},
		{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`},
		{Name: "LoginCount", Type: "pgtype.Int4", Tag: `json:"login_count"`},
		{Name: "OrgId", Type: "int64", Tag: `json:"org_id"`},
		{Name: "ReferrerId", Type: "pgtype.Int8", Tag: `json:"referrer_id"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
		{Name: "Verified", Type: "bool", Tag: `json:"verified"`},
	})
	assert.Equal(t, scans, []string{
		"&i.Id", "&i.Name", "&i.Email", "&i.Age", "&i.Status",
		"&i.RoleId", "&i.LoginCount", "&i.OrgId", "&i.ReferrerId",
		"&i.Active", "&i.Verified",
	})
}

// Test nullable vs not-null booleans in RETURNING
func TestResolveReturning_BooleanNullability(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET name = $1 WHERE users.id = $2 RETURNING active, verified;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
		{Name: "Verified", Type: "bool", Tag: `json:"verified"`},
	})
	assert.Equal(t, scans, []string{"&i.Active", "&i.Verified"})
}

// Test same RETURNING columns across all three DML types for consistency
func TestResolveReturning_ConsistentAcrossInsertUpdateDelete(t *testing.T) {
	c := testCliWithUsersSchema(t)

	returningCols := "RETURNING id, name, age, role_id, active;"

	insertSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email) VALUES ($1, $2) ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse INSERT: %v", err)
	}
	updateSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET email = $1 WHERE users.id = $2 ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse UPDATE: %v", err)
	}
	deleteSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM users WHERE users.id = $1 ` + returningCols)
	if err != nil {
		t.Fatalf("failed to parse DELETE: %v", err)
	}

	expectedFields := []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
	}
	expectedScans := []string{"&i.Id", "&i.Name", "&i.Age", "&i.RoleId", "&i.Active"}

	insertFields, insertScans, err := c.resolveReturning(insertSQL)
	assert.Nil(t, err)
	assert.Equal(t, insertFields, expectedFields)
	assert.Equal(t, insertScans, expectedScans)

	updateFields, updateScans, err := c.resolveReturning(updateSQL)
	assert.Nil(t, err)
	assert.Equal(t, updateFields, expectedFields)
	assert.Equal(t, updateScans, expectedScans)

	deleteFields, deleteScans, err := c.resolveReturning(deleteSQL)
	assert.Nil(t, err)
	assert.Equal(t, deleteFields, expectedFields)
	assert.Equal(t, deleteScans, expectedScans)
}

// Test alternating nullable/not-null int sizes in RETURNING
func TestResolveReturning_AlternatingNullableIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	// age(nullable smallint), status(not-null smallint), login_count(nullable int), role_id(not-null int), referrer_id(nullable bigint), org_id(not-null bigint)
	parsedSQL, err := postgresparser.ParseSQLStrict(`DELETE FROM users WHERE users.id = $1 RETURNING age, status, login_count, role_id, referrer_id, org_id;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},
		{Name: "Status", Type: "int16", Tag: `json:"status"`},
		{Name: "LoginCount", Type: "pgtype.Int4", Tag: `json:"login_count"`},
		{Name: "RoleId", Type: "int32", Tag: `json:"role_id"`},
		{Name: "ReferrerId", Type: "pgtype.Int8", Tag: `json:"referrer_id"`},
		{Name: "OrgId", Type: "int64", Tag: `json:"org_id"`},
	})
	assert.Equal(t, scans, []string{"&i.Age", "&i.Status", "&i.LoginCount", "&i.RoleId", "&i.ReferrerId", "&i.OrgId"})
}

// Test UPDATE returning the column being set + other columns
func TestResolveReturning_UpdateReturnsSetColumnAndOthers(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET name = $1, active = $2 WHERE users.id = $3 RETURNING id, name, email, active, verified;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveReturning(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "Email", Type: "string", Tag: `json:"email"`},
		{Name: "Active", Type: "pgtype.Bool", Tag: `json:"active"`},
		{Name: "Verified", Type: "bool", Tag: `json:"verified"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Email", "&i.Active", "&i.Verified"})
}

// JOIN tests

func TestResolveProjections_InnerJoinColumnsFromBothTables(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.name, p.id as post_id, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
		{Name: "PostId", Type: "int64", Tag: `json:"post_id"`},
		{Name: "Title", Type: "string", Tag: `json:"title"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.PostId", "&i.Title"})
}

func TestResolveProjections_LeftJoinForcesNullableOnJoinedTable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.name, p.id as post_id, p.title, p.published FROM users u LEFT JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},                     // users.id — NOT NULL, base table
		{Name: "Name", Type: "string", Tag: `json:"name"`},                // users.name — NOT NULL, base table
		{Name: "PostId", Type: "pgtype.Int8", Tag: `json:"post_id"`},      // posts.id — NOT NULL in schema but LEFT JOIN makes it nullable
		{Name: "Title", Type: "pgtype.Text", Tag: `json:"title"`},         // posts.title — NOT NULL in schema but LEFT JOIN makes it nullable
		{Name: "Published", Type: "pgtype.Bool", Tag: `json:"published"`}, // posts.published — NOT NULL but LEFT JOIN makes it nullable
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.PostId", "&i.Title", "&i.Published"})
}

func TestResolveProjections_LeftJoinNullableColumnStaysNullable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	// posts.body is already nullable in schema, LEFT JOIN should still produce pgtype.Text
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, p.body FROM users u LEFT JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Body", Type: "pgtype.Text", Tag: `json:"body"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Body"})
}

func TestResolveProjections_InnerJoinDoesNotForceNullable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, p.title, p.published FROM users u INNER JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Title", Type: "string", Tag: `json:"title"`},       // NOT NULL, INNER JOIN doesn't force nullable
		{Name: "Published", Type: "bool", Tag: `json:"published"`}, // NOT NULL, INNER JOIN doesn't force nullable
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Title", "&i.Published"})
}

func TestResolveProjections_RightJoinForcesNullableOnBaseTable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.name, p.title FROM users u RIGHT JOIN posts p ON u.id = p.user_id WHERE p.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "pgtype.Int8", Tag: `json:"id"`},     // users.id — NOT NULL but RIGHT JOIN makes base table nullable
		{Name: "Name", Type: "pgtype.Text", Tag: `json:"name"`}, // users.name — NOT NULL but RIGHT JOIN makes base table nullable
		{Name: "Title", Type: "string", Tag: `json:"title"`},    // posts.title — NOT NULL, joined table in RIGHT JOIN keeps types
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name", "&i.Title"})
}

func TestResolveProjections_JoinWithMixedNullability(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	// INNER JOIN: nullable columns stay nullable, not-null stay not-null
	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, u.age, p.title, p.body FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},           // users.id — NOT NULL
		{Name: "Age", Type: "pgtype.Int2", Tag: `json:"age"`},   // users.age — nullable in schema
		{Name: "Title", Type: "string", Tag: `json:"title"`},    // posts.title — NOT NULL
		{Name: "Body", Type: "pgtype.Text", Tag: `json:"body"`}, // posts.body — nullable in schema
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Age", "&i.Title", "&i.Body"})
}

// Subquery tests — projections

func TestResolveProjections_WhereInSubqueryColumns(t *testing.T) {
	// WHERE IN subquery: parser only exposes outer table columns in Columns
	// The subquery's table (posts) is not in Tables, only users is
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.title = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveProjections_ExistsSubqueryColumns(t *testing.T) {
	// EXISTS subquery: parser only exposes outer table columns
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE EXISTS (SELECT 1 FROM posts WHERE posts.user_id = users.id) AND users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "Name", Type: "string", Tag: `json:"name"`},
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.Name"})
}

func TestResolveProjections_ScalarSubqueryInSelect(t *testing.T) {
	// Scalar subquery in SELECT: the entire subquery becomes a column expression
	// Our code won't find a table.column pattern, so it falls through to "string"
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, (SELECT COUNT(*) FROM posts WHERE posts.user_id = users.id) as post_count FROM users WHERE users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	fields, scans, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.Nil(t, err)
	assert.Equal(t, fields, []gen.Field{
		{Name: "Id", Type: "int64", Tag: `json:"id"`},
		{Name: "PostCount", Type: "string", Tag: `json:"post_count"`}, // scalar subquery falls through to default string
	})
	assert.Equal(t, scans, []string{"&i.Id", "&i.PostCount"})
}

func TestResolveProjections_FromSubqueryErrors(t *testing.T) {
	// FROM subquery: parser creates a "sub" table (type=subquery) which is not in schema
	// Columns reference sub.id and sub.name which won't resolve
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT sub.id, sub.name FROM (SELECT users.id, users.name FROM users WHERE users.age > $1) sub WHERE sub.id = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, _, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `table sub not found`)
}
