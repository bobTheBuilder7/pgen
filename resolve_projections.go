package main

import (
	"fmt"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/utils"
	"github.com/valkdb/postgresparser"
)

// resolveProjections resolves the Go struct fields and scan field references
// from the SELECT column projections. Returns the struct fields for code generation
// and the scan fields (e.g. "&i.Id") for row.Scan().
func (c *cli) resolveProjections(columns []postgresparser.SelectColumn, tables []postgresparser.TableRef) ([]gen.Field, []string, error) {
	var structFields []gen.Field
	var scanFields []string

	for _, col := range columns {
		goType, err := c.resolveColumnGoType(col, tables)
		if err != nil {
			return nil, nil, err
		}

		fieldName := col.Alias
		if fieldName == "" {
			fieldName = col.Expression
			if dotIdx := strings.LastIndex(fieldName, "."); dotIdx != -1 {
				fieldName = fieldName[dotIdx+1:]
			}
		}
		fieldName = utils.ToPascalCase(fieldName)

		structFields = append(structFields, gen.Field{Name: fieldName, Type: goType})
		scanFields = append(scanFields, "&i."+fieldName)
	}

	return structFields, scanFields, nil
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
			return pgTypeToGoType(ddlCol.Type, ddlCol.Nullable), nil
		}
	}

	return "string", nil
}
