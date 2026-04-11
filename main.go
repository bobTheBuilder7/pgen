package main

import (
	"context"
	"flag"
	"fmt"
)

const sqlConstSuffix = "SQL"

func main() {
	std := flag.Bool("std", false, "generate base file for database/sql instead of pgx")
	debug := flag.Bool("d", false, "print generated files to stdout instead of writing to disk")
	dbDir := flag.String("db", "db", "path to the db directory")
	queriesDir := flag.String("queries", "queries", "subdirectory under db containing query files")
	migrationsDir := flag.String("migrations", "migrations", "subdirectory under db containing migration files")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx, *std, *debug, *dbDir, *queriesDir, *migrationsDir)
	if err != nil {
		fmt.Println(err.Error())
	}
}
