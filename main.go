package main

import (
	"context"
	"flag"
	"fmt"
)

const dbDirectory = "db"
const queriesDirectory = "queries"
const migrationsDirectory = "migrations"
const sqlConstSuffix = "SQL"

func main() {
	std := flag.Bool("std", false, "generate base file for database/sql instead of pgx")
	debug := flag.Bool("d", false, "print generated files to stdout instead of writing to disk")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx, *std, *debug)
	if err != nil {
		fmt.Println(err.Error())
	}
}
