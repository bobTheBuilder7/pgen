package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
	"github.com/bobTheBuilder7/pgen/syncmap"
	_ "github.com/bradfitz/gopglite"
	"github.com/valkdb/postgresparser"
)

type cli struct {
	tablesCol syncmap.Map[string, []postgresparser.DDLColumn]
	std       bool
	db        *sql.DB
}

func run(ctx context.Context, std bool) error {
	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		return errors.Join(err, errors.New("pglite db failed"))
	}

	c := &cli{std: std, db: db}

	files, err := os.ReadDir(filepath.Join(dbDirectory, schemaDirectory))
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
		f, err := os.Open(filepath.Join(dbDirectory, schemaDirectory, name))
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

	queryFiles, err := os.ReadDir(filepath.Join(dbDirectory, queriesDirectory))
	if err != nil {
		return err
	}

	for _, file := range queryFiles {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		f, err := os.Open(filepath.Join(dbDirectory, queriesDirectory, filename))
		if err != nil {
			return err
		}

		queries, err := parseFileToQueries(ctx, f)
		f.Close()
		if err != nil {
			return err
		}

		buf := bytesbufferpool.Get()
		err = c.generateCode(queries, buf)
		if err != nil {
			bytesbufferpool.Put(buf)
			return err
		}

		err = os.WriteFile(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)), buf.Bytes(), 0644)
		bytesbufferpool.Put(buf)
		if err != nil {
			return err
		}
	}

	baseFile, err := os.Create("./db/db.go")
	if err != nil {
		return err
	}
	defer baseFile.Close()

	err = generateBaseFile(baseFile, std)
	if err != nil {
		return errors.Join(err, errors.New("generateBaseFile failed"))
	}

	return nil
}
