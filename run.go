package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
)

func run(ctx context.Context, std bool) error {
	c := &cli{}

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
		err = c.parseSchema(ctx, f)
		f.Close()
		if err != nil {
			return err
		}
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

	return generateBaseFile(baseFile, std)
}
