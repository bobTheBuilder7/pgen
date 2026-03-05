package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	c := &cli{}
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

			return c.parseSchema(ctx, f)
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

			buf := bytesbufferpool.Get()
			err = c.generateCode(queries, buf)
			if err != nil {
				return err
			}

			// formatted, err := format.Source(buf.Bytes())
			// if err != nil {
			// 	return fmt.Errorf("formatting %s: %w", filename, err)
			// }

			err = os.WriteFile(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)), buf.Bytes(), 0644)
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
