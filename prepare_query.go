package main

import (
	"context"
	"fmt"
)

func (c *cli) testQueryAgainstDB(ctx context.Context, query query) error {
	_, err := c.db.ExecContext(ctx, fmt.Sprintf("PREPARE pgen_test_%s as %s", query.name, query.sql))
	if err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("invalid query: %s", query.name)
	}

	return nil
}
