# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is pgen

pgen is a Go code generator (similar to sqlc) that reads PostgreSQL schema definitions and query files, then generates type-safe Go code for database access using pgx/v5.

## Build, Run & Test

```bash
go build ./...          # build everything
go run .                # run the generator (reads db/schema/ and db/query/, outputs to db/)
go vet ./...            # lint
go test ./...           # run all tests
go test -run TestName   # run a single test
```

## Architecture

### The `cli` struct

The central struct is `cli` (defined in `main.go`). It holds `tablesCol`, a generic concurrent map (`syncmap.Map[string, []postgresparser.DDLColumn]`) that stores parsed schema column metadata keyed by table name. All methods that need schema info are receivers on `cli`.

### Two-phase pipeline (`run.go`)

`run()` creates a `cli` instance and orchestrates two sequential phases using `errgroup` for concurrency within each phase:

1. **Schema parsing** — Reads `db/schema/*.sql`, parses each `CREATE TABLE` via `postgresparser.ParseSQLStrict`, stores column metadata into `c.tablesCol`. One table per file.

2. **Query parsing & code generation** — Reads `db/query/*.sql`, parses query annotations into `Query` structs, then calls `c.generateCode()` to produce Go source written to `db/<name>.go`.

Phase 1 must fully complete before phase 2 starts (phase 2 depends on `tablesCol` being populated).

### File-by-file breakdown

| File | Purpose |
|---|---|
| `main.go` | Entry point. Defines `cli` struct, directory constants (`db`, `query`, `schema`), `sqlConstSuffix` const. |
| `run.go` | Orchestrates the two-phase pipeline. Creates `cli`, runs schema parsing then query codegen. |
| `parse_schema.go` | `c.parseSchema()` — reads a schema `.sql` file, validates it's a single `CREATE TABLE`, stores columns in `c.tablesCol`. |
| `parse_query.go` | `parseFileToQueries()` — reads a query `.sql` file, splits on `-- name:` annotations, returns `[]Query` (each has `name`, `t` for query type, `sql`). |
| `generator.go` | `c.generateCode()` — main code generation. Contains `pgTypeToGoType()`, `generateExec()`, `filterColumns()`, `buildCallArgsString()`, `stringersFromStrings()`. Handles SELECT `:one`/`:many`/`:exec`/`:execresult`, delegates INSERT/UPDATE/DELETE to `generateExec`. |
| `resolve_projections.go` | `c.resolveProjections()` — resolves Go struct fields and scan fields from SELECT columns. Also contains `resolveColumnGoType()`, `isOuterJoinNullable()`, `resolveReturning()`. |
| `resolve_params.go` | `c.resolveParams()` — resolves function parameter names and Go types. Handles SELECT/DELETE (WHERE filters), UPDATE (SET + WHERE), INSERT (positional via `InsertColumns`). Also contains `resolveInsertParams()`. |
| `utils/utils.go` | `ToPascalCase` helper function. |
| `db/db.go` | Hand-written base for generated code: `DBTX` interface, `Queries` struct, `New()`, `WithTx()`. Do not edit generated files in `db/` other than this one. |
| `experiments/main.go` | Test harness for exploring `postgresparser` behavior interactively. Listed in `.gitignore`. |

### Query annotation format

```sql
-- name: FunctionName :queryType
SELECT ...;
```

`queryType` determines the generated method's return type and behavior:

| Query Type | Applies To | Returns |
|---|---|---|
| `:one` | SELECT, INSERT/UPDATE/DELETE with RETURNING | `(RowStruct, error)` — uses `QueryRow` + `Scan` |
| `:many` | SELECT, INSERT/UPDATE/DELETE with RETURNING | `([]RowStruct, error)` — uses `Query` + rows iteration |
| `:exec` | SELECT, INSERT, UPDATE, DELETE | `error` — uses `Exec`, discards result |
| `:execresult` | SELECT, INSERT, UPDATE, DELETE | `(pgconn.CommandTag, error)` — uses `Exec`, returns command tag |

### What `:one` SELECT generates

For a query like:
```sql
-- name: GetUserByID :one
SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $2;
```

The generator produces:
- A **result struct** (`GetUserByIDRow`) with fields derived from projected columns, types resolved via schema
- A **SQL const** (`GetUserByIDSQL`)
- A **method** on `*Queries` using `QueryRow` + `Scan`

### How type resolution works

1. **Column type resolution** (`c.resolveColumnGoType` in `resolve_projections.go`): For each SELECT column, parses `table.column` or `alias.column`, looks up the table in `c.tablesCol`, finds the DDL column, maps its PG type to Go via `pgTypeToGoType`. Handles outer join nullability via `isOuterJoinNullable`.

