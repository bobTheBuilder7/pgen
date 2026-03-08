package main

import (
	"context"
	"errors"
	"io"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
)

type Query struct {
	name string
	t    string
	sql  string
}

func parseFileToQueries(ctx context.Context, reader io.Reader) ([]Query, error) {
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
