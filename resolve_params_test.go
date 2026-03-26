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

const postsSchemaSQL = `CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT,
    user_id BIGINT NOT NULL,
    published BOOLEAN NOT NULL DEFAULT false
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

func testCliWithUsersAndPostsSchema(t *testing.T) *cli {
	t.Helper()
	c := testCliWithUsersSchema(t)
	parsed, err := postgresparser.ParseSQLStrict(postsSchemaSQL)
	if err != nil {
		t.Fatalf("failed to parse posts schema: %v", err)
	}
	c.tablesCol.Store("posts", parsed.DDLActions[0].ColumnDetails)
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
	assert.Equal(t, types, []string{"int64", "pgtype.Bool"})
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
	assert.Equal(t, types, []string{"pgtype.Int8"})
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
	assert.Equal(t, types, []string{"int32", "string", "pgtype.Bool"})
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

func TestResolveParams_UpdateSetAndWhere(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET name = $1 WHERE users.id = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "id"})
	assert.Equal(t, types, []string{"string", "int64"})
}

func TestResolveParams_UpdateMultipleSetColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET name = $1, email = $2, active = $3 WHERE users.id = $4;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "active", "id"})
	assert.Equal(t, types, []string{"string", "string", "pgtype.Bool", "int64"})
}

func TestResolveParams_UpdateMultipleWhereColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`UPDATE users SET verified = $1 WHERE users.id = $2 AND users.org_id = $3;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"verified", "id", "org_id"})
	assert.Equal(t, types, []string{"bool", "int64", "int64"})
}

func TestResolveParams_InsertSimple(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, status) VALUES ($1, $2, $3);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "status"})
	assert.Equal(t, types, []string{"string", "string", "int16"})
}

func TestResolveParams_InsertWithNullableColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, age, active) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "age", "active"})
	assert.Equal(t, types, []string{"string", "string", "pgtype.Int2", "pgtype.Bool"})
}

func TestResolveParams_InsertAllIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, status, role_id, org_id) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "status", "role_id", "org_id"})
	assert.Equal(t, types, []string{"string", "string", "int16", "int32", "int64"})
}

func TestResolveParams_InsertSingleColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name) VALUES ($1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_InsertBooleans(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, active, verified) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "active", "verified"})
	assert.Equal(t, types, []string{"string", "string", "pgtype.Bool", "bool"})
}

func TestResolveParams_InsertNullableIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, email, age, login_count, referrer_id) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "email", "age", "login_count", "referrer_id"})
	assert.Equal(t, types, []string{"string", "string", "pgtype.Int2", "pgtype.Int4", "pgtype.Int8"})
}

func TestResolveParams_InsertAllColumns(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (id, name, email, age, status, role_id, login_count, org_id, referrer_id, active, verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "name", "email", "age", "status", "role_id", "login_count", "org_id", "referrer_id", "active", "verified"})
	assert.Equal(t, types, []string{"int64", "string", "string", "pgtype.Int2", "int16", "int32", "pgtype.Int4", "int64", "pgtype.Int8", "pgtype.Bool", "bool"})
}

