# Phase 5 Session — Knowledge Base

## Changes This Session

### Test migration to real db/migrations schema

Replaced all hardcoded `users`/`posts` test tables with the actual project schema from `db/migrations/`.

---

### `testhelpers_test.go` — run migrations from directory

Removed the `usersSchemaSQL`/`postsSchemaSQL` string constants and their inline `runMigration` calls. `TestMain` now applies all `*.up.sql` files from `db/migrations/` using the same loop as `run.go`:

```go
files, err := os.ReadDir(filepath.Join(dbDirectory, migrationsDirectory))
// sort, open each file, runMigration, close
```

Added `path/filepath` import; `strings` retained for `HasSuffix`.

---

### `resolve_query.go` — added `date` type support

`pgTypeToGoType` had no case for `date`, which fell through to `"string"`. Added:

```go
case "date":
    return "pgtype.Date"  // both nullable and not-null use pgtype.Date (no native Go date type)
```

This was a pre-existing gap exposed by the schema migration (the real migrations have `date` columns; the old `users`/`posts` schema did not).

---

### Test files rewritten — table substitution

**Primary tables**:
- `movies` replaces `users` (bigserial id, text name, date when_released, integer box_office)
- `trailers` replaces `posts` (serial id, bigint movie_id FK, text url, date when_released)
- `cities` used for smallint types (state_id smallint NULL → pgtype.Int2, country_id smallint NN → int16)
- `actors` referenced in plan but `movies`/`trailers`/`cities` were sufficient

**Column type mapping** (old → new):

| Old | New | Go type |
|-----|-----|---------|
| users.id (bigserial NN) | movies.id | int64 |
| users.name (text NN) | movies.name | string |
| users.email (varchar NN) | trailers.url | string |
| users.status (smallint NN) | cities.country_id | int16 |
| users.role_id (integer NN) | trailers.id | int32 |
| users.age (smallint NULL) | cities.state_id | pgtype.Int2 |
| users.login_count (int NULL) | movies.box_office | pgtype.Int4 |
| users.referrer_id (bigint NULL) | trailers.movie_id | pgtype.Int8 |
| users.active (bool NULL) | movies.when_released | pgtype.Date |
| users.verified (bool NN) | cities.country_id or movies.name | int16 / string |
| posts.body (text NULL) | trailers.when_released | pgtype.Date |
| posts.published (bool NN) | repurposed | — |

**Boolean types**: The real schema has no boolean columns. Tests that checked `bool`/`pgtype.Bool` were repurposed to test `pgtype.Date` (nullable date) instead — same nullability behavior, different concrete type.

**JOIN pair**: `movies m LEFT/RIGHT/INNER JOIN trailers t ON t.movie_id = m.id` replaces `users u JOIN posts p ON u.id = p.user_id`.

**CTE tests**: Updated to use `movies`/`trailers` instead of `users`/`posts`.

---

### Files changed

| File | Change |
|------|--------|
| `testhelpers_test.go` | Runs `db/migrations/*.up.sql` directory instead of hardcoded SQL |
| `resolve_query.go` | Added `date → pgtype.Date` in `pgTypeToGoType` |
| `generator_test.go` | All SQL rewritten to use `movies`/`trailers` |
| `parse_query_file_test.go` | SQL fixtures updated to `movies` |
| `resolve_params_test.go` | Full rewrite (~300 lines) — movies/trailers/cities |
| `resolve_projections_test.go` | Full rewrite (~650 lines) — movies/trailers/cities |
