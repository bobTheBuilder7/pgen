package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/assert"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testSharedCli *cli
)

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	ctx := context.Background()

	postgresC, err := testcontainers.Run(
		ctx, "postgres:18",
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		fmt.Printf("failed to start postgres: %v\n", err)
		return 1
	}
	defer postgresC.Terminate(ctx)

	endpoint, err := postgresC.PortEndpoint(ctx, "5432/tcp", "")
	if err != nil {
		fmt.Printf("failed to get endpoint: %v\n", err)
		return 1
	}

	connPool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", endpoint))
	if err != nil {
		fmt.Printf("failed to create connection pool: %v\n", err)
		return 1
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		fmt.Printf("failed to ping postgres: %v\n", err)
		return 1
	}

	testSharedCli = &cli{db: connPool}

	files, err := os.ReadDir(filepath.Join("db", "migrations"))
	if err != nil {
		fmt.Printf("failed to read migrations: %v\n", err)
		return 1
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
			return 1
		}
		err = testSharedCli.runMigration(ctx, name, f)
		f.Close()
		if err != nil {
			fmt.Printf("failed to run migration %s: %v\n", name, err)
			return 1
		}
	}

	if err := testSharedCli.loadSchemaFromDB(ctx); err != nil {
		fmt.Printf("failed to load schema: %v\n", err)
		return 1
	}

	if err := testSharedCli.loadEnumsFromDB(ctx); err != nil {
		fmt.Printf("failed to load enums: %v\n", err)
		return 1
	}

	return m.Run()
}

func testCliWithEmptyDB(t *testing.T) *cli {
	t.Helper()

	postgresC, err := testcontainers.Run(
		t.Context(), "postgres:18",
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	defer testcontainers.CleanupContainer(t, postgresC)
	assert.Nil(t, err)

	endpoint, err := postgresC.PortEndpoint(t.Context(), "5432/tcp", "")
	assert.Nil(t, err)

	connPool, err := pgxpool.New(t.Context(), fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", endpoint))
	assert.Nil(t, err)
	t.Cleanup(connPool.Close)

	if err := connPool.Ping(t.Context()); err != nil {
		t.Fatalf("failed to ping postgres: %v", err)
	}

	return &cli{db: connPool}
}
