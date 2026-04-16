# Phase 9 Session — Knowledge Base

## Changes This Session

### PostgreSQL array type support

Added Go type mappings for PostgreSQL array columns (e.g. `text[]`, `integer[]`), which previously fell through to the `default` case in `pgTypeToGoType` and silently returned `"string"`.

---

### `load_schema.go` — fetch `udt_name` for array element detection

Added `udt_name` to the `information_schema.columns` SELECT. PostgreSQL reports array columns as `data_type = 'ARRAY'` with `udt_name = '_text'`, `'_int4'`, etc. (underscore-prefixed element type name). When an array is detected, the element type is encoded directly into `dbColumn.Type` as `"text[]"`, `"int4[]"`, etc.:

```go
colType := dataType
if strings.EqualFold(dataType, "ARRAY") {
    colType = strings.TrimPrefix(udtName, "_") + "[]"
}
```

No `dbColumn` struct changes needed — the `[]`-suffixed type string flows through the existing pipeline unchanged.

Added `"strings"` import (only new import across all changed files).

---

### `resolve_query.go` — array branch in `pgTypeToGoType`

Added at the top of `pgTypeToGoType`, before the `switch`:

```go
if strings.HasSuffix(pgType, "[]") {
    elemGoType := pgTypeToGoType(pgType[:len(pgType)-2], false)
    if nullable {
        return "pgtype.Array[" + elemGoType + "]"
    }
    return "[]" + elemGoType
}
```

The recursion always passes `nullable=false` for the element type — only the outer container differs between null and non-null arrays.

**Type mapping:**

| PG type | Nullable | Go type |
|---|---|---|
| `text[]` | NO | `[]string` |
| `text[]` | YES | `pgtype.Array[string]` |
| `integer[]` | NO | `[]int32` |
| `integer[]` | YES | `pgtype.Array[int32]` |
| `bigint[]` | NO | `[]int64` |
| `bigint[]` | YES | `pgtype.Array[int64]` |
| `boolean[]` | NO | `[]bool` |
| `boolean[]` | YES | `pgtype.Array[bool]` |

The recursion reuses all existing element-type cases — no additional switch cases needed.

---

### `db/migrations/015_add_arrays.up.sql` — new migration

```sql
ALTER TABLE movies
    ADD COLUMN tags        text[]    NOT NULL DEFAULT '{}',
    ADD COLUMN scores      integer[],
    ADD COLUMN flags       boolean[] NOT NULL DEFAULT '{}',
    ADD COLUMN related_ids bigint[];
```

Covers: the four most common element types, both NOT NULL (`[]T`) and NULL (`pgtype.Array[T]`) output forms.

---

### Tests added

**Unit tests** (`resolve_projections_test.go`, call `pgTypeToGoType` directly):
- `TestPgTypeToGoType_TextArrayNotNull` / `_Nullable`
- `TestPgTypeToGoType_IntArrayNotNull` / `_Nullable`
- `TestPgTypeToGoType_BigintArrayNotNull` / `_Nullable`
- `TestPgTypeToGoType_BoolArrayNotNull` / `_Nullable`

**Schema-backed projection tests** (`resolve_projections_test.go`, use `testSharedCli`):
- `TestResolveProjections_TextArrayNotNullColumn` — `movies.tags` → `[]string`
- `TestResolveProjections_IntArrayNullableColumn` — `movies.scores` → `pgtype.Array[int32]`
- `TestResolveProjections_BoolArrayNotNullColumn` — `movies.flags` → `[]bool`
- `TestResolveProjections_BigintArrayNullableColumn` — `movies.related_ids` → `pgtype.Array[int64]`
- `TestResolveProjections_MixedArrayAndScalarColumns` — scalar + array columns in one query

**Schema-backed param tests** (`resolve_params_test.go`):
- `TestResolveParams_InsertTextArrayParam` — INSERT `tags` → `[]string`
- `TestResolveParams_InsertNullableIntArrayParam` — INSERT `scores` → `pgtype.Array[int32]`
- `TestResolveParams_WhereTextArrayParam` — WHERE `tags = $1` → `[]string`
- `TestResolveParams_WhereNullableArrayParam` — WHERE `scores = $1` → `pgtype.Array[int32]`

---

## Known Limitations

- **Multidimensional arrays**: `information_schema` always returns `udt_name = '_text'` regardless of array dimensions — treated as 1D. `text[][]` generates the same type as `text[]`.
- **Arrays of uncommon element types** (e.g. `numeric[]`, `uuid[]`, `date[]`): the recursive call hits the existing switch case and resolves correctly — these are supported but not explicitly tested.
