package main

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/valkdb/postgresparser"
)

var namedParamRegex = regexp.MustCompile(`@(\w+)`)
var positionalParamRegex = regexp.MustCompile(`\$\d+`)

// convertNamedParams detects @name style parameters in SQL and converts them
// to positional $N parameters. Returns the converted SQL and the ordered list
// of parameter names. If no named params are found, returns the original SQL
// with nil paramNames. Errors if both @name and $N styles are mixed.
func convertNamedParams(sql string) (string, []string, error) {
	hasNamed := namedParamRegex.MatchString(sql)
	hasPositional := positionalParamRegex.MatchString(sql)

	if hasNamed && hasPositional {
		return "", nil, fmt.Errorf("cannot mix named (@name) and positional ($N) parameters in the same query")
	}

	if !hasNamed {
		return sql, nil, nil
	}

	nameToPos := make(map[string]int)
	var paramNames []string
	converted := namedParamRegex.ReplaceAllStringFunc(sql, func(match string) string {
		name := match[1:] // strip @
		if pos, seen := nameToPos[name]; seen {
			return fmt.Sprintf("$%d", pos)
		}
		pos := len(paramNames) + 1
		nameToPos[name] = pos
		paramNames = append(paramNames, name)
		return fmt.Sprintf("$%d", pos)
	})

	return converted, paramNames, nil
}

func pgTypeToGoType(pgType string, nullable bool) string {
	switch strings.ToLower(pgType) {
	case "bigserial", "bigint", "int8":
		if nullable {
			return "pgtype.Int8"
		}
		return "int64"
	case "serial", "integer", "int", "int4":
		if nullable {
			return "pgtype.Int4"
		}
		return "int32"
	case "smallserial", "smallint", "int2":
		if nullable {
			return "pgtype.Int2"
		}
		return "int16"
	case "boolean", "bool":
		if nullable {
			return "pgtype.Bool"
		}
		return "bool"
	case "real", "float4":
		if nullable {
			return "pgtype.Float4"
		}
		return "float32"
	case "double precision", "float8":
		if nullable {
			return "pgtype.Float8"
		}
		return "float64"
	case "text", "varchar", "character varying", "char", "character":
		if nullable {
			return "pgtype.Text"
		}
		return "string"
	case "uuid":
		if nullable {
			return "pgtype.UUID"
		}
		return "string"
	default:
		return "string"
	}
}

func (c *cli) generateCode(queries []Query, output io.Writer) error {
	generatedFile := gen.NewFile("db")

	generatedFile.AddBlock(gen.Import("", "context"))
	generatedFile.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgconn"))
	generatedFile.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgtype"))

	for _, query := range queries {
		switch query.t {
		case "one", "many", "exec", "execresult":
		default:
			return fmt.Errorf("query %q: unknown query type %q, must be one of: one, many, exec, execresult", query.name, query.t)
		}

		// Convert @name params to $N if present
		sqlForParsing, namedParams, err := convertNamedParams(query.sql)
		if err != nil {
			return err
		}
		// Use converted SQL for both parsing and the generated const (pgx needs $N)
		sqlForConst := sqlForParsing

		parsedSQL, err := postgresparser.ParseSQLStrict(sqlForParsing)
		if err != nil {
			return err
		}

		if parsedSQL.Command == postgresparser.QueryCommandUpdate || parsedSQL.Command == postgresparser.QueryCommandDelete {
			hasFilter := false
			for _, cu := range parsedSQL.ColumnUsage {
				if cu.UsageType == postgresparser.ColumnUsageTypeFilter {
					hasFilter = true
					break
				}
			}
			if !hasFilter {
				return fmt.Errorf("query %q: %s without WHERE clause is not allowed", query.name, parsedSQL.Command)
			}
		}

		switch command := parsedSQL.Command; command {
		case postgresparser.QueryCommandSelect:
			structFields, scanFields, err := c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables)
			if err != nil {
				return err
			}

			rowStructName := query.name + "Row"

			// Resolve function parameters from WHERE clause
			paramNames, paramTypes, err := c.resolveParams(parsedSQL)
			if err != nil {
				return err
			}

			// Override param names if using named params
			if namedParams != nil {
				paramNames = namedParams
			}

			// Build function signature params
			funcParams := "ctx context.Context"
			var callArgs []fmt.Stringer
			callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+sqlConstSuffix))
			for i, name := range paramNames {
				funcParams += ", " + name + " " + paramTypes[i]
				callArgs = append(callArgs, gen.Arg(name))
			}

			// Generate struct
			generatedFile.AddBlock(gen.Struct(rowStructName, structFields...))

			// Generate SQL const
			generatedFile.AddBlock(gen.Const(query.name+sqlConstSuffix, gen.String(sqlForConst)))

			switch query.t {
			case "one":
				generatedFile.AddBlock(
					gen.MethodFunc("q *Queries", query.name, funcParams, "("+rowStructName+", error)",
						gen.Call("row", "q.db.QueryRow", callArgs...),
						gen.Line("var i "+rowStructName),
						gen.Call("err", "row.Scan", stringersFromStrings(scanFields)...),
						gen.Line("return i, err"),
					),
				)
			case "many":
				body := []fmt.Stringer{
					gen.Call("rows, err", "q.db.Query", callArgs...),
					gen.Line("if err != nil {"),
					gen.Line("return nil, err"),
					gen.Line("}"),
					gen.Line("defer rows.Close()"),
					gen.Line("var items []" + rowStructName),
					gen.Line("for rows.Next() {"),
					gen.Line("var i " + rowStructName),
				}
				body = append(body,
					gen.Line("if err := rows.Scan("+strings.Join(scanFields, ", ")+"); err != nil {"),
					gen.Line("return nil, err"),
					gen.Line("}"),
					gen.Line("items = append(items, i)"),
					gen.Line("}"),
					gen.Line("return items, rows.Err()"),
				)
				generatedFile.AddBlock(
					gen.MethodFunc("q *Queries", query.name, funcParams, "([]"+rowStructName+", error)", body...),
				)
			case "exec", "execresult":
				return c.generateExec(generatedFile, query, parsedSQL, sqlForConst, namedParams)
			default:
				return fmt.Errorf("query type %s not supported for SELECT", query.t)
			}

		case postgresparser.QueryCommandInsert:
			if err := c.generateExec(generatedFile, query, parsedSQL, sqlForConst, namedParams); err != nil {
				return err
			}
		case postgresparser.QueryCommandUpdate:
			if err := c.generateExec(generatedFile, query, parsedSQL, sqlForConst, namedParams); err != nil {
				return err
			}
		case postgresparser.QueryCommandDelete:
			if err := c.generateExec(generatedFile, query, parsedSQL, sqlForConst, namedParams); err != nil {
				return err
			}
		default:
			return errors.New("not implemented")
		}
	}

	_, err := generatedFile.WriteTo(output)
	if err != nil {
		return err
	}

	return nil
}

