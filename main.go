package main

import (
	"context"

	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/valkdb/postgresparser"
)

type Query struct {
	name  string
	t     string
	sql   string
	table string
}

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "schema"

var tablesCol = syncmap.Map[string, []postgresparser.DDLColumn]{}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx)
	if err != nil {
		panic(err.Error())
	}
}
