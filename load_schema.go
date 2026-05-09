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
		} else if strings.EqualFold(dataType, "USER-DEFINED") {
			colType = udtName
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

func (c *cli) loadEnumsFromDB(ctx context.Context) error {
	rows, err := c.db.Query(ctx, `
		SELECT t.typname, e.enumlabel
		FROM pg_type t
		JOIN pg_enum e ON e.enumtypid = t.oid
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE n.nspname = 'public'
		ORDER BY t.typname, e.enumsortorder
	`)
	if err != nil {
		return fmt.Errorf("querying pg_enum: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var typName, enumLabel string
		if err := rows.Scan(&typName, &enumLabel); err != nil {
			return err
		}
		vals, _ := c.enums.Load(typName)
		c.enums.Store(typName, append(vals, enumLabel))
	}

	return rows.Err()
}
