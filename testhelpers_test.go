package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/bobTheBuilder7/gopglite"
)

const usersSchemaSQL = `CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	email VARCHAR(255) NOT NULL,
	age SMALLINT,
	status SMALLINT NOT NULL,
	role_id INTEGER NOT NULL,
	login_count INTEGER,
	org_id BIGINT NOT NULL,
	referrer_id BIGINT,
	active BOOLEAN DEFAULT TRUE,
	verified BOOLEAN NOT NULL DEFAULT FALSE
);`

const postsSchemaSQL = `CREATE TABLE posts (id BIGSERIAL PRIMARY KEY, title TEXT NOT NULL, body TEXT, user_id BIGINT NOT NULL, published BOOLEAN NOT NULL DEFAULT FALSE);`

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

	if err := testSharedCli.runMigration(ctx, "users.up.sql", strings.NewReader(usersSchemaSQL)); err != nil {
		fmt.Printf("failed to run users migration: %v\n", err)
		os.Exit(1)
	}

	if err := testSharedCli.runMigration(ctx, "posts.up.sql", strings.NewReader(postsSchemaSQL)); err != nil {
		fmt.Printf("failed to run users migration: %v\n", err)
		os.Exit(1)
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
