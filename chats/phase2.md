# Phase 2 Session — Knowledge Base

## Changes This Session

### `-d` (debug) flag

Added a `-d` flag to the CLI. When passed, generated files are printed to stdout instead of being written to disk.

**Files changed:**
- `main.go` — added `debug := flag.Bool("d", false, ...)`, passed `*debug` to `run()`
- `run.go` — `run()` signature changed from `run(ctx, std bool)` to `run(ctx, std bool, debug bool)`

**Behaviour in `run.go`:**
- For each query `.sql` file: if `debug`, write generated Go to `os.Stdout`; otherwise `os.Create(...)` as before
- For `db/db.go`: if `debug`, write to `os.Stdout`; otherwise `os.Create("./db/db.go")` as before

**Usage:**
```
go run . -d        # print all generated files to stdout, no disk writes
go run .           # normal mode, writes db/*.go and db/db.go
go run . -std      # normal mode with database/sql base file
go run . -d -std   # debug mode with database/sql base file
```
