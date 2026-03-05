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
| `main.go` | Entry point. Defines `cli` struct, `ToPascalCase` helper, directory constants (`db`, `query`, `schema`). |
| `run.go` | Orchestrates the two-phase pipeline. Creates `cli`, runs schema parsing then query codegen. |
| `parse_schema.go` | `c.parseSchema()` — reads a schema `.sql` file, validates it's a single `CREATE TABLE`, stores columns in `c.tablesCol`. |
| `parse_to_query.go` | `parseFileToQueries()` — reads a query `.sql` file, splits on `-- name:` annotations, returns `[]Query` (each has `name`, `t` for query type, `sql`). |
| `generator.go` | `c.generateCode()` — main code generation. For each query: parses SQL, resolves types, generates Go structs + methods using the `gen` library. Also contains `pgTypeToGoType()`, `c.resolveColumnGoType()`, `filterColumns()`, `findTable()`. |
| `resolver.go` | `c.resolveParams()` — resolves function parameter names and Go types from WHERE clause filter columns matched to `$1`, `$2`, etc. |
| `db/db.go` | Hand-written base for generated code: `DBTX` interface, `Queries` struct, `New()`, `WithTx()`. Do not edit generated files in `db/` other than this one. |

### Query annotation format

```sql
-- name: FunctionName :queryType
SELECT ...;
```

`queryType` is `:one`, `:many`, etc. Currently only `:one` is implemented for SELECT.

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

1. **Column type resolution** (`c.resolveColumnGoType` in `generator.go`): For each SELECT column, parses `table.column` or `alias.column`, looks up the table in `c.tablesCol`, finds the DDL column, maps its PG type to Go via `pgTypeToGoType`.

2. **Parameter type resolution** (`c.resolveParams` in `resolver.go`): For each `$N` parameter, matches it to the Nth filter-type `ColumnUsage` entry, then resolves the Go type through the same schema lookup chain.

3. **PG to Go type mapping** (`pgTypeToGoType` in `generator.go`):
   - `bigserial`, `bigint`, `int8` -> `int64`
   - `serial`, `integer`, `int`, `int4` -> `int32`
   - `smallserial`, `smallint`, `int2` -> `int16`
   - `boolean`, `bool` -> `bool`
   - `real`, `float4` -> `float32`
   - `double precision`, `float8` -> `float64`
   - `text`, `varchar`, `character varying`, `char`, `character`, `uuid` -> `string`
   - Anything else defaults to `string`

### What is NOT yet implemented

- Query types `:many`, `:exec`, `:execresult`
- INSERT / UPDATE / DELETE code generation (switch cases exist but are empty)
- Multi-table / JOIN queries (currently errors on `len(tables) != 1`)
- Nullable column handling (nullable columns should map to pointer types or `sql.Null*`)
- `go/format` on generated output (code is in `run.go` but commented out)

### Internal utility packages

- `syncmap/` — Generic type-safe wrapper around `sync.Map` with full API (`Load`, `Store`, `Delete`, `Range`, etc.)
- `bytesbufferpool/` — `sync.Pool`-backed `bytes.Buffer` pool (`Get`/`Put`)
- `assert/` — Lightweight generic test assertion library (see Testing section)

### Key dependencies

- `github.com/valkdb/postgresparser` — SQL parser. Key types: `ParsedQuery` (with `Columns`, `Tables`, `ColumnUsage`, `Parameters`, `DDLActions`), `DDLColumn` (with `Name`, `Type`, `Nullable`), `SelectColumn` (with `Expression`, `Alias`), `ColumnUsage` (with `TableAlias`, `Column`, `UsageType`), `Parameter` (with `Position`)
- `github.com/bobTheBuilder7/gen` — Go code generation. Key functions: `NewFile()`, `AddBlock()`, `Struct()`, `Const()`, `MethodFunc()`, `Call()`, `ErrCheck()`, `Line()`, `Arg()`, `Import()`, `Field{Name, Type, Tag}`, `WriteTo()`
- `github.com/jackc/pgx/v5` — PostgreSQL driver used in generated code

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