func (c *cli) generateExec(generatedFile *gen.File, query Query, parsedSQL *postgresparser.ParsedQuery, sqlForConst string, namedParams []string) error {
	hasReturning := len(parsedSQL.Returning) > 0

	if !hasReturning && query.t != "exec" && query.t != "execresult" {
		return fmt.Errorf("query type %s not supported for %s without RETURNING", query.t, parsedSQL.Command)
	}
	if hasReturning && query.t != "one" && query.t != "many" {
		return fmt.Errorf("query type %s not supported for %s with RETURNING (use :one or :many)", query.t, parsedSQL.Command)
	}

	paramNames, paramTypes, err := c.resolveParams(parsedSQL)
	if err != nil {
		return err
	}

	// Override param names if using named params
	if namedParams != nil {
		paramNames = namedParams
	}

	funcParams := "ctx context.Context"
	var callArgs []fmt.Stringer
	callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+sqlConstSuffix))
	for i, name := range paramNames {
		funcParams += ", " + name + " " + paramTypes[i]
		callArgs = append(callArgs, gen.Arg(name))
	}

	generatedFile.AddBlock(gen.Const(query.name+sqlConstSuffix, gen.String(sqlForConst)))

	if hasReturning {
		structFields, scanFields, err := c.resolveReturning(parsedSQL)
		if err != nil {
			return err
		}

		rowStructName := query.name + "Row"
		generatedFile.AddBlock(gen.Struct(rowStructName, structFields...))

		switch query.t {
		case "one":
			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "("+rowStructName+", error)",
					gen.Call("row", "q.db.QueryRow", callArgs...),
					gen.Line("var i "+rowStructName),
					gen.Call("err", "row.Scan", stringersFromStrings(scanFields)...),
					gen.Line("return i, err"),
				),
			)
		case "many":
			body := []fmt.Stringer{
				gen.Call("rows, err", "q.db.Query", callArgs...),
				gen.Line("if err != nil {"),
				gen.Line("return nil, err"),
				gen.Line("}"),
				gen.Line("defer rows.Close()"),
				gen.Line("var items []" + rowStructName),
				gen.Line("for rows.Next() {"),
				gen.Line("var i " + rowStructName),
			}
			body = append(body,
				gen.Line("if err := rows.Scan("+strings.Join(scanFields, ", ")+"); err != nil {"),
				gen.Line("return nil, err"),
				gen.Line("}"),
				gen.Line("items = append(items, i)"),
				gen.Line("}"),
				gen.Line("return items, rows.Err()"),
			)
			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "([]"+rowStructName+", error)", body...),
			)
		}
		return nil
	}

	switch query.t {
	case "exec":
		generatedFile.AddBlock(
			gen.MethodFunc("q *Queries", query.name, funcParams, "error",
				gen.Call("_, err", "q.db.Exec", callArgs...),
				gen.Line("return err"),
			),
		)
	case "execresult":
		generatedFile.AddBlock(
			gen.MethodFunc("q *Queries", query.name, funcParams, "(pgconn.CommandTag, error)",
				gen.Line("return q.db.Exec("+buildCallArgsString(callArgs)+")"),
			),
		)
	}

	return nil
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

func filterColumns(columns []postgresparser.ColumnUsage, usage postgresparser.ColumnUsageType) []postgresparser.ColumnUsage {
	var cols []postgresparser.ColumnUsage

	for _, col := range columns {
		if col.UsageType == usage {
			cols = append(cols, col)
		}
	}

	return cols
}