func TestResolveParams_InsertMixedNullability(t *testing.T) {
	c := testCliWithUsersSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`INSERT INTO users (name, age, status, login_count, role_id) VALUES ($1, $2, $3, $4, $5);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "age", "status", "login_count", "role_id"})
	assert.Equal(t, types, []string{"string", "pgtype.Int2", "int16", "pgtype.Int4", "int32"})
}

// JOIN param tests

func TestResolveParams_JoinSingleParam(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_JoinParamsFromBothTables(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1 AND p.title = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "title"})
	assert.Equal(t, types, []string{"int64", "string"})
}

func TestResolveParams_LeftJoinParamFromJoinedTable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT u.id, p.title FROM users u LEFT JOIN posts p ON u.id = p.user_id WHERE u.id = $1 AND p.published = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id", "published"})
	// published is NOT NULL in schema but LEFT JOIN makes posts columns nullable
	assert.Equal(t, types, []string{"int64", "pgtype.Bool"})
}

// Subquery param tests

func TestResolveParams_ExistsSubqueryParamOnOuterQuery(t *testing.T) {
	// EXISTS subquery with param on the outer WHERE clause — resolves correctly
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE EXISTS (SELECT 1 FROM posts WHERE posts.user_id = users.id) AND users.id = $1;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
}

func TestResolveParams_WhereInSubqueryParam(t *testing.T) {
	// WHERE IN subquery: $1 is inside the subquery (posts.title = $1)
	// Should resolve to posts.name (string), not the outer users.id (int64)
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.title = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"title"})
	assert.Equal(t, types, []string{"string"})
}

func TestResolveParams_NotInSubqueryParam(t *testing.T) {
	// NOT IN subquery: $1 is inside the subquery (posts.published = $1)
	// Should resolve to posts.published (bool), not users.id (int64)
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.id NOT IN (SELECT posts.user_id FROM posts WHERE posts.published = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"published"})
	assert.Equal(t, types, []string{"bool"})
}

func TestResolveParams_MixedOuterAndSubqueryParams(t *testing.T) {
	// $1 is on outer WHERE (users.name), $2 is in subquery (posts.name)
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.name = $1 AND users.id IN (SELECT posts.user_id FROM posts WHERE posts.title = $2);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name", "title"})
	assert.Equal(t, types, []string{"string", "string"})
}

func TestResolveParams_SubqueryParamWithNullableColumn(t *testing.T) {
	// $1 is inside subquery referencing a nullable column (posts.body)
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.body = $1);`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"body"})
	assert.Equal(t, types, []string{"pgtype.Text"})
}

func TestResolveParams_SubqueryParamAndOuterParamReversed(t *testing.T) {
	// $1 is in subquery (posts.name), $2 is on outer WHERE (users.name)
	c := testCliWithUsersAndPostsSchema(t)

	parsedSQL, err := postgresparser.ParseSQLStrict(`SELECT users.id, users.name FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.title = $1) AND users.name = $2;`)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"title", "name"})
	assert.Equal(t, types, []string{"string", "string"})
}

// Named parameter tests — convertNamedParams

func TestConvertNamedParams_NoParams(t *testing.T) {
	sql := `SELECT users.id FROM users;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, sql)
	assert.Nil(t, names)
}

func TestConvertNamedParams_PositionalOnly(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.id = $1;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, sql)
	assert.Nil(t, names)
}

func TestConvertNamedParams_SingleNamed(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.id = $1;`)
	assert.Equal(t, names, []string{"id"})
}

func TestConvertNamedParams_MultipleNamed(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.id = @id AND users.name = @name;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.id = $1 AND users.name = $2;`)
	assert.Equal(t, names, []string{"id", "name"})
}

func TestConvertNamedParams_Underscore(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.role_id = @role_id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.role_id = $1;`)
	assert.Equal(t, names, []string{"role_id"})
}

func TestConvertNamedParams_MixedError(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.id = $1 AND users.name = @name;`
	_, _, err := convertNamedParams(sql)
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `cannot mix`)
}

func TestConvertNamedParams_InsertNamed(t *testing.T) {
	sql := `INSERT INTO users (name, email, age) VALUES (@name, @email, @age);`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO users (name, email, age) VALUES ($1, $2, $3);`)
	assert.Equal(t, names, []string{"name", "email", "age"})
}

func TestConvertNamedParams_UpdateNamed(t *testing.T) {
	sql := `UPDATE users SET name = @name WHERE users.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `UPDATE users SET name = $1 WHERE users.id = $2;`)
	assert.Equal(t, names, []string{"name", "id"})
}

func TestConvertNamedParams_DeleteNamed(t *testing.T) {
	sql := `DELETE FROM users WHERE users.id = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `DELETE FROM users WHERE users.id = $1;`)
	assert.Equal(t, names, []string{"id"})
}

