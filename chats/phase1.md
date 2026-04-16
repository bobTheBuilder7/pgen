# Phase 1 Session — Knowledge Base

## Architecture (as of this session)

### Entry Point & Central Struct
- `main.go` — defines `cli` struct with `tablesCol syncmap.Map[string, []dbColumn]` and `db *sql.DB`
- `dbColumn` struct (defined in `run.go`): `Name string`, `Type string`, `Nullable bool`
- `cli` no longer uses `postgresparser.DDLColumn` — schema is loaded from a live pglite DB (see below)

### Pipeline (run.go)
The pipeline changed significantly from the original CLAUDE.md description:

1. **Migration phase** — reads `db/migrations/*.up.sql` files (sorted), runs each against an in-memory pglite DB (`bobTheBuilder7/gopglite` driver via `database/sql`)
2. **Schema load** — `loadSchemaFromDB()` queries `information_schema.columns` from pglite and populates `c.tablesCol`
3. **Query phase** — reads `db/queries/*.sql`, parses to `[]query`, validates each via `testQueryAgainstDB()` (PREPARE statement against pglite), then generates Go code

**Important:** Phase 1 (migrations) must complete before phase 2 (queries) because `tablesCol` must be populated first. Migration files sorted by filename (alphabetically, so numeric prefix ordering works).

### Directory Layout
```
db/migrations/*.up.sql   — migration files, applied in sorted order
db/queries/*.sql         — query files, each generates a db/<name>.go
db/db.go                 — generated base file (DBTX interface, Queries, New, WithTx)
```

### Key Files
| File | Purpose |
|---|---|
| `main.go` | Entry point; `cli` struct; directory/const definitions |
| `run.go` | Pipeline orchestration + `dbColumn` struct |
| `migrations.go` | `runMigration()` — executes a migration file against pglite |
| `load_schema.go` | `loadSchemaFromDB()` — loads schema from pglite information_schema |
| `parse_query_file.go` | `parseFileToQueries()` — splits SQL file on `-- name:` annotations |
| `prepare_query.go` | `testQueryAgainstDB()` — validates query via PREPARE against pglite |
| `generator.go` | `generateCode()`, `generateExec()`, `convertNamedParams()`, `pgTypeToGoType()` |
| `resolve_projections.go` | `resolveProjections()`, `resolveColumnGoType()`, CTE resolution |
| `resolve_params.go` | `resolveParams()`, `buildSubqueryParamMap()`, `buildCTEParamMap()` |
| `base_file.go` | `generateBaseFile()` — generates db/db.go (pgx and std modes) |
| `sort.go` | `sortMigrations()` — sorts migration filenames |

---

## Features Implemented This Session

### CTE Support (full)

**Problem:** `go run .` failed with `table city_directors not found in schema` because a query used a CTE alias as a column prefix in SELECT, and the code tried to look it up as a real table.

**Root cause:** Two separate issues:
1. `resolveProjections` / `resolveColumnGoType` had no concept of CTEs — when a column like `city_directors.id` was encountered, it looked up `city_directors` in `tablesCol` which only contains real schema tables.
2. `resolveParams` had no concept of CTEs — parameters inside a CTE's WHERE clause appear in top-level `parsedSQL.Parameters` but have no matching entry in the outer query's filter `ColumnUsage`. Without CTE param resolution, they were incorrectly assigned as LIMIT/OFFSET params.

---

## postgresparser CTE API (discovered via experiments)

### `parsedSQL.CTEs []postgresparser.CTE`
Each CTE definition is exposed here:
```go
type CTE struct {
    Name         string
    Query        string        // raw SQL string of the CTE body
    ParsedQuery  *ParsedQuery  // fully parsed CTE body
    Materialized string        // "", "MATERIALIZED", or "NOT MATERIALIZED"
}
```

### CTE tables in `parsedSQL.Tables`
CTE aliases appear alongside real tables in the outer query's `Tables` slice, but with `Type: "cte"` instead of `Type: "base"`:
```json
{"Name": "city_directors", "Type": "cte", ...}
{"Name": "directors",      "Type": "base", ...}
```

### CTE parameters
Parameters inside a CTE's WHERE clause:
- Appear in top-level `parsedSQL.Parameters` (same as regular params)
- Have **no** matching filter `ColumnUsage` in the outer query's `parsedSQL.ColumnUsage`
- Are found in `cte.ParsedQuery.ColumnUsage` with `UsageType == "filter"` and a `Context` field like `"users.active = $1"`

Multiple params in one CTE produce separate filter ColumnUsage entries, one per condition:
```json
{"Column":"active","UsageType":"filter","Context":"users.active = $1"}
{"Column":"age",   "UsageType":"filter","Context":"users.age > $2"}
```

### CTE column output names
A CTE output column's name is: the `Alias` if set, otherwise the last segment of `Expression` after the final `.`:
```
{Expression: "users.id",    Alias: ""}        → output name: "id"
{Expression: "users.name",  Alias: "user_nm"} → output name: "user_nm"
```

