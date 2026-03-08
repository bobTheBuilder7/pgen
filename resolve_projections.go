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
		expr := strings.TrimSpace(col.Expression)
		if expr == "*" || strings.HasSuffix(expr, ".*") {
			return nil, nil, fmt.Errorf("SELECT * is not supported: explicitly list the columns you need")
		}

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
			nullable := ddlCol.Nullable || isOuterJoinNullable(tableName, tables)
			return pgTypeToGoType(ddlCol.Type, nullable), nil
		}
	}

	return "string", nil
}

// isOuterJoinNullable returns true if columns from the given table should be
// treated as nullable due to an outer join. LEFT JOIN makes the joined table
// nullable, RIGHT JOIN makes the base table nullable, FULL makes both nullable.
func isOuterJoinNullable(tableName string, tables []postgresparser.TableRef) bool {
	for _, t := range tables {
		if t.Name == tableName {
			// This table is LEFT or FULL joined → its columns are nullable
			if t.JoinType == "LEFT" || t.JoinType == "FULL" {
				return true
			}
		} else {
			// Another table is RIGHT or FULL joined → base table columns are nullable
			if t.JoinType == "RIGHT" || t.JoinType == "FULL" {
				return true
			}
		}
	}
	return false
}

// resolveReturning resolves the Go struct fields and scan field references
// from RETURNING clause columns in INSERT/UPDATE/DELETE statements.
func (c *cli) resolveReturning(parsedSQL *postgresparser.ParsedQuery) ([]gen.Field, []string, error) {
	returningCols := filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeReturning)
	if len(returningCols) == 0 {
		return nil, nil, fmt.Errorf("no RETURNING columns found")
	}

	var structFields []gen.Field
	var scanFields []string

	for _, col := range returningCols {
		// Resolve table name from alias, default to first table if unqualified
		var tableName string
		if col.TableAlias != "" {
			for _, t := range parsedSQL.Tables {
				if t.Alias == col.TableAlias || t.Name == col.TableAlias {
					tableName = t.Name
					break
				}
			}
		} else if len(parsedSQL.Tables) == 1 {
			tableName = parsedSQL.Tables[0].Name
		}

		if tableName == "" {
			return nil, nil, fmt.Errorf("could not resolve table for returning column %s", col.Column)
		}

		ddlColumns, ok := c.tablesCol.Load(tableName)
		if !ok {
			return nil, nil, fmt.Errorf("table %s not found in schema", tableName)
		}

		var goType string
		for _, ddlCol := range ddlColumns {
			if ddlCol.Name == col.Column {
				goType = pgTypeToGoType(ddlCol.Type, ddlCol.Nullable)
				break
			}
		}
		if goType == "" {
			return nil, nil, fmt.Errorf("column %s not found in table %s", col.Column, tableName)
		}

		fieldName := utils.ToPascalCase(col.Column)
		structFields = append(structFields, gen.Field{Name: fieldName, Type: goType})
		scanFields = append(scanFields, "&i."+fieldName)
	}

	return structFields, scanFields, nil
}