func TestConvertNamedParams_ThreeParams(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.role_id = @role_id AND users.name = @user_name AND users.active = @is_active;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.role_id = $1 AND users.name = $2 AND users.active = $3;`)
	assert.Equal(t, names, []string{"role_id", "user_name", "is_active"})
}

func TestConvertNamedParams_ReturningNamed(t *testing.T) {
	sql := `INSERT INTO users (name, email) VALUES (@name, @email) RETURNING id, name;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name;`)
	assert.Equal(t, names, []string{"name", "email"})
}

// End-to-end named param tests — convertNamedParams + ParseSQLStrict + resolveParams

func TestNamedParams_SelectSingleParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `SELECT users.id, users.name FROM users WHERE users.id = @user_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	names, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	// resolveParams returns column-derived names, override with named params
	assert.Equal(t, names, []string{"id"})
	assert.Equal(t, types, []string{"int64"})
	// namedParams has the user-specified names
	assert.Equal(t, namedParams, []string{"user_id"})
}

func TestNamedParams_SelectMultipleParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `SELECT users.id FROM users WHERE users.id = @user_id AND users.name = @user_name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"user_id", "user_name"})
}

func TestNamedParams_SelectWithAliasedTable(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `SELECT u.id, u.email FROM users u WHERE u.referrer_id = @ref_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"pgtype.Int8"})
	assert.Equal(t, namedParams, []string{"ref_id"})
}

func TestNamedParams_SelectNullableColumn(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `SELECT users.id FROM users WHERE users.active = @is_active AND users.age = @min_age;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"pgtype.Bool", "pgtype.Int2"})
	assert.Equal(t, namedParams, []string{"is_active", "min_age"})
}

func TestNamedParams_InsertParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `INSERT INTO users (name, email, status) VALUES (@user_name, @user_email, @user_status);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "string", "int16"})
	assert.Equal(t, namedParams, []string{"user_name", "user_email", "user_status"})
}

func TestNamedParams_InsertNullableParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `INSERT INTO users (name, email, age, active) VALUES (@name, @email, @age, @active);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "string", "pgtype.Int2", "pgtype.Bool"})
	assert.Equal(t, namedParams, []string{"name", "email", "age", "active"})
}

func TestNamedParams_UpdateSetAndWhere(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `UPDATE users SET name = @new_name WHERE users.id = @user_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "int64"})
	assert.Equal(t, namedParams, []string{"new_name", "user_id"})
}

func TestNamedParams_UpdateMultipleSets(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `UPDATE users SET name = @new_name, email = @new_email, active = @is_active WHERE users.id = @user_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "string", "pgtype.Bool", "int64"})
	assert.Equal(t, namedParams, []string{"new_name", "new_email", "is_active", "user_id"})
}

func TestNamedParams_DeleteSingleParam(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `DELETE FROM users WHERE users.id = @user_id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64"})
	assert.Equal(t, namedParams, []string{"user_id"})
}

func TestNamedParams_DeleteMultipleParams(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `DELETE FROM users WHERE users.id = @user_id AND users.name = @user_name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"user_id", "user_name"})
}

func TestNamedParams_AllIntSizes(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `SELECT users.email FROM users WHERE users.id = @big AND users.status = @small AND users.role_id = @medium AND users.verified = @flag;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "int16", "int32", "bool"})
	assert.Equal(t, namedParams, []string{"big", "small", "medium", "flag"})
}

func TestNamedParams_JoinParams(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	sql := `SELECT u.id, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = @user_id AND p.title = @post_title;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "string"})
	assert.Equal(t, namedParams, []string{"user_id", "post_title"})
}

func TestNamedParams_LeftJoinForcesNullable(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	sql := `SELECT u.id, p.title FROM users u LEFT JOIN posts p ON u.id = p.user_id WHERE u.id = @user_id AND p.published = @is_published;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"int64", "pgtype.Bool"})
	assert.Equal(t, namedParams, []string{"user_id", "is_published"})
}