---

## Implementation Details

### resolve_projections.go — CTE column type resolution

`resolveProjections`, `resolveColumnGoType`, `resolveAggregationType`, and `resolveSimpleColumnType` all gained a `ctes []postgresparser.CTE` parameter.

In `resolveColumnGoType`, after resolving `tableName`, before looking up in `tablesCol`:
```go
for _, cte := range ctes {
    if cte.Name == tableName {
        for _, cteCol := range cte.ParsedQuery.Columns {
            if cteOutputColName(cteCol) == colName {
                return c.resolveColumnGoType(cteCol, cte.ParsedQuery.Tables, cte.ParsedQuery.CTEs)
            }
        }
        return "", fmt.Errorf("column %q not found in CTE %q", colName, tableName)
    }
}
```

`cteOutputColName` helper:
```go
func cteOutputColName(col postgresparser.SelectColumn) string {
    if col.Alias != "" {
        return col.Alias
    }
    expr := col.Expression
    if i := strings.LastIndex(expr, "."); i >= 0 {
        return expr[i+1:]
    }
    return expr
}
```

The table-existence check (`if !ok { return "table not found" }`) was **removed** — pglite's PREPARE-based validation already catches invalid table references, making the check redundant.

### resolve_params.go — CTE parameter resolution

Added `buildCTEParamMap` (mirrors `buildSubqueryParamMap`):
```go
func buildCTEParamMap(ctes []postgresparser.CTE) map[int]subqueryParamInfo {
    // iterates CTE filter ColumnUsage entries, extracts $N position from Context
    // recurses into nested CTEs (cte.ParsedQuery.CTEs)
}
```

In `resolveParams`, the lookup order is now:
1. `subqMap` (IN/NOT IN subquery params)
2. `cteMap` (CTE WHERE params) ← new
3. Top-level filter ColumnUsage
4. LIMIT/OFFSET overflow

### generator.go — call site update
```go
// Before:
c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
// After:
c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
```

---

## Schema Loading (pglite-based)

The old approach (parsing `db/schema/*.sql` via `postgresparser.ParseSQLStrict`) was replaced with:

1. Run all `db/migrations/*.up.sql` files against an in-memory pglite instance
2. Query `information_schema.columns` with `WHERE table_schema = 'pg_catalog'` — this is intentional: pglite exposes user-created tables under `pg_catalog`, not `public`. The filter is correct for pglite and should not be changed to `'public'`.
3. Populate `c.tablesCol` from the result

**pglite setup:**
```go
db, _ := sql.Open("pglite", ":memory:")
db.SetMaxOpenConns(1)
db.SetConnMaxIdleTime(0)
db.SetConnMaxLifetime(0)
```

### Query Validation
Every query is validated via PREPARE before code generation:
```go
c.db.ExecContext(ctx, fmt.Sprintf("PREPARE pgen_test_%s as %s", query.name, query.sql))
```
This catches invalid table/column references at generation time, which is why redundant schema-level existence checks were removed from type resolution.

---

## Named Parameters (@name syntax)

`convertNamedParams()` in `generator.go` converts `@name` style to `$N` before parsing:
- Detects `@name` vs `$N` — errors if mixed
- Same `@name` used twice maps to the same `$N`
- Returns the ordered list of original names (used to override resolved param names in generated function signatures)

---

## Test Infrastructure

- `TestMain` in `testhelpers_test.go` sets up a shared pglite DB with `users` and `posts` tables
- `testSharedCli` is the shared `*cli` used by most tests
- `testCliWithEmptyDB` creates a fresh isolated pglite DB for tests that need clean state
- No file I/O in tests — all SQL is hardcoded strings

### Test schemas
```sql
-- users
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY, name TEXT NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL, age SMALLINT,
    status SMALLINT NOT NULL, role_id INTEGER NOT NULL,
    login_count INTEGER, org_id BIGINT NOT NULL,
    referrer_id BIGINT, active BOOLEAN DEFAULT TRUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);
-- posts
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY, title TEXT NOT NULL,
    body TEXT, user_id BIGINT NOT NULL,
    published BOOLEAN NOT NULL DEFAULT FALSE
);
```

---

## Known Issues / Gaps

- `load_schema.go` queries `WHERE table_schema = 'pg_catalog'` — **this is correct and intentional**. pglite exposes user-created tables under `pg_catalog`, not `public`. Do not change this to `'public'`.
- `db/db.go` imports `github.com/jackc/pgx/v5` but that package is not in `go.mod` (the generator itself uses pglite). The generated `db/` package is intended to be used in a consumer project that has pgx as a dependency.
- Duplicate query names within a file are not detected — **this is intentional**, not a bug. The behaviour is acceptable and there are no plans to add detection.
- No params struct generated for 3+ parameters (noted in original CLAUDE.md).
