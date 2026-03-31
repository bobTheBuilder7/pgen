package main

import (
	"context"
	"flag"
	"fmt"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const migrationsDirectory = "migrations"
const sqlConstSuffix = "SQL"

func main() {
	std := flag.Bool("std", false, "generate base file for database/sql instead of pgx")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx, *std)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("success")
	}
}
