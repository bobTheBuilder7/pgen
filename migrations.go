package main

import (
	"context"
	"errors"
	"fmt"
	"io"
)

func (c *cli) runMigration(ctx context.Context, name string, r io.Reader) error {
	migration, err := io.ReadAll(r)
	if err != nil {
		return errors.Join(err, errors.New("reading migration failed"))
	}

	_, err = c.db.ExecContext(ctx, string(migration))
	if err != nil {
		return fmt.Errorf("migration %s: %w", name, err)
	}

	return nil
}
