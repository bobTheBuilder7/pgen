# pgen

Type-safe Go code generation from PostgreSQL queries. Write SQL, get Go — no runtime reflection, no manual scanning.

Similar to [sqlc](https://sqlc.dev), but built specifically for [pgx/v5](https://github.com/jackc/pgx).

## How it works

1. Define your schema in SQL
2. Write your queries with simple annotations
3. Run `pgen`
4. Get type-safe Go code with proper pgx types

```
db/
├── schema/
│   └── users.sql          -- CREATE TABLE definitions
├── query/
│   └── users.sql          -- annotated SQL queries
├── db.go                  -- base file (generated)
└── users.go               -- generated code
```

## Quick start

```bash
go install github.com/bobTheBuilder7/pgen@latest
```

### 1. Define your schema

```sql
-- db/schema/users.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    age INTEGER
);
```

### 2. Write your queries

Annotate each query with `-- name: FunctionName :queryType`.

You can use **positional parameters** (`$1`, `$2`) or **named parameters** (`@name`, `@user_id`):

```sql
-- db/query/users.sql

-- name: GetUserByID :one
SELECT users.id, users.name
FROM users
WHERE users.id = $1 AND users.name = $2;

-- name: GetUserByName :one
SELECT users.id, users.name
FROM users
WHERE users.id = @user_id AND users.name = @user_name;

-- name: ListUsers :many
SELECT users.id, users.name FROM users;

-- name: CreateUser :exec
INSERT INTO users (name, age) VALUES ($1, $2);

-- name: DeleteUserByID :exec
DELETE FROM users WHERE users.id = $1;

-- name: UpdateUserName :exec
UPDATE users SET name = $1 WHERE users.id = $2;

-- name: InsertUserReturning :one
INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name;

-- name: DeleteUserByName :execresult
DELETE FROM users WHERE users.name = $1;
```

### 3. Generate

```bash
pgen
```

### 4. Use the generated code

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v5"
    "your/project/db"
)

func main() {
    ctx := context.Background()

    conn, err := pgx.Connect(ctx, "postgres://localhost:5432/mydb")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(ctx)

    queries := db.New(conn)

    // :one — returns a single row
    user, err := queries.GetUserByID(ctx, 1, "alice")
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("User: %s (id=%d)", user.Name, user.Id)

    // :many — returns a slice
    users, err := queries.ListUsers(ctx)
    if err != nil {
        log.Fatal(err)
    }
    for _, u := range users {
        log.Printf("- %s", u.Name)
    }

    // :exec — returns only an error
    err = queries.CreateUser(ctx, "bob", pgtype.Int4{Int32: 25, Valid: true})

    // :execresult — returns pgconn.CommandTag + error
    tag, err := queries.DeleteUserByName(ctx, "bob")
    log.Printf("Deleted %d rows", tag.RowsAffected())
}
```

## Generated output

For a `:one` SELECT query, pgen generates a result struct, a SQL constant, and a method on `*Queries`:

```go
type GetUserByIDRow struct {
    Id   int64  `json:"id"`
    Name string `json:"name"`
}

const GetUserByIDSQL = "SELECT users.id, users.name FROM users WHERE users.id = $1 and users.name = $2;"

func (q *Queries) GetUserByID(ctx context.Context, id int64, name string) (GetUserByIDRow, error) {
    row := q.db.QueryRow(ctx, GetUserByIDSQL, id, name)
    var i GetUserByIDRow
    err := row.Scan(&i.Id, &i.Name)
    return i, err
}
```

When using named parameters, the `@name` values become the function parameter names:

```go
// Query: WHERE users.id = @user_id AND users.name = @user_name

const GetUserByNameSQL = "SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $2;"

