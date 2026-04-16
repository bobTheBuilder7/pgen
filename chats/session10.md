# Phase 10 Session — Knowledge Base

## Changes This Session

### `bytea` type support

Added Go type mappings for PostgreSQL `bytea` columns, which previously fell through to the `default` case in `pgTypeToGoType` and silently returned `"string"`.

---

### `resolve_query.go` — new case in `pgTypeToGoType`

Added immediately after the `"jsonb"` case, before `"numeric"/"decimal"`:

```go
case "bytea":
    if nullable {
        return "pgtype.Bytea"
    }
    return "[]byte"
```

**Type mapping rationale:**
- Non-null → `[]byte` (idiomatic; works directly with binary data, same pattern as non-null json/jsonb)
- Nullable → `pgtype.Bytea`
- No new imports needed — `pgtype` is already unconditionally imported in non-std mode

---

### `db/migrations/016_add_bytea.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN checksum bytea NOT NULL DEFAULT '',
    ADD COLUMN thumbnail bytea;
```

| Column | Type | Nullable | Go type |
|---|---|---|---|
| `movies.checksum` | `bytea` | NO | `[]byte` |
| `movies.thumbnail` | `bytea` | YES | `pgtype.Bytea` |

Follows the same pattern as migrations 013–015 (ALTER TABLE on `movies`).

---

### `resolve_projections_test.go` — new tests

**Unit tests** (call `pgTypeToGoType` directly):
- `TestPgTypeToGoType_ByteaNotNull` — `("bytea", false)` → `"[]byte"`
- `TestPgTypeToGoType_ByteaNullable` — `("bytea", true)` → `"pgtype.Bytea"`

**Schema-backed tests** (use `testSharedCli`):
- `TestResolveProjections_ByteaNotNullColumn` — `movies.checksum` → `[]byte`
- `TestResolveProjections_ByteaNullableColumn` — `movies.thumbnail` → `pgtype.Bytea`

---

### `resolve_params_test.go` — new tests

**WHERE param tests:**
- `TestResolveParams_ByteaNotNullParam` — `WHERE checksum = $1` → `[]byte`
- `TestResolveParams_ByteaNullableParam` — `WHERE thumbnail = $1` → `pgtype.Bytea`

**INSERT test:**
- `TestResolveParams_InsertBytea` — `INSERT INTO movies (name, checksum, thumbnail) VALUES ($1, $2, $3)` → `["string", "[]byte", "pgtype.Bytea"]`
