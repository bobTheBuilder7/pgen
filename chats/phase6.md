# Phase 6 Session — Knowledge Base

## Changes This Session

### Directory paths moved to CLI flags

Replaced the three package-level constants (`dbDirectory`, `queriesDirectory`, `migrationsDirectory`) with CLI flags stored on the `cli` struct.

---

### `main.go` — constants removed, flags added

Removed:
```go
const dbDirectory = "db"
const queriesDirectory = "queries"
const migrationsDirectory = "migrations"
```

Added three flags (defaults match the old constant values):
```go
dbDir         := flag.String("db",         "db",          "path to the db directory")
queriesDir    := flag.String("queries",    "queries",     "subdirectory under db containing query files")
migrationsDir := flag.String("migrations", "migrations",  "subdirectory under db containing migration files")
```

Passed to `run()`:
```go
err := run(ctx, *std, *debug, *dbDir, *queriesDir, *migrationsDir)
```

`sqlConstSuffix` was not touched — it is a code-generation constant, not a filesystem path.

---

### `run.go` — fields added to `cli`, `run()` signature updated

New fields on `cli`:
```go
type cli struct {
    tablesCol           syncmap.Map[string, []dbColumn]
    db                  *sql.DB
    dbDirectory         string
    queriesDirectory    string
    migrationsDirectory string
}
```

`run()` signature:
```go
func run(ctx context.Context, std bool, debug bool, dbDirectory string, queriesDirectory string, migrationsDirectory string) error
```

`c` initialised with the new fields:
```go
c := &cli{
    db:                  db,
    dbDirectory:         dbDirectory,
    queriesDirectory:    queriesDirectory,
    migrationsDirectory: migrationsDirectory,
}
```

All path construction updated to use `c.*`:
- `filepath.Join(dbDirectory, migrationsDirectory, ...)` → `filepath.Join(c.dbDirectory, c.migrationsDirectory, ...)`
- `filepath.Join(dbDirectory, queriesDirectory, ...)` → `filepath.Join(c.dbDirectory, c.queriesDirectory, ...)`
- Hardcoded `"./db/db.go"` → `filepath.Join(c.dbDirectory, "db.go")`

---

### `testhelpers_test.go` — string literals replace constants

The constants are gone, so the two references in `TestMain` became inline string literals:
```go
// before
filepath.Join(dbDirectory, migrationsDirectory)
// after
filepath.Join("db", "migrations")
```
