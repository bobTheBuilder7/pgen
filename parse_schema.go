package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
	"github.com/valkdb/postgresparser"
)

func (c *cli) parseSchema(_ context.Context, reader io.Reader) error {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.ReadFrom(reader)

	stmts := strings.SplitSeq(b.String(), ";")
	for stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		parsed, err := postgresparser.ParseSQLStrict(stmt)
		if err != nil {
			// Non-parseable statements (comments, unsupported syntax) are skipped
			continue
		}

		for _, action := range parsed.DDLActions {
			switch action.Type {
			case postgresparser.DDLCreateTable:
				if _, ok := c.tablesCol.Load(action.ObjectName); ok {
					return fmt.Errorf("table %s defined twice", action.ObjectName)
				}
				c.tablesCol.Store(action.ObjectName, action.ColumnDetails)

			case postgresparser.DDLDropColumn:
				cols, ok := c.tablesCol.Load(action.ObjectName)
				if !ok {
					return fmt.Errorf("DROP COLUMN on unknown table %s", action.ObjectName)
				}
				for _, colName := range action.Columns {
					cols = removeColumn(cols, colName)
				}
				c.tablesCol.Store(action.ObjectName, cols)
			}
		}
	}

	return nil
}

// removeColumn returns a new slice with the named column removed.
func removeColumn(cols []postgresparser.DDLColumn, name string) []postgresparser.DDLColumn {
	result := cols[:0:0]
	for _, col := range cols {
		if col.Name != name {
			result = append(result, col)
		}
	}
	return result
}
