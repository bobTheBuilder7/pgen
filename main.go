package main

import (
	"context"
	"flag"

	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/valkdb/postgresparser"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "migrations"
const sqlConstSuffix = "SQL"

type cli struct {
	tablesCol syncmap.Map[string, []postgresparser.DDLColumn]
	std       bool
}

func main() {
	std := flag.Bool("std", false, "generate base file for database/sql instead of pgx")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx, *std)
	if err != nil {
		panic(err.Error())
	}
}
