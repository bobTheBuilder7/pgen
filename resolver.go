package main

import (
	"fmt"

	"github.com/valkdb/postgresparser"
)

// resolveParams builds the function parameter list from parsed filter columns and parameters.
// Returns param names and their Go types, ordered by parameter position ($1, $2, ...).
func (c *cli) resolveParams(parsedSQL *postgresparser.ParsedQuery) ([]string, []string, error) {
	filters := filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeFilter)

	var names []string
	var types []string

	for _, param := range parsedSQL.Parameters {
		_ = param
		// Match this parameter to a filter column.
		// For simple queries, filters and parameters are in the same order.
		if param.Position-1 >= len(filters) {
			return nil, nil, fmt.Errorf("parameter $%d has no matching filter column", param.Position)
		}
		filter := filters[param.Position-1]

		colName := filter.Column
		tableAlias := filter.TableAlias

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
				goType = pgTypeToGoType(ddlCol.Type)
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
