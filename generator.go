package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/utils"
)

func (c *cli) generateCode(_ context.Context, queries []resolvedQuery, output io.Writer, std bool) error {
	f := gen.NewFile("db")

	f.AddBlock(gen.Import("", "context"))
	if std {
		f.AddBlock(gen.Import("", "database/sql"))
	} else {
		f.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgconn"))
		f.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgtype"))
	}
	if usesTimeType(queries) {
		f.AddBlock(gen.Import("", "time"))
	}

	for _, rq := range queries {
		emitQuery(f, rq, std)
	}

	_, err := f.WriteTo(output)
	return err
}

// emitQuery writes Go code for a single resolved query into f.
// It has no knowledge of SQL, schema, or type resolution.
func emitQuery(f *gen.File, rq resolvedQuery, std bool) {
	f.AddBlock(gen.Const(rq.name+sqlConstSuffix, gen.String(rq.sqlConst)))

	funcParams := "ctx context.Context"
	var callArgs []fmt.Stringer
	callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(rq.name+sqlConstSuffix))

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

	rowStructName := rq.name + "Row"

	switch rq.t {
	case "one":
		queryRowMethod := "q.db.QueryRow"
		if std {
			queryRowMethod = "q.db.QueryRowContext"
		}
		f.AddBlock(gen.Struct(rowStructName, rq.rowFields...))
		f.AddBlock(
			gen.MethodFunc("q *Queries", rq.name, funcParams, "("+rowStructName+", error)",
				gen.Call("row", queryRowMethod, callArgs...),
				gen.Line("var i "+rowStructName),
				gen.Call("err", "row.Scan", stringersFromStrings(rq.scanFields)...),
				gen.Line("return i, err"),
			),
		)

	case "many":
		queryManyMethod := "q.db.Query"
		if std {
			queryManyMethod = "q.db.QueryContext"
		}
		body := []fmt.Stringer{
			gen.Call("rows, err", queryManyMethod, callArgs...),
			gen.Line("if err != nil {"),
			gen.Line("return nil, err"),
			gen.Line("}"),
			gen.Line("defer rows.Close()"),
			gen.Line("var items []" + rowStructName),
			gen.Line("for rows.Next() {"),
			gen.Line("var i " + rowStructName),
			gen.Line("if err := rows.Scan(" + strings.Join(rq.scanFields, ", ") + "); err != nil {"),
			gen.Line("return nil, err"),
			gen.Line("}"),
			gen.Line("items = append(items, i)"),
			gen.Line("}"),
			gen.Line("return items, rows.Err()"),
		}
		f.AddBlock(gen.Struct(rowStructName, rq.rowFields...))
		f.AddBlock(
			gen.MethodFunc("q *Queries", rq.name, funcParams, "([]"+rowStructName+", error)", body...),
		)

	case "exec":
		execMethod := "q.db.Exec"
		if std {
			execMethod = "q.db.ExecContext"
		}
		f.AddBlock(
			gen.MethodFunc("q *Queries", rq.name, funcParams, "error",
				gen.Call("_, err", execMethod, callArgs...),
				gen.Line("return err"),
			),
		)

	case "execresult":
		execMethod := "q.db.Exec"
		execResultType := "pgconn.CommandTag"
		if std {
			execMethod = "q.db.ExecContext"
			execResultType = "sql.Result"
		}
		f.AddBlock(
			gen.MethodFunc("q *Queries", rq.name, funcParams, "("+execResultType+", error)",
				gen.Line("return "+execMethod+"("+buildCallArgsString(callArgs)+")"),
			),
		)
	}
}

func buildCallArgsString(args []fmt.Stringer) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = a.String()
	}
	return strings.Join(parts, ", ")
}

func stringersFromStrings(ss []string) []fmt.Stringer {
	out := make([]fmt.Stringer, len(ss))
	for i, s := range ss {
		out[i] = gen.Arg(s)
	}
	return out
}

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
