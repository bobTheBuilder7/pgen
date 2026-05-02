package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type dbColumn struct {
	Name     string
	Type     string
	Nullable bool
}

type cli struct {
	tablesCol           syncmap.Map[string, []dbColumn]
	db                  *pgxpool.Pool
	dbDirectory         string
	queriesDirectory    string
	migrationsDirectory string
}

func run(ctx context.Context, std bool, debug bool, dbDirectory string, queriesDirectory string, migrationsDirectory string) error {
	postgresC, err := testcontainers.Run(
		ctx, "postgres:18",
		testcontainers.WithLogger(log.New(io.Discard, "", 0)),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_PASSWORD": "password",
		}),
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
			wait.ForListeningPort("5432/tcp"),
		),
	)
	if err != nil {
		return fmt.Errorf("starting postgres container: %w", err)
	}
	defer postgresC.Terminate(ctx)

	endpoint, err := postgresC.PortEndpoint(ctx, "5432/tcp", "")
	if err != nil {
		return fmt.Errorf("getting postgres endpoint: %w", err)
	}

	connPool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://postgres:password@%s/postgres", endpoint))
	if err != nil {
		return fmt.Errorf("creating connection pool: %w", err)
	}
	defer connPool.Close()

	err = connPool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("pinging postgres: %w", err)
	}

	c := &cli{
		db:                  connPool,
		dbDirectory:         dbDirectory,
		queriesDirectory:    queriesDirectory,
		migrationsDirectory: migrationsDirectory,
	}

	files, err := os.ReadDir(filepath.Join(c.dbDirectory, c.migrationsDirectory))
	if err != nil {
		return err
	}
	var migrationNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			migrationNames = append(migrationNames, file.Name())
		}
	}
	sortMigrations(migrationNames)

	for _, name := range migrationNames {
		f, err := os.Open(filepath.Join(c.dbDirectory, c.migrationsDirectory, name))
		if err != nil {
			return err
		}
		err = c.runMigration(ctx, name, f)
		f.Close()
		if err != nil {
			return err
		}
	}

	if err := c.loadSchemaFromDB(ctx); err != nil {
		return err
	}

	queryFiles, err := os.ReadDir(filepath.Join(c.dbDirectory, c.queriesDirectory))
	if err != nil {
		return err
	}

	for _, file := range queryFiles {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			return fmt.Errorf("%s shouldn't be in queries directory", filename)
		}

		f, err := os.Open(filepath.Join(c.dbDirectory, c.queriesDirectory, filename))
		if err != nil {
			return err
		}

		defer f.Close()

		queries, err := parseFileToQueries(ctx, f)
		if err != nil {
			return err
		}

		for _, query := range queries {
			err = c.testQueryAgainstDB(ctx, query)
			if err != nil {
				return err
			}
		}

		var resolved []resolvedQuery
		for _, q := range queries {
			rq, err := c.resolveQuery(q)
			if err != nil {
				return err
			}

			resolved = append(resolved, rq)
		}

		var out *os.File
		if debug {
			out = os.Stdout
		} else {
			out, err = os.Create(filepath.Join(c.dbDirectory, strings.Replace(filename, ".sql", ".go", 1)))
			if err != nil {
				return err
			}
			defer out.Close()
		}

		err = c.generateCode(ctx, resolved, out, std)
		if err != nil {
			return err
		}
	}

	var baseFile *os.File
	if debug {
		baseFile = os.Stdout
	} else {
		baseFile, err = os.Create(filepath.Join(c.dbDirectory, "db.go"))
		if err != nil {
			return err
		}
		defer baseFile.Close()
	}

	err = generateBaseFile(baseFile, std)
	if err != nil {
		return errors.Join(err, errors.New("generateBaseFile failed"))
	}

	return nil
}
