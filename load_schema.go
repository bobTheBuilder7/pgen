package main

import (
	"context"
	"fmt"
)

func (c *cli) loadSchemaFromDB(ctx context.Context) error {
	rows, err := c.db.QueryContext(ctx, `
		SELECT table_name, column_name, data_type, is_nullable, table_schema
		FROM information_schema.columns
		WHERE table_schema = 'pg_catalog'
		ORDER BY table_name, ordinal_position
	`)
	if err != nil {
		return fmt.Errorf("querying information_schema: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, dataType, isNullable, tableSchema string
		if err := rows.Scan(&tableName, &columnName, &dataType, &isNullable, &tableSchema); err != nil {
			return err
		}

		col := Column{
			Name:     columnName,
			Type:     dataType,
			Nullable: isNullable == "YES",
		}
		cols, _ := c.tablesCol.Load(tableName)
		c.tablesCol.Store(tableName, append(cols, col))
	}

	return rows.Err()
}
