# Session 13 — Replace gopglite with testcontainers-go

## Summary

Replaced `github.com/bobTheBuilder7/gopglite` (in-memory pglite) with `github.com/testcontainers/testcontainers-go` + `github.com/jackc/pgx/v5/pgxpool` throughout the codebase.

## Changes

### `run.go`
- `cli.db` field changed from `*sql.DB` to `*pgxpool.Pool`
- Removed `database/sql` and `_ gopglite` imports
- `run()` now spins up a `postgres:18` testcontainer, gets endpoint via `PortEndpoint`, creates a `pgxpool.Pool`, pings, then tears down via `defer postgresC.Terminate(ctx)`
- Default credentials in `run()`: `POSTGRES_PASSWORD=password`, user `postgres`, db `postgres`

### `migrations.go`
- `c.db.ExecContext(ctx, ...)` → `c.db.Exec(ctx, ...)`

### `load_schema.go`
- `c.db.QueryContext(ctx, ...)` → `c.db.Query(ctx, ...)`
- **Critical fix**: `WHERE table_schema = 'pg_catalog'` → `WHERE table_schema = 'public'`
  - pglite stored user tables under `pg_catalog`; real Postgres puts them in `public`

### `prepare_query.go`
- `c.db.ExecContext(ctx, ...)` → `c.db.Exec(ctx, ...)`

### `testhelpers_test.go`
- Removed `database/sql` and `_ gopglite` imports
- `TestMain` now uses a `testMain(m) int` helper (so `defer` runs before `os.Exit`) that starts a testcontainer, creates a pool, runs migrations, loads schema
- `testCliWithEmptyDB` starts a per-test testcontainer using `t.Context()` and `testcontainers.CleanupContainer(t, postgresC)`
- Both use credentials `test/test`, db `test`

### `migrations_test.go`
- `c.db.QueryRowContext(...)` → `c.db.QueryRow(...)`
- `c.db.QueryContext(...)` → `c.db.Query(...)`

### `go.mod`
- Added direct: `github.com/testcontainers/testcontainers-go`, `github.com/jackc/pgx/v5`
- Removed: `github.com/bobTheBuilder7/gopglite`, `github.com/tetratelabs/wazero`

## Key notes
- `wait.ForListeningPort("5432/tcp")` was removed from wait strategies — `postgres:18` has no shell so testcontainers can't exec into it; `wait.ForLog(...).WithOccurrence(2)` is sufficient
- `testCliWithEmptyDB` creates one container per test call (isolated, slow but correct)
- `TestMain` reuses one container for all schema-backed tests via `testSharedCli`
