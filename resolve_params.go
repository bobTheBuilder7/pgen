package main

import (
	"fmt"

	"github.com/valkdb/postgresparser"
)

// resolveParams builds the function parameter list from parsed columns and parameters.
// Returns param names and their Go types, ordered by parameter position ($1, $2, ...).
// For SELECT/DELETE, parameters match filter (WHERE) columns.
// For UPDATE, parameters match dml_set (SET) columns first, then filter (WHERE) columns.
func (c *cli) resolveParams(parsedSQL *postgresparser.ParsedQuery) ([]string, []string, error) {
	var usages []postgresparser.ColumnUsage

	// For INSERT, parameters map directly to InsertColumns by position.
	// For UPDATE, SET columns come first, then WHERE filter columns.
	// For SELECT/DELETE, parameters match WHERE filter columns.
	if parsedSQL.Command == postgresparser.QueryCommandInsert {
		return c.resolveInsertParams(parsedSQL)
	}

	switch parsedSQL.Command {
	case postgresparser.QueryCommandUpdate:
		usages = append(usages, filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeDMLSet)...)
		usages = append(usages, filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeFilter)...)
	default:
		usages = filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeFilter)
	}

	var names []string
	var types []string

	for _, param := range parsedSQL.Parameters {
		if param.Position-1 >= len(usages) {
			return nil, nil, fmt.Errorf("parameter $%d has no matching column usage", param.Position)
		}
		usage := usages[param.Position-1]

		colName := usage.Column
		tableAlias := usage.TableAlias

		// Resolve table name from alias
		var tableName string
		for _, t := range parsedSQL.Tables {
			if t.Alias == tableAlias || t.Name == tableAlias {
				tableName = t.Name
				break
			}
		}

		if tableName == "" {
			return nil, nil, fmt.Errorf("could not resolve table for alias %s", tableAlias)
		}

		ddlColumns, ok := c.tablesCol.Load(tableName)
		if !ok {
			return nil, nil, fmt.Errorf("table %s not found in schema", tableName)
		}

		var goType string
		for _, ddlCol := range ddlColumns {
			if ddlCol.Name == colName {
				nullable := ddlCol.Nullable || isOuterJoinNullable(tableName, parsedSQL.Tables)
				goType = pgTypeToGoType(ddlCol.Type, nullable)
				break
			}
		}
		if goType == "" {
			return nil, nil, fmt.Errorf("column %s not found in table %s", colName, tableName)
		}

		names = append(names, colName)
		types = append(types, goType)
	}

	return names, types, nil
}

// resolveInsertParams resolves parameter names and Go types for INSERT statements.
// Parameters map positionally to InsertColumns: $1 → InsertColumns[0], $2 → InsertColumns[1], etc.
func (c *cli) resolveInsertParams(parsedSQL *postgresparser.ParsedQuery) ([]string, []string, error) {
	if len(parsedSQL.Tables) == 0 {
		return nil, nil, fmt.Errorf("no table found in INSERT statement")
	}

	tableName := parsedSQL.Tables[0].Name

	ddlColumns, ok := c.tablesCol.Load(tableName)
	if !ok {
		return nil, nil, fmt.Errorf("table %s not found in schema", tableName)
	}

	var names []string
	var types []string

	for _, param := range parsedSQL.Parameters {
		if param.Position-1 >= len(parsedSQL.InsertColumns) {
			return nil, nil, fmt.Errorf("parameter $%d has no matching insert column", param.Position)
		}
		colName := parsedSQL.InsertColumns[param.Position-1]

		var goType string
		for _, ddlCol := range ddlColumns {
			if ddlCol.Name == colName {
				goType = pgTypeToGoType(ddlCol.Type, ddlCol.Nullable)
				break
			}
		}
		if goType == "" {
			return nil, nil, fmt.Errorf("column %s not found in table %s", colName, tableName)
		}

		names = append(names, colName)
		types = append(types, goType)
	}

	return names, types, nil
}
