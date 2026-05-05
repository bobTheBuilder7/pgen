# Phase 14 Session — Knowledge Base

## Changes This Session

### `EXISTS` support

Added `SELECT EXISTS(...)` projection support and verified `WHERE EXISTS (...)` param resolution already worked.

---

### `resolve_projections.go` — new case in `resolveAggregationType`

Added immediately before the `default` case:

```go
case "EXISTS":
    return "bool", false, nil
```

**Type mapping rationale:**
- EXISTS always evaluates to a boolean, never null → always `bool`
- Same pattern as `COUNT` (always non-null, fixed Go type regardless of nullable)
- No new imports needed

Updated the `default` error message to include EXISTS:
```go
return "", false, fmt.Errorf("unsupported aggregation function %q: only COUNT, SUM, AVG, MIN, MAX, COALESCE, and EXISTS are supported", fn)
```

Updated the function comment to mention EXISTS.

---

### How EXISTS is handled

**`SELECT EXISTS(...) AS alias`** — the expression `EXISTS(...)` contains `(` and doesn't start with `(`, so it:
1. Hits the alias guard in `resolveProjections` — **alias is required**, same as COUNT/SUM/etc.
2. Hits the aggregation branch in `resolveColumnGoType`
3. `aggregationRegex` matches with `fn = "EXISTS"`, new case returns `"bool", false`
4. `pgTypeToGoType("bool", false)` → `"bool"`

**Params inside `SELECT EXISTS(...)`** — the EXISTS subquery lands in `parsedSQL.Subqueries`, so `buildSubqueryParamMap` already resolved these correctly with no changes needed.

**`WHERE EXISTS (...)`** — already worked with zero changes. The parser produces no phantom top-level filter ColumnUsage for WHERE EXISTS (unlike IN subqueries), and the `$N` params are resolved through `buildSubqueryParamMap`.

---

### Tests added

**`resolve_projections_test.go`:**
- `TestResolveProjections_SelectExistsReturnsBool` — `EXISTS(...) AS exists_trailers` → `bool` field named `ExistsTrailers`
- `TestResolveProjections_SelectExistsRequiresAlias` — `EXISTS(...)` without alias → error containing `"alias"`

**`resolve_params_test.go`:**
- `TestResolveParams_ExistsSubqueryParam` — `$1` inside `EXISTS(SELECT 1 FROM trailers WHERE trailers.movie_id = $1)` resolves to `"movie_id"` / `"pgtype.Int8"`

**`generator_test.go`:**
- `TestGenerateCode_SelectExistsGeneratesBoolField` — end-to-end: generated row struct contains `ExistsTrailers bool`
- `TestGenerateCode_SelectExistsWithoutAliasReturnsError` — end-to-end alias enforcement
