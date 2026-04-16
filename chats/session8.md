# Phase 8 Session — Knowledge Base

## Changes This Session

### `json` / `jsonb` / `numeric` / `decimal` type support

Added Go type mappings for PostgreSQL JSON and numeric types, which previously fell through to the `default` case in `pgTypeToGoType` and silently returned `"string"`.

---

### `resolve_query.go` — new cases in `pgTypeToGoType`

Added immediately after the `"timestamptz"` case, before `default`:

```go
case "json":
    if nullable {
        return "pgtype.JSON"
    }
    return "[]byte"
case "jsonb":
    if nullable {
        return "pgtype.JSONB"
    }
    return "[]byte"
case "numeric", "decimal":
    return "pgtype.Numeric"
```

**Type mapping rationale:**
- Non-null json/jsonb → `[]byte` (idiomatic; works directly with `json.Unmarshal`, no custom scanner)
- Nullable json → `pgtype.JSON`, nullable jsonb → `pgtype.JSONB`
- `json` and `jsonb` are kept as **separate cases** (not merged) because their nullable types differ (`pgtype.JSON` vs `pgtype.JSONB`), matching the `"timestamp"` / `"timestamptz"` pattern
- numeric/decimal → `pgtype.Numeric` for both nullable and non-null (no native Go decimal type — same pattern as `date` → always `pgtype.Date`)
- `"decimal"` included as an alias for `"numeric"` for robustness; `information_schema.data_type` always returns `"numeric"` in practice, but the case must exist
- No new imports needed — `pgtype` is already unconditionally imported in non-std mode

---

### `db/migrations/014_add_json_numeric.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN metadata     jsonb,
    ADD COLUMN config       json NOT NULL DEFAULT '{}',
    ADD COLUMN price        numeric NOT NULL DEFAULT 0,
    ADD COLUMN rating       numeric;
```

Covers all four cases in the type matrix:

| Column | Type | Nullable | Go type |
|---|---|---|---|
| `movies.metadata` | `jsonb` | YES | `pgtype.JSONB` |
| `movies.config` | `json` | NO | `[]byte` |
| `movies.rating` | `numeric` | YES | `pgtype.Numeric` |
| `movies.price` | `numeric` | NO | `pgtype.Numeric` |

Follows the same pattern as migration 013 (ALTER TABLE on `movies`).

---

### `resolve_projections_test.go` — new tests

**Schema-backed tests** (use `testSharedCli` against real migrations):
- `TestResolveProjections_JsonbNullableColumn` — `movies.metadata` → `pgtype.JSONB`
- `TestResolveProjections_JsonNotNullColumn` — `movies.config` → `[]byte`
- `TestResolveProjections_NumericNullableColumn` — `movies.rating` → `pgtype.Numeric`
- `TestResolveProjections_NumericNotNullColumn` — `movies.price` → `pgtype.Numeric`

**Direct unit tests** (call `pgTypeToGoType` directly, adjacent to `TestPgTypeToGoType_Timestamp*`):
- `TestPgTypeToGoType_JsonbNotNull` — `("jsonb", false)` → `"[]byte"`
- `TestPgTypeToGoType_JsonbNullable` — `("jsonb", true)` → `"pgtype.JSONB"`
- `TestPgTypeToGoType_JsonNotNull` — `("json", false)` → `"[]byte"`
- `TestPgTypeToGoType_JsonNullable` — `("json", true)` → `"pgtype.JSON"`
- `TestPgTypeToGoType_NumericNotNull` — `("numeric", false)` → `"pgtype.Numeric"`
- `TestPgTypeToGoType_NumericNullable` — `("numeric", true)` → `"pgtype.Numeric"`
- `TestPgTypeToGoType_DecimalAlias` — `("decimal", false/true)` → `"pgtype.Numeric"` (pure unit test; `information_schema` always returns `"numeric"` for DECIMAL columns but the switch case must exist)

---

### `resolve_params_test.go` — new tests

Projection and param resolution are separate code paths (`resolveColumnGoType` vs `resolveParamFromUsage` + `resolveInsertParams`), so both paths need coverage.

**Schema-backed WHERE-param tests:**
- `TestResolveParams_JsonbParam` — `movies.metadata = $1` → `pgtype.JSONB`
- `TestResolveParams_JsonNotNullParam` — `movies.config = $1` → `[]byte`
- `TestResolveParams_NumericNullableParam` — `movies.rating = $1` → `pgtype.Numeric`
- `TestResolveParams_NumericNotNullParam` — `movies.price = $1` → `pgtype.Numeric`

**INSERT test** (exercises `resolveInsertParams` code path):
- `TestResolveParams_InsertJsonbAndNumeric` — `INSERT INTO movies (name, metadata, config, rating, price) VALUES ($1, $2, $3, $4, $5)` → `["string", "pgtype.JSONB", "[]byte", "pgtype.Numeric", "pgtype.Numeric"]`
