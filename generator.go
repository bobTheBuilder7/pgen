package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/valkdb/postgresparser"
)

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
	generatedFile.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgtype"))

	for _, query := range queries {
		parsedSQL, err := postgresparser.ParseSQLStrict(query.sql)
		if err != nil {
			return err
		}

		if len(parsedSQL.Tables) != 1 {
			return errors.New("complex query: multiple tables not yet supported")
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

			// Build function signature params
			funcParams := "ctx context.Context"
			var callArgs []fmt.Stringer
			callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+"SQL"))
			for i, name := range paramNames {
				funcParams += ", " + name + " " + paramTypes[i]
				callArgs = append(callArgs, gen.Arg(name))
			}

			// Generate struct
			generatedFile.AddBlock(gen.Struct(rowStructName, structFields...))

			// Generate SQL const
			generatedFile.AddBlock(gen.Const(query.name+"SQL", gen.String(query.sql)))

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
			default:
				return fmt.Errorf("query type %s not supported for SELECT", query.t)
			}

		case postgresparser.QueryCommandInsert:
			if query.t != "exec" {
				return fmt.Errorf("query type %s not supported for INSERT", query.t)
			}

			paramNames, paramTypes, err := c.resolveParams(parsedSQL)
			if err != nil {
				return err
			}

			funcParams := "ctx context.Context"
			var callArgs []fmt.Stringer
			callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+"SQL"))
			for i, name := range paramNames {
				funcParams += ", " + name + " " + paramTypes[i]
				callArgs = append(callArgs, gen.Arg(name))
			}

			generatedFile.AddBlock(gen.Const(query.name+"SQL", gen.String(query.sql)))

			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "error",
					gen.Call("_, err", "q.db.Exec", callArgs...),
					gen.Line("return err"),
				),
			)
		case postgresparser.QueryCommandUpdate:
			if query.t != "exec" {
				return fmt.Errorf("query type %s not supported for UPDATE", query.t)
			}

			paramNames, paramTypes, err := c.resolveParams(parsedSQL)
			if err != nil {
				return err
			}

			funcParams := "ctx context.Context"
			var callArgs []fmt.Stringer
			callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+"SQL"))
			for i, name := range paramNames {
				funcParams += ", " + name + " " + paramTypes[i]
				callArgs = append(callArgs, gen.Arg(name))
			}

			generatedFile.AddBlock(gen.Const(query.name+"SQL", gen.String(query.sql)))

			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "error",
					gen.Call("_, err", "q.db.Exec", callArgs...),
					gen.Line("return err"),
				),
			)
		case postgresparser.QueryCommandDelete:
			if query.t != "exec" {
				return fmt.Errorf("query type %s not supported for DELETE", query.t)
			}

			paramNames, paramTypes, err := c.resolveParams(parsedSQL)
			if err != nil {
				return err
			}

			funcParams := "ctx context.Context"
			var callArgs []fmt.Stringer
			callArgs = append(callArgs, gen.Arg("ctx"), gen.Arg(query.name+"SQL"))
			for i, name := range paramNames {
				funcParams += ", " + name + " " + paramTypes[i]
				callArgs = append(callArgs, gen.Arg(name))
			}

			generatedFile.AddBlock(gen.Const(query.name+"SQL", gen.String(query.sql)))

			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "error",
					gen.Call("_, err", "q.db.Exec", callArgs...),
					gen.Line("return err"),
				),
			)
		default:
			return errors.New("not implemented")
		}
	}

	err := generatedFile.WriteTo(output)
	if err != nil {
		return err
	}

	return nil
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

func findTable(tables []postgresparser.TableRef, search string) (postgresparser.TableRef, error) {
	for _, table := range tables {
		if table.Name == search || table.Alias == search {
			return table, nil
		}
	}

	return postgresparser.TableRef{}, fmt.Errorf("table %s not found", search)
}
