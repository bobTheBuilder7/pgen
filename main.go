package main

import (
	"context"
	"flag"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "migrations"
const sqlConstSuffix = "SQL"

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
