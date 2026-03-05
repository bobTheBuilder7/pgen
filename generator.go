package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/utils"
	"github.com/valkdb/postgresparser"
)

func pgTypeToGoType(pgType string) string {
	switch strings.ToLower(pgType) {
	case "bigserial", "bigint", "int8":
		return "int64"
	case "serial", "integer", "int", "int4":
		return "int32"
	case "smallserial", "smallint", "int2":
		return "int16"
	case "boolean", "bool":
		return "bool"
	case "real", "float4":
		return "float32"
	case "double precision", "float8":
		return "float64"
	case "text", "varchar", "character varying", "char", "character", "uuid":
		return "string"
	default:
		return "string"
	}
}

// resolveColumnGoType resolves the Go type for a SELECT projection column.
// It handles table-qualified columns (e.g. u.id) and string literals (e.g. 'foo').
func (c *cli) resolveColumnGoType(col postgresparser.SelectColumn, tables []postgresparser.TableRef) (string, error) {
	expr := strings.TrimSpace(col.Expression)

	// String literal
	if strings.HasPrefix(expr, "'") {
		return "string", nil
	}

	// Parse table_alias.column_name or just column_name
	var tableAlias, colName string
	if parts := strings.SplitN(expr, ".", 2); len(parts) == 2 {
		tableAlias = parts[0]
		colName = parts[1]
	} else {
		colName = parts[0]
	}

	// Find the real table name from the alias
	var tableName string
	if tableAlias != "" {
		for _, t := range tables {
			if t.Alias == tableAlias || t.Name == tableAlias {
				tableName = t.Name
				break
			}
		}
	} else if len(tables) == 1 {
		tableName = tables[0].Name
	}

	if tableName == "" {
		return "string", nil
	}

	// Look up schema columns
	ddlColumns, ok := c.tablesCol.Load(tableName)
	if !ok {
		return "", fmt.Errorf("table %s not found in schema", tableName)
	}

	for _, ddlCol := range ddlColumns {
		if ddlCol.Name == colName {
			return pgTypeToGoType(ddlCol.Type), nil
		}
	}

	return "string", nil
}

func (c *cli) generateCode(queries []Query, output io.Writer) error {
	generatedFile := gen.NewFile("db")

	generatedFile.AddBlock(gen.Import("", "context"))

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
			if query.t != "one" {
				return fmt.Errorf("query type %s not yet supported", query.t)
			}

			// Resolve projected column types for the result struct
			var structFields []gen.Field
			var scanFields []string
			for _, col := range parsedSQL.Columns {
				goType, err := c.resolveColumnGoType(col, parsedSQL.Tables)
				if err != nil {
					return err
				}

				// Use alias if available, otherwise extract the column name from the expression
				fieldName := col.Alias
				if fieldName == "" {
					fieldName = col.Expression
					// Strip table qualifier (e.g. "users.id" -> "id")
					if dotIdx := strings.LastIndex(fieldName, "."); dotIdx != -1 {
						fieldName = fieldName[dotIdx+1:]
					}
				}
				fieldName = utils.ToPascalCase(fieldName)

				structFields = append(structFields, gen.Field{Name: fieldName, Type: goType})
				scanFields = append(scanFields, "&i."+fieldName)
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

			// Generate method
			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, funcParams, "("+rowStructName+", error)",
					gen.Call("row", "q.db.QueryRow", callArgs...),
					gen.Line("var i "+rowStructName),
					gen.Call("err", "row.Scan", stringersFromStrings(scanFields)...),
					gen.Line("return i, err"),
				),
			)

		case postgresparser.QueryCommandInsert:
		case postgresparser.QueryCommandUpdate:
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
