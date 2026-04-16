# Phase 12 Session — Knowledge Base

## Changes This Session

### Params struct for 2+ parameters

When a query has 2 or more parameters, `emitQuery` now generates a `QueryNameParams` struct and accepts it as a single `arg QueryNameParams` argument. Queries with 0 or 1 parameters keep individual-arg behaviour.

---

### `generator.go` — conditional param struct in `emitQuery`

Added `"github.com/bobTheBuilder7/pgen/utils"` import.

Replaced the unconditional param-building loop with a branch on `len(rq.paramNames)`:

```go
if len(rq.paramNames) >= 2 {
    paramsStructName := rq.name + "Params"
    var paramFields []gen.Field
    for i, name := range rq.paramNames {
        fieldName := utils.ToPascalCase(name)
        paramFields = append(paramFields, gen.Field{Name: fieldName, Type: rq.paramTypes[i], Tag: `json:"` + name + `"`})
        callArgs = append(callArgs, gen.Arg("arg."+fieldName))
    }
    f.AddBlock(gen.Struct(paramsStructName, paramFields...))
    funcParams += ", arg " + paramsStructName
} else {
    for i, name := range rq.paramNames {
        funcParams += ", " + name + " " + rq.paramTypes[i]
        callArgs = append(callArgs, gen.Arg(name))
    }
}
```

**Design decisions:**
- Threshold is **2+** (not 3+ as originally noted in CLAUDE.md)
- Parameter name for the struct arg is `arg` (matches sqlc convention)
- Struct fields are PascalCase (via `utils.ToPascalCase`), with `json:"..."` tags matching the original param names — consistent with row structs in `resolveProjections`
- Params struct is emitted via `f.AddBlock(gen.Struct(...))` before `gen.MethodFunc`, so output order is: SQL const → params struct → method
- No changes to `resolvedQuery`, `resolveQuery()`, or `resolveParams()` — `paramNames`/`paramTypes` remain parallel arrays; the struct shape is derived at emit time

**Generated output example for 2-param query:**
```go
type GetMovieParams struct {
    Id   int64  `json:"id"`
    Name string `json:"name"`
}

func (q *Queries) GetMovie(ctx context.Context, arg GetMovieParams) (GetMovieRow, error) {
    row := q.db.QueryRow(ctx, GetMovieSQL, arg.Id, arg.Name)
    // ...
}
```

---

### `generator_test.go` — 6 new tests

Added `"strings"` import (for `strings.Contains` in negative assertions).

Group: `// --- Params struct (2+ params) ---`

| Test | Asserts |
|------|---------|
| `TestGenerateCode_TwoParamsGeneratesParamsStruct` | `GetMovieParams` struct name in output |
| `TestGenerateCode_TwoParamsUsesArgInSignature` | `arg GetMovieParams` in function signature |
| `TestGenerateCode_TwoParamsUnpacksFields` | `arg.Id` and `arg.Name` in method body |
| `TestGenerateCode_TwoParamsExecQuery` | Works for `:exec` queries (`DeleteMovieParams`) |
| `TestGenerateCode_OneParamNoParamsStruct` | No `Params` in output for 1-param query |
| `TestGenerateCode_ZeroParamsNoParamsStruct` | No `Params` in output for 0-param query |

Negative assertions use `assert.False(t, strings.Contains(out, "Params"))` — the `assert` package has no `NotMatchesRegexp`.
