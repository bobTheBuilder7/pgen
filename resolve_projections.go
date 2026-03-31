package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/utils"
	"github.com/valkdb/postgresparser"
)

var aggregationRegex = regexp.MustCompile(`(?i)^(\w+)\((.+)\)$`)

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

		if strings.Contains(expr, "(") && !strings.HasPrefix(expr, "(") && col.Alias == "" {
			return nil, nil, fmt.Errorf("aggregation %q requires an alias (e.g. ... AS my_name)", expr)
		}

		goType, err := c.resolveColumnGoType(col, tables)
		if err != nil {
			return nil, nil, err
		}

		jsonName := col.Alias
		if jsonName == "" {
			jsonName = col.Expression
			if dotIdx := strings.LastIndex(jsonName, "."); dotIdx != -1 {
				jsonName = jsonName[dotIdx+1:]
			}
		}
		fieldName := utils.ToPascalCase(jsonName)

		structFields = append(structFields, gen.Field{Name: fieldName, Type: goType, Tag: `json:"` + jsonName + `"`})
		scanFields = append(scanFields, "&i."+fieldName)
	}

	return structFields, scanFields, nil
}

// resolveColumnGoType resolves the Go type for a SELECT projection column.
// It handles table-qualified columns (e.g. u.id), string literals (e.g. 'foo'),
// and aggregation functions (COUNT, SUM, COALESCE).
func (c *cli) resolveColumnGoType(col postgresparser.SelectColumn, tables []postgresparser.TableRef) (string, error) {
	expr := strings.TrimSpace(col.Expression)

	// Aggregation function (not a subquery — those start with '(')
	if strings.Contains(expr, "(") && !strings.HasPrefix(expr, "(") {
		pgType, nullable, err := c.resolveAggregationType(expr, tables)
		if err != nil {
			return "", err
		}
		return pgTypeToGoType(pgType, nullable), nil
	}

	// String literal or scalar subquery — type is string
	if strings.HasPrefix(expr, "'") || strings.HasPrefix(expr, "(") {
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
		return "", fmt.Errorf("could not resolve table for column %q", expr)
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

	return "", fmt.Errorf("column %q not found in table %q", colName, tableName)
}

// resolveAggregationType resolves the PG type and nullability for an aggregation
// function expression. Supports COUNT (always int64, never null), SUM (column
// type, always nullable), and COALESCE (inner type, forced non-nullable).
func (c *cli) resolveAggregationType(expr string, tables []postgresparser.TableRef) (pgType string, nullable bool, err error) {
	m := aggregationRegex.FindStringSubmatch(expr)
	if m == nil {
		return "", false, fmt.Errorf("unsupported expression %q: only COUNT, SUM, and COALESCE are supported", expr)
	}

	fn := strings.ToUpper(m[1])
	inner := strings.TrimSpace(m[2])

	switch fn {
	case "COUNT":
		return "int8", false, nil
	case "SUM":
		pgType, err := c.resolveSimpleColumnType(inner, tables)
		if err != nil {
			return "", false, fmt.Errorf("SUM: %w", err)
		}
		return pgType, true, nil
	case "AVG":
		return "float8", true, nil
	case "MIN", "MAX":
		pgType, err := c.resolveSimpleColumnType(inner, tables)
		if err != nil {
			return "", false, fmt.Errorf("%s: %w", fn, err)
		}
		return pgType, true, nil
	case "COALESCE":
		args := splitTopLevelArgs(inner)
		if len(args) == 0 {
			return "", false, fmt.Errorf("COALESCE requires at least one argument")
		}
		pgType, _, err := c.resolveAggregationType(strings.TrimSpace(args[0]), tables)
		if err != nil {
			return "", false, fmt.Errorf("COALESCE: %w", err)
		}
		return pgType, false, nil
	default:
		return "", false, fmt.Errorf("unsupported aggregation function %q: only COUNT, SUM, AVG, MIN, MAX, and COALESCE are supported", fn)
	}
}

// resolveSimpleColumnType resolves a table.column or alias.column expression
// to its raw PostgreSQL type string (e.g. "smallint"), without Go mapping.
func (c *cli) resolveSimpleColumnType(expr string, tables []postgresparser.TableRef) (string, error) {
	expr = strings.TrimSpace(expr)
	var tableAlias, colName string
	if parts := strings.SplitN(expr, ".", 2); len(parts) == 2 {
		tableAlias = parts[0]
		colName = parts[1]
	} else {
		colName = parts[0]
	}

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
		return "", fmt.Errorf("could not resolve table for column %q", expr)
	}

	ddlColumns, ok := c.tablesCol.Load(tableName)
	if !ok {
		return "", fmt.Errorf("table %s not found in schema", tableName)
	}
	for _, ddlCol := range ddlColumns {
		if ddlCol.Name == colName {
			return ddlCol.Type, nil
		}
	}

	return "", fmt.Errorf("column %s not found in table %s", colName, tableName)
}

// splitTopLevelArgs splits a comma-separated argument string, respecting nested
// parentheses. E.g. "SUM(users.age), 0" → ["SUM(users.age)", " 0"].
func splitTopLevelArgs(s string) []string {
	var args []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, s[start:i])
				start = i + 1
			}
		}
	}
	args = append(args, s[start:])
	return args
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
		structFields = append(structFields, gen.Field{Name: fieldName, Type: goType, Tag: `json:"` + col.Column + `"`})
		scanFields = append(scanFields, "&i."+fieldName)
	}

	return structFields, scanFields, nil
}
