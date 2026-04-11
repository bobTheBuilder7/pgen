# Phase 3 Session — Knowledge Base

## Changes This Session

### RETURNING generator tests

Added 11 tests to `generator_test.go` covering the `generateExec` RETURNING path, which was previously untested at the code-generation level (`resolveReturning` had unit tests but nothing exercised it through `generateCode`).

**Happy path — method name checks:**
- `TestGenerateCode_InsertReturningOneUsesQueryRow` — INSERT RETURNING :one → `QueryRow` (pgx mode)
- `TestGenerateCode_InsertReturningManyUsesQuery` — INSERT RETURNING :many → `.Query` (pgx mode)
- `TestGenerateCode_StdModeInsertReturningOneUsesQueryRowContext` — INSERT RETURNING :one → `QueryRowContext` (std mode)
- `TestGenerateCode_StdModeInsertReturningManyUsesQueryContext` — INSERT RETURNING :many → `QueryContext` (std mode)
- `TestGenerateCode_UpdateReturningOneSucceeds` — UPDATE RETURNING :one succeeds
- `TestGenerateCode_DeleteReturningOneSucceeds` — DELETE RETURNING :one succeeds

**Output shape checks:**
- `TestGenerateCode_InsertReturningGeneratesRowStruct` — `CreateUserRow` struct is emitted
- `TestGenerateCode_InsertReturningGeneratesSQLConst` — `CreateUserSQL` const is emitted
- `TestGenerateCode_InsertReturningGeneratesScanCall` — `.Scan(` call is emitted

**Error path — wrong type + RETURNING combinations:**
- `TestGenerateCode_InsertReturningWithExecTypeReturnsError` — INSERT RETURNING :exec → error mentioning `RETURNING`
- `TestGenerateCode_InsertReturningWithExecResultTypeReturnsError` — INSERT RETURNING :execresult → same
- `TestGenerateCode_InsertWithoutReturningWithOneTypeReturnsError` — INSERT no RETURNING :one → error mentioning `one`
- `TestGenerateCode_InsertWithoutReturningWithManyTypeReturnsError` — INSERT no RETURNING :many → error mentioning `many`

---

## Known Gaps (identified, not yet addressed)

### Not covered in generator_test.go
1. **RETURNING field type projection** — no test verifies that a nullable column produces `pgtype.Int2` etc. in generated output. `resolveReturning` is unit-tested in `resolve_projections_test.go` but the wiring through `generateExec` is unchecked at the type level.

### Not covered anywhere
2. **float/uuid types** — `pgTypeToGoType` handles `float4`, `float8`, `real`, `double precision`, `uuid` but the shared test schema has none of those columns. Never exercised.
3. **CTE end-to-end in generateCode** — CTE resolution is unit-tested at `resolveProjections`/`resolveParams` level but no `TestGenerateCode_CTE*` test runs a full CTE query through `generateCode`.
4. **String literal projection** — `SELECT 'constant' AS label FROM users` → the `'foo'` → `string` branch in `resolveColumnGoType` has no test.
5. **Unsupported aggregation error** — the default error branch in `resolveAggregationType` (e.g. `STDDEV(...)`) is untested.
