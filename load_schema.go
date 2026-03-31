package main

import (
	"context"
	"fmt"

	"github.com/valkdb/postgresparser"
)

func (c *cli) loadSchemaFromDB(ctx context.Context) error {
	rows, err := c.db.QueryContext(ctx, `
		SELECT table_name, column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public'
		ORDER BY table_name, ordinal_position
	`)
	if err != nil {
		return fmt.Errorf("querying information_schema: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tableName, columnName, dataType, isNullable string
		if err := rows.Scan(&tableName, &columnName, &dataType, &isNullable); err != nil {
			return err
		}
		col := postgresparser.DDLColumn{
			Name:     columnName,
			Type:     dataType,
			Nullable: isNullable == "YES",
		}
		cols, _ := c.tablesCol.Load(tableName)
		c.tablesCol.Store(tableName, append(cols, col))
	}

	return rows.Err()
}
