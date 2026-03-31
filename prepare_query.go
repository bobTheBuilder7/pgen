package main

import (
	"context"
)

func (c *cli) prepareQuery(ctx context.Context, sql string) error {
	// fmt.Println(sql)

	// _, err := c.db.ExecContext(ctx, "EXPLAIN "+sql)
	// if err != nil {
	// 	return fmt.Errorf("invalid query: %w", err)
	// }

	return nil
}