func TestNamedParams_InsertWithReturning(t *testing.T) {
	c := testCliWithUsersSchema(t)

	sql := `INSERT INTO users (name, email) VALUES (@user_name, @user_email) RETURNING id, name;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string", "string"})
	assert.Equal(t, namedParams, []string{"user_name", "user_email"})
}

func TestNamedParams_SubqueryParam(t *testing.T) {
	c := testCliWithUsersAndPostsSchema(t)

	sql := `SELECT users.id, users.name FROM users WHERE users.id IN (SELECT posts.user_id FROM posts WHERE posts.title = @post_title);`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)

	parsedSQL, err := postgresparser.ParseSQLStrict(converted)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}

	_, types, err := c.resolveParams(parsedSQL)
	assert.Nil(t, err)
	assert.Equal(t, types, []string{"string"})
	assert.Equal(t, namedParams, []string{"post_title"})
}

// --- Duplicate param tests ---

func TestConvertNamedParams_DuplicateNameMapsToSamePosition(t *testing.T) {
	// @val used twice → both should map to $1 (same slot)
	sql := `SELECT users.id FROM users WHERE users.name = @val AND users.email = @val;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.name = $1 AND users.email = $1;`)
	assert.Equal(t, names, []string{"val"})
}

func TestConvertNamedParams_DuplicateNameInUpdate(t *testing.T) {
	// @id used in two WHERE conditions → same $2 slot
	sql := `UPDATE users SET name = @name WHERE users.id = @id AND users.age = @id;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `UPDATE users SET name = $1 WHERE users.id = $2 AND users.age = $2;`)
	assert.Equal(t, names, []string{"name", "id"})
}

func TestConvertNamedParams_DuplicateNameInInsert(t *testing.T) {
	// All distinct names — no dedup expected
	sql := `INSERT INTO users (name, email) VALUES (@name, @email);`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `INSERT INTO users (name, email) VALUES ($1, $2);`)
	assert.Equal(t, names, []string{"name", "email"})
}

func TestConvertNamedParams_ThreeDistinctNames(t *testing.T) {
	sql := `SELECT users.id FROM users WHERE users.name = @a AND users.email = @b AND users.age = @c;`
	converted, names, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, converted, `SELECT users.id FROM users WHERE users.name = $1 AND users.email = $2 AND users.age = $3;`)
	assert.Equal(t, names, []string{"a", "b", "c"})
}

func TestResolveParams_DuplicatePositionalParam(t *testing.T) {
	// $1 used for two different columns — should generate only one function param
	c := testCliWithUsersSchema(t)
	const sql = `SELECT users.id FROM users WHERE users.name = $1 AND users.email = $1;`
	parsed, err := postgresparser.ParseSQLStrict(sql)
	assert.Nil(t, err)
	names, types, err := c.resolveParams(parsed)
	assert.Nil(t, err)
	assert.Equal(t, names, []string{"name"})
	assert.Equal(t, types, []string{"string"})
}

func TestNamedParams_DuplicateNameSelectOneFunctionParam(t *testing.T) {
	// @val used twice in WHERE → one function param, SQL uses $1 twice
	c := testCliWithUsersSchema(t)
	const sql = `SELECT users.id FROM users WHERE users.name = @val AND users.email = @val;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"val"})
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	names, types, err := c.resolveParams(parsed)
	assert.Nil(t, err)
	// namedParams overrides resolved names
	names = namedParams
	assert.Equal(t, names, []string{"val"})
	assert.Equal(t, types, []string{"string"})
}

func TestNamedParams_DuplicateNameUpdateOneFunctionParam(t *testing.T) {
	// @id used twice in WHERE → two $N slots collapse to one function param after dedup
	c := testCliWithUsersSchema(t)
	const sql = `UPDATE users SET name = @new_name WHERE users.id = @id AND users.age = @id;`
	converted, namedParams, err := convertNamedParams(sql)
	assert.Nil(t, err)
	assert.Equal(t, namedParams, []string{"new_name", "id"})
	assert.Equal(t, converted, `UPDATE users SET name = $1 WHERE users.id = $2 AND users.age = $2;`)
	parsed, err := postgresparser.ParseSQLStrict(converted)
	assert.Nil(t, err)
	names, types, err := c.resolveParams(parsed)
	assert.Nil(t, err)
	// Override with named param names
	names = namedParams
	assert.Equal(t, names, []string{"new_name", "id"})
	assert.Equal(t, types, []string{"string", "int64"})
}
