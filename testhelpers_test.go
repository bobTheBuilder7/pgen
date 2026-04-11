package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/bobTheBuilder7/gopglite"
)

var (
	testSharedCli *cli
)

func TestMain(m *testing.M) {
	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		fmt.Printf("failed to open pglite: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	ctx := context.Background()

	testSharedCli = &cli{db: db}

	files, err := os.ReadDir(filepath.Join("db", "migrations"))
	if err != nil {
		fmt.Printf("failed to read migrations: %v\n", err)
		os.Exit(1)
	}
	var migrationNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationNames = append(migrationNames, file.Name())
		}
	}
	sortMigrations(migrationNames)
	for _, name := range migrationNames {
		f, err := os.Open(filepath.Join("db", "migrations", name))
		if err != nil {
			fmt.Printf("failed to open migration %s: %v\n", name, err)
			os.Exit(1)
		}
		err = testSharedCli.runMigration(ctx, name, f)
		f.Close()
		if err != nil {
			fmt.Printf("failed to run migration %s: %v\n", name, err)
			os.Exit(1)
		}
	}

	if err := testSharedCli.loadSchemaFromDB(ctx); err != nil {
		fmt.Printf("failed to load schema: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func testCliWithEmptyDB(t *testing.T) *cli {
	t.Helper()

	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open pglite: %v", err)
	}

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	t.Cleanup(func() { db.Close() })

	return &cli{db: db}
}
