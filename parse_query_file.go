package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/bobTheBuilder7/pgen/bytesbufferpool"
)

type query struct {
	name string
	t    string
	sql  string
}

func parseFileToQueries(_ context.Context, reader io.Reader) ([]query, error) {
	b := bytesbufferpool.Get()
	defer bytesbufferpool.Put(b)

	var queries []query

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

		query := query{}

		if strings.HasPrefix(line, "-- name:") {
			line = strings.ReplaceAll(line, " ", "")
			parts := strings.Split(line, ":")
			if len(parts) != 3 {
				return queries, errors.New("invalid header")
			}

			query.name = parts[1]
			query.t = parts[2]
		} else {
			return queries, errors.New("query doesn't have a header")
		}

		sql, err := b.ReadString(';')
		if err != nil {
			return queries, fmt.Errorf("query %s doesn't have a semicolon", query.name)
		}

		query.sql = sql

		queries = append(queries, query)
	}

	return queries, nil
}