func (q *Queries) GetUserByName(ctx context.Context, user_id int64, user_name string) (GetUserByNameRow, error) {
    // ...
}
```

Named parameters are converted to positional `$N` in the SQL constant (pgx requires positional params), but your function signature uses the readable names you chose.

For `:many`, it returns a slice and iterates over `rows`:

```go
func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
    rows, err := q.db.Query(ctx, ListUsersSQL)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var items []ListUsersRow
    for rows.Next() {
        var i ListUsersRow
        if err := rows.Scan(&i.Id, &i.Name); err != nil {
            return nil, err
        }
        items = append(items, i)
    }
    return items, rows.Err()
}
```

## Query types

| Annotation | Returns | Use case |
|---|---|---|
| `:one` | `(Row, error)` | Single row — SELECT, INSERT/UPDATE/DELETE with RETURNING |
| `:many` | `([]Row, error)` | Multiple rows — SELECT, or RETURNING multiple rows |
| `:exec` | `error` | No result needed — INSERT, UPDATE, DELETE |
| `:execresult` | `(pgconn.CommandTag, error)` | When you need affected row count |

## Features

- **Named parameters** — Use `@name` style params for readable function signatures, or `$N` positional params
- **JSON tags** — Generated structs include `json:"snake_case"` tags
- **JOINs** — Multi-table queries with aliases are fully supported
- **RETURNING** — INSERT/UPDATE/DELETE with RETURNING clauses generate row structs
- **Subqueries** — WHERE IN, EXISTS, and scalar subqueries are supported
- **Nullable columns** — Nullable columns map to pgx types (`pgtype.Int4`, `pgtype.Text`, etc.) instead of raw Go types
- **All DML** — SELECT, INSERT, UPDATE, and DELETE are all supported
- **Concurrent parsing** — Schema and query files are parsed in parallel using errgroups
- **pgx/v5 native** — Generated code uses pgx directly, no `database/sql` wrapper

## Named parameters

Instead of positional `$1, $2, $3`, you can use `@name` style parameters for more readable queries and generated code:

```sql
-- Positional: parameter names are derived from column names
-- name: GetUser :one
SELECT users.id, users.name FROM users WHERE users.id = $1;
-- generates: func GetUser(ctx context.Context, id int64)

-- Named: you choose the parameter names
-- name: GetUser :one
SELECT users.id, users.name FROM users WHERE users.id = @user_id;
-- generates: func GetUser(ctx context.Context, user_id int64)
```

Named parameters work with all query types — SELECT, INSERT, UPDATE, DELETE:

```sql
-- name: CreateUser :exec
INSERT INTO users (name, age) VALUES (@name, @age);

-- name: UpdateUserName :exec
UPDATE users SET name = @new_name WHERE users.id = @user_id;

-- name: DeleteUser :exec
DELETE FROM users WHERE users.id = @user_id;
```

**Note:** You cannot mix `@name` and `$N` in the same query — pick one style per query.

## Type mapping

| PostgreSQL | Go | Go (nullable) |
|---|---|---|
| `bigserial`, `bigint`, `int8` | `int64` | `pgtype.Int8` |
| `serial`, `integer`, `int4` | `int32` | `pgtype.Int4` |
| `smallserial`, `smallint`, `int2` | `int16` | `pgtype.Int2` |
| `boolean`, `bool` | `bool` | `pgtype.Bool` |
| `real`, `float4` | `float32` | `pgtype.Float4` |
| `double precision`, `float8` | `float64` | `pgtype.Float8` |
| `text`, `varchar`, `char` | `string` | `pgtype.Text` |
| `uuid` | `string` | `pgtype.UUID` |

## Project structure

```
pgen/
├── main.go               # entry point, cli struct
├── run.go                # two-phase pipeline orchestration
├── parse_schema.go       # schema file parsing
├── parse_query.go        # query file parsing with annotation splitting
├── generator.go          # code generation engine
├── resolve_params.go     # WHERE clause parameter resolution
├── resolve_projections.go # SELECT column type resolution
├── syncmap/              # generic type-safe sync.Map wrapper
├── bytesbufferpool/      # pooled bytes.Buffer
└── assert/               # lightweight generic test assertions
```

## License

MIT
