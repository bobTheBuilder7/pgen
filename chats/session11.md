# Phase 11 Session — Knowledge Base

## Changes This Session

### `time` / `timetz` / `interval` type support + `varchar` verification

Added Go type mappings for PostgreSQL time-of-day and interval types, which previously fell through to the `default` case in `pgTypeToGoType` and silently returned `"string"`. Also verified `varchar`/`character varying` (already in the switch) works correctly through the real pglite schema by adding a migration and schema-backed tests.

---

### `resolve_query.go` — new cases in `pgTypeToGoType`

Added immediately after the `"timestamptz"` case, before `"json"`:

```go
case "time", "time without time zone":
    return "pgtype.Time"
case "timetz", "time with time zone":
    return "pgtype.Time"
case "interval":
    return "pgtype.Interval"
```

**Type mapping rationale:**
- No native Go time-of-day type → always use `pgtype.Time` (same pattern as `date`/`numeric` — no nullable/not-null distinction)
- `"time without time zone"` and `"time with time zone"` are the long forms returned by pglite's `information_schema.data_type`; short aliases included for robustness
- `pgtype.Time` is the same type for timetz — pgx v5 has no separate timetz wrapper
- `pgtype.Interval` for both nullable and not-null interval
- No new imports needed — `pgtype` is already unconditionally imported in non-std mode

**varchar was already handled** at line 187 (`case "text", "varchar", "character varying", "char", "character":`). The migration + schema-backed tests confirmed pglite returns `"character varying"` from `information_schema.data_type` for VARCHAR columns.

---

### `db/migrations/017_add_varchar.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN slug        varchar(100) NOT NULL DEFAULT '',
    ADD COLUMN description varchar(500);
```

| Column | Type | Nullable | Go type |
|---|---|---|---|
| `movies.slug` | `character varying` | NO | `string` |
| `movies.description` | `character varying` | YES | `pgtype.Text` |

---

### `db/migrations/018_add_time_interval.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN show_time   time     NOT NULL DEFAULT '00:00:00',
    ADD COLUMN show_timetz timetz,
    ADD COLUMN duration    interval NOT NULL DEFAULT '0',
    ADD COLUMN break_time  interval;
```

| Column | Type | Nullable | Go type |
|---|---|---|---|
| `movies.show_time` | `time without time zone` | NO | `pgtype.Time` |
| `movies.show_timetz` | `time with time zone` | YES | `pgtype.Time` |
| `movies.duration` | `interval` | NO | `pgtype.Interval` |
| `movies.break_time` | `interval` | YES | `pgtype.Interval` |

---

### Tests added

**Unit tests** (`resolve_projections_test.go`, call `pgTypeToGoType` directly):

*varchar:*
- `TestPgTypeToGoType_VarcharNotNull` — `("character varying", false)` → `"string"`
- `TestPgTypeToGoType_VarcharNullable` — `("character varying", true)` → `"pgtype.Text"`

*time:*
- `TestPgTypeToGoType_TimeNotNull` / `_TimeNullable` — both → `"pgtype.Time"`
- `TestPgTypeToGoType_TimeLongForm` — `("time without time zone", false)` → `"pgtype.Time"`
- `TestPgTypeToGoType_TimetzNotNull` / `_TimetzLongForm` — both → `"pgtype.Time"`

*interval:*
- `TestPgTypeToGoType_IntervalNotNull` / `_IntervalNullable` — both → `"pgtype.Interval"`

**Schema-backed projection tests** (`resolve_projections_test.go`, use `testSharedCli`):
- `TestResolveProjections_VarcharNotNullColumn` — `movies.slug` → `string`
- `TestResolveProjections_VarcharNullableColumn` — `movies.description` → `pgtype.Text`
- `TestResolveProjections_TimeNotNullColumn` — `movies.show_time` → `pgtype.Time`
- `TestResolveProjections_TimetzNullableColumn` — `movies.show_timetz` → `pgtype.Time`
- `TestResolveProjections_IntervalNotNullColumn` — `movies.duration` → `pgtype.Interval`
- `TestResolveProjections_IntervalNullableColumn` — `movies.break_time` → `pgtype.Interval`

**Schema-backed param tests** (`resolve_params_test.go`):

WHERE params:
- `TestResolveParams_VarcharNotNullParam` / `_VarcharNullableParam`
- `TestResolveParams_TimeNotNullParam` / `_TimetzNullableParam`
- `TestResolveParams_IntervalNotNullParam` / `_IntervalNullableParam`

INSERT tests:
- `TestResolveParams_InsertVarchar` — `INSERT INTO movies (name, slug, description)` → `["string", "string", "pgtype.Text"]`
- `TestResolveParams_InsertTimeInterval` — `INSERT INTO movies (name, show_time, show_timetz, duration, break_time)` → `["string", "pgtype.Time", "pgtype.Time", "pgtype.Interval", "pgtype.Interval"]`
