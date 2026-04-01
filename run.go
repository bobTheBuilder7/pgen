package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/bobTheBuilder7/gopglite"
	"github.com/bobTheBuilder7/pgen/syncmap"
)

type dbColumn struct {
	Name     string
	Type     string
	Nullable bool
}

type cli struct {
	tablesCol syncmap.Map[string, []dbColumn]
	db        *sql.DB
}

func run(ctx context.Context, std bool) error {
	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		return errors.Join(err, errors.New("pglite db failed"))
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	c := &cli{db: db}

	files, err := os.ReadDir(filepath.Join(dbDirectory, migrationsDirectory))
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
		f, err := os.Open(filepath.Join(dbDirectory, migrationsDirectory, name))
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
			return fmt.Errorf("%s shouldn't be in queries directory", filename)
		}

		f, err := os.Open(filepath.Join(dbDirectory, queriesDirectory, filename))
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

		out, err := os.Create(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)))
		if err != nil {
			return err
		}
		defer out.Close()

		err = c.generateCode(ctx, queries, out, std)
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
