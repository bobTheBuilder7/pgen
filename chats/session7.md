# Phase 7 Session — Knowledge Base

## Changes This Session

### `timestamp` / `timestamptz` type support

Added Go type mappings for PostgreSQL timestamp types, which previously fell through to the `default` case in `pgTypeToGoType` and silently returned `"string"`.

---

### `resolve_query.go` — new cases in `pgTypeToGoType`

Added immediately after the `"date"` case:

```go
case "timestamp", "timestamp without time zone":
    if nullable {
        return "pgtype.Timestamp"
    }
    return "time.Time"
case "timestamptz", "timestamp with time zone":
    if nullable {
        return "pgtype.Timestamptz"
    }
    return "time.Time"
```

Both the short alias (`"timestamp"`) and the long `information_schema` form (`"timestamp without time zone"`) are handled because pglite returns the long form from `information_schema.columns`.

**Type mapping rationale:**
- Non-null → `time.Time` (native Go type, same pattern as `string` for text, `int64` for bigint)
- Nullable → `pgtype.Timestamp` / `pgtype.Timestamptz`
- Contrast with `date`, which uses `pgtype.Date` for both — there is no native Go date type, but `time.Time` is the native Go timestamp type.

---

### `generator.go` — conditional `"time"` import

Added `usesTimeType` helper and conditional import in `generateCode`:

```go
if usesTimeType(queries) {
    f.AddBlock(gen.Import("", "time"))
}
```

```go
func usesTimeType(queries []resolvedQuery) bool {
    for _, rq := range queries {
        for _, f := range rq.rowFields {
            if f.Type == "time.Time" {
                return true
            }
        }
        for _, pt := range rq.paramTypes {
            if pt == "time.Time" {
                return true
            }
        }
    }
    return false
}
```

The `"time"` import is only emitted when at least one row field or param type is `time.Time`, avoiding unused import errors in generated files with no timestamp columns.

---

### `db/migrations/013_add_timestamps.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMPTZ;
```

- `movies.created_at` — `timestamptz NOT NULL` → `time.Time`
- `movies.updated_at` — `timestamptz NULL` → `pgtype.Timestamptz`

---

### `resolve_projections_test.go` — 4 new tests

**Schema-backed tests** (use `testSharedCli` against real migrations):
- `TestResolveProjections_TimestamptzNotNullColumn` — `movies.created_at` → `time.Time`
- `TestResolveProjections_TimestamptzNullableColumn` — `movies.updated_at` → `pgtype.Timestamptz`

**Direct unit tests** (call `pgTypeToGoType` directly, no DB):
- `TestPgTypeToGoType_TimestampNotNull` — both `"timestamp"` and `"timestamp without time zone"` → `"time.Time"`
- `TestPgTypeToGoType_TimestampNullable` — both → `"pgtype.Timestamp"`
