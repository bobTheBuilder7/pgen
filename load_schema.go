package main

import (
	"context"
	"fmt"
	"strings"
)

func (c *cli) loadSchemaFromDB(ctx context.Context) error {
	rows, err := c.db.Query(ctx, `
		SELECT table_name, column_name, data_type, udt_name, is_nullable, table_schema
		FROM information_schema.columns
		WHERE table_schema = 'public'
		ORDER BY table_name, ordinal_position
	`)
	if err != nil {
		return fmt.Errorf("querying information_schema: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, dataType, udtName, isNullable, tableSchema string
		if err := rows.Scan(&tableName, &columnName, &dataType, &udtName, &isNullable, &tableSchema); err != nil {
			return err
		}

		colType := dataType
		if strings.EqualFold(dataType, "ARRAY") {
			colType = strings.TrimPrefix(udtName, "_") + "[]"
		}

		col := dbColumn{
			Name:     columnName,
			Type:     colType,
			Nullable: isNullable == "YES",
		}
		cols, _ := c.tablesCol.Load(tableName)
		c.tablesCol.Store(tableName, append(cols, col))
	}

	return rows.Err()
}
