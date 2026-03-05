package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
	"github.com/valkdb/postgresparser"
)

func (c *cli) parseSchema(ctx context.Context, reader io.Reader) error {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.ReadFrom(reader)

	sql, err := b.ReadString(';')
	if err != nil {
		return err
	}

	parsedSQL, err := postgresparser.ParseSQLStrict(sql)
	if err != nil {
		return err
	}

	if parsedSQL.Command != postgresparser.QueryCommandDDL {
		return errors.New("has to be create table")
	}

	if len(parsedSQL.Tables) != 1 {
		return errors.New("the amout of tables per file is not 1")
	}

	if len(parsedSQL.DDLActions) != 1 {
		return errors.New("DDLActions is not one")
	}

	_, ok := c.tablesCol.Load(parsedSQL.Tables[0].Name)
	if ok {
		return fmt.Errorf("table %s defined twice", parsedSQL.Tables[0].Name)
	}

	c.tablesCol.Store(parsedSQL.Tables[0].Name, parsedSQL.DDLActions[0].ColumnDetails)

	return nil
}
