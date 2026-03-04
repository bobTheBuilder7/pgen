package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	files, err := os.ReadDir(filepath.Join(dbDirectory, schemaDirectory))
	if err != nil {
		return err
	}

	for _, file := range files {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		eg.Go(func() error {
			f, err := os.Open(filepath.Join(dbDirectory, schemaDirectory, filename))
			if err != nil {
				return err
			}
			defer f.Close()

			return parseSchema(ctx, f)
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	files, err = os.ReadDir(filepath.Join(dbDirectory, queriesDirectory))
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		eg.Go(func() error {
			f, err := os.Open(filepath.Join(dbDirectory, queriesDirectory, filename))
			if err != nil {
				return err
			}
			defer f.Close()

			queries, err := parseFileToQueries(ctx, f)
			if err != nil {
				return err
			}

			output, err := os.Create(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)))
			if err != nil {
				return err
			}
			defer output.Close()

			_, err = generateCode(queries, output)
			if err != nil {
				return err
			}

			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}
