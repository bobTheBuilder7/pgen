package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
	"github.com/valkdb/postgresparser"
	"golang.org/x/sync/errgroup"
)

type Query struct {
	name string
	t    string
	sql  string
}

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "schema"

func main() {
	files, err := os.ReadDir(filepath.Join(dbDirectory, queriesDirectory))
	if err != nil {
		panic(err.Error())
	}

	g := new(errgroup.Group)

	for _, file := range files {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		g.Go(func() error {
			f, err := os.Open(filepath.Join(dbDirectory, queriesDirectory, filename))
			if err != nil {
				return err
			}
			defer f.Close()

			queries, err := parseFileToQueries(f)
			if err != nil {
				return err
			}

			output, err := os.Create(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)))
			if err != nil {
				return err
			}
			defer output.Close()

			err = generateCode(queries, output)
			if err != nil {
				return err
			}

			return nil
		})
	}

	files, err = os.ReadDir(filepath.Join(dbDirectory, schemaDirectory))
	if err != nil {
		panic(err.Error())
	}

	for _, file := range files {
		filename := file.Name()

		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		g.Go(func() error {
			f, err := os.Open(filepath.Join(dbDirectory, schemaDirectory, filename))
			if err != nil {
				return err
			}
			defer f.Close()

			return parseSchema(f)
		})

	}

	err = g.Wait()
	if err != nil {
		fmt.Println("error happened", err.Error())
	}
}

func parseFileToQueries(reader io.Reader) ([]Query, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	var queries []Query

	b.ReadFrom(reader)

	for {
		line, err := b.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		query := Query{}

		if strings.HasPrefix(line, "-- name:") {
			line = strings.ReplaceAll(line, " ", "")
			parts := strings.Split(line, ":")
			if len(parts) != 3 {
				return queries, errors.New("invalid header")
			}

			query.name = parts[1]
			query.t = parts[2]
		}

		sql, err := b.ReadString(';')
		if err != nil {
			break
		}

		query.sql = strings.ReplaceAll(sql, "\n", " ")

		queries = append(queries, query)
	}

	return queries, nil
}

func parseSchema(reader io.Reader) error {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	b.ReadFrom(reader)

	sql, err := b.ReadString(';')
	if err != nil {
		return err
	}

	parsedSQL, err := postgresparser.ParseSQL(sql)
	if err != nil {
		return err
	}

	if parsedSQL.Command != postgresparser.QueryCommandDDL {
		return errors.New("has to be create table")
	}

	for _, action := range parsedSQL.DDLActions[0].ColumnDetails {
		fmt.Println(action.Nullable)
	}

	return nil
}