2. **Parameter type resolution** (`c.resolveParams` in `resolve_params.go`): For SELECT/DELETE, matches `$N` to the Nth filter-type `ColumnUsage`. For UPDATE, SET columns come first, then filter columns. For INSERT, uses `InsertColumns` positionally. Handles outer join nullability.

3. **RETURNING resolution** (`c.resolveReturning` in `resolve_projections.go`): For INSERT/UPDATE/DELETE with RETURNING, uses `ColumnUsage` entries with `UsageType == "returning"` to resolve struct fields.

4. **PG to Go type mapping** (`pgTypeToGoType` in `generator.go`):
   - `bigserial`, `bigint`, `int8` -> `int64`
   - `serial`, `integer`, `int`, `int4` -> `int32`
   - `smallserial`, `smallint`, `int2` -> `int16`
   - `boolean`, `bool` -> `bool`
   - `real`, `float4` -> `float32`
   - `double precision`, `float8` -> `float64`
   - `text`, `varchar`, `character varying`, `char`, `character` -> `string` (nullable: `pgtype.Text`)
   - `uuid` -> `string` (nullable: `pgtype.UUID`)
   - Anything else defaults to `string`

### Multi-table / JOIN support

Queries can reference multiple tables via JOINs. Table aliases are resolved throughout the pipeline. Outer join nullability is handled automatically:

- `LEFT JOIN` → joined table's columns forced nullable
- `RIGHT JOIN` → base table's columns forced nullable
- `FULL JOIN` → both sides forced nullable
- `JOIN` / `INNER JOIN` / `CROSS JOIN` → use schema nullability as-is

### Validation

- `SELECT *` and `table.*` are rejected with an error — columns must be explicitly listed
- String literals in SELECT (e.g. `'foo'`) resolve to `string`

### What is NOT yet implemented

- Error on duplicate query names within a file
- Params struct for 3+ parameters (like sqlc generates)

### Internal utility packages

- `syncmap/` — Generic type-safe wrapper around `sync.Map` with full API (`Load`, `Store`, `Delete`, `Range`, etc.)
- `bytesbufferpool/` — `sync.Pool`-backed `bytes.Buffer` pool (`Get`/`Put`)
- `assert/` — Lightweight generic test assertion library (see Testing section)

### Key dependencies

- `github.com/valkdb/postgresparser` — SQL parser. Key types: `ParsedQuery` (with `Columns`, `Tables`, `ColumnUsage`, `Parameters`, `DDLActions`, `InsertColumns`, `Returning`, `Command`), `DDLColumn` (with `Name`, `Type`, `Nullable`), `SelectColumn` (with `Expression`, `Alias`), `ColumnUsage` (with `TableAlias`, `Column`, `UsageType`), `TableRef` (with `Name`, `Alias`, `JoinType`), `Parameter` (with `Position`). `ColumnUsageType` constants: `ColumnUsageTypeFilter`, `ColumnUsageTypeDMLSet`, `ColumnUsageTypeReturning`.
- `github.com/bobTheBuilder7/gen` — Go code generation. Key functions: `NewFile()`, `AddBlock()`, `Struct()`, `Const()`, `MethodFunc()`, `Call()`, `Line()`, `Arg()`, `Import()`, `String()`, `Field{Name, Type, Tag}`, `WriteTo()`
- `github.com/jackc/pgx/v5` — PostgreSQL driver used in generated code
- `github.com/jackc/pgx/v5/pgtype` — Nullable Go types: `pgtype.Int2`, `Int4`, `Int8`, `Bool`, `Float4`, `Float8`, `Text`, `UUID`
- `github.com/jackc/pgx/v5/pgconn` — `pgconn.CommandTag` used by `:execresult` return type

### Directory conventions

- `db/schema/*.sql` — One `CREATE TABLE` per file
- `db/query/*.sql` — Query files with `-- name:` annotations; each generates a corresponding `db/<name>.go`
- `db/db.go` — Hand-written, not generated

## Testing

We use our own `assert/` package (`github.com/bobTheBuilder7/pgen/assert`). Available functions:

```go
assert.Equal(t, got, want)       // deep equality via reflect.DeepEqual
assert.NotEqual(t, got, want)
assert.True(t, got)
assert.False(t, got)
assert.Nil(t, got)
assert.NotNil(t, got)
assert.ErrorIs(t, got, want)
assert.ErrorAs(t, got, target)
assert.MatchesRegexp(t, got, pattern)
```

All functions are generic where applicable. Always use this package instead of manual `if` checks in tests.

### Test patterns

Tests create a `cli` instance directly and populate `tablesCol` by parsing a hardcoded SQL schema string:

```go
const usersSchemaSQL = `CREATE TABLE users (...);`

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
```

Queries in tests are also hardcoded as SQL strings and parsed with `postgresparser.ParseSQLStrict`. No file I/O in tests.
