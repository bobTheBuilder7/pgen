package main

import (
	"context"
	"fmt"
)

func (c *cli) prepareQuery(ctx context.Context, sql, name string, i int) error {
	_, err := c.db.ExecContext(ctx, fmt.Sprintf("PREPARE pgen_test_%d as %s", i, sql))
	if err != nil {
		return fmt.Errorf("invalid query: %s", name)
	}

	return nil
}
