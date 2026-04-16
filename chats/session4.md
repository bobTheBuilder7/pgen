# Phase 4 Session — Knowledge Base

## Changes This Session

### Generator refactor — separation of resolution and emission

The goal was to split `generateCode` into distinct responsibilities: type resolution happens before `generateCode` is called, and `generateCode` itself only emits Go code.

---

### New file: `resolve_query.go`

Contains two things:

**`resolvedQuery` struct** — all pre-resolved data needed to emit one query:
```go
type resolvedQuery struct {
    name       string
    sqlConst   string      // SQL with $N params (after named param conversion)
    paramNames []string
    paramTypes []string
    t          string      // "one", "many", "exec", "execresult"
    rowFields  []gen.Field // set for one/many queries
    scanFields []string    // set for one/many queries
}
```

**`(c *cli) resolveQuery(q query) (resolvedQuery, error)`** — does all validation and type resolution:
1. Query type guard (`one/many/exec/execresult`, else error)
2. Named param conversion (`@name` → `$N`)
3. `postgresparser.ParseSQLStrict`
4. Sequential param check (no gaps in `$1, $2, ...`)
5. UPDATE/DELETE WHERE guard
6. RETURNING + query type compatibility check
7. `c.resolveParams` → param names/types
8. `c.resolveProjections` (SELECT :one/:many) or `c.resolveReturning` (RETURNING) → row fields/scan fields

Also moved `convertNamedParams`, `pgTypeToGoType`, and the named/positional param regexes into this file since they are part of the resolution step.

---

### `generator.go` — now purely a code emitter

- **`(c *cli) generateCode(_ context.Context, queries []resolvedQuery, output io.Writer, std bool) error`** — signature changed: takes `[]resolvedQuery` instead of `[]query`. No resolution logic inside; just iterates `emitQuery`.
- **`emitQuery(f *gen.File, rq resolvedQuery, std bool)`** — unchanged; pure structural code generation switching on `rq.t`.
- **`generateExec` deleted** in a prior step; its logic was already split into `resolveQuery` (validation/resolution) and `emitQuery` (emission).
- Kept `*cli` receiver and `context.Context` parameter on `generateCode` for future use.

---

### `run.go` — resolution loop added

Between `testQueryAgainstDB` and the output file creation, a resolution loop was added:

```go
var resolved []resolvedQuery
for _, q := range queries {
    rq, err := c.resolveQuery(q)
    if err != nil {
        return err
    }
    resolved = append(resolved, rq)
}
// ...
err = c.generateCode(ctx, resolved, out, std)
```

One `generateCode` call per query file, receiving fully resolved data.

---

### `generator_test.go` — test helpers updated

`generateQuery` and `generateQueryOutput` now explicitly call `c.resolveQuery` then `c.generateCode`, matching the new two-step flow in `run.go`.
