package main

import (
	"context"

	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/valkdb/postgresparser"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "migrations"
const sqlConstSuffix = "SQL"

type cli struct {
	tablesCol syncmap.Map[string, []postgresparser.DDLColumn]
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx)
	if err != nil {
		panic(err.Error())
	}
}
