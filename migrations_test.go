package main

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/assert"
	_ "github.com/bradfitz/gopglite"
)

func testCliWithDB(t *testing.T) *cli {
	t.Helper()

	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open pglite: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	return &cli{db: db}
}

func TestRunMigration_CreateTable(t *testing.T) {
	c := testCliWithDB(t)
	err := c.runMigration(context.Background(), "001_create_users.up.sql", strings.NewReader(`
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);
	`))
	assert.Nil(t, err)

	// verify table exists
	var count int
	err = c.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'users'`,
	).Scan(&count)
	assert.Nil(t, err)
	assert.Equal(t, count, 1)
}

func TestRunMigration_InvalidSQLReturnsError(t *testing.T) {
	c := testCliWithDB(t)
	err := c.runMigration(context.Background(), "001_bad.up.sql", strings.NewReader(`THIS IS NOT SQL;`))
	assert.NotNil(t, err)
}

func TestRunMigration_ErrorIncludesMigrationName(t *testing.T) {
	c := testCliWithDB(t)
	err := c.runMigration(context.Background(), "001_bad.up.sql", strings.NewReader(`THIS IS NOT SQL;`))
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `001_bad\.up\.sql`)
}

// runMigrations is a test helper that applies a sequence of (name, sql) pairs in order.
func runMigrations(t *testing.T, c *cli, migrations [][2]string) error {
	t.Helper()
	for _, m := range migrations {
		if err := c.runMigration(context.Background(), m[0], strings.NewReader(m[1])); err != nil {
			return err
		}
	}
	return nil
}

func tableColumns(t *testing.T, c *cli, table string) []string {
	t.Helper()
	rows, err := c.db.QueryContext(context.Background(),
		`SELECT column_name FROM information_schema.columns WHERE table_name = $1 ORDER BY ordinal_position`,
		table,
	)
	if err != nil {
		t.Fatalf("querying columns: %v", err)
	}
	defer rows.Close()
	var cols []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			t.Fatalf("scanning column: %v", err)
		}
		cols = append(cols, col)
	}
	return cols
}

// --- multiple migration files ---

func TestMultipleMigrations_TablesCreatedInOrder(t *testing.T) {
	c := testCliWithDB(t)
	err := runMigrations(t, c, [][2]string{
		{"001_create_products.up.sql", `CREATE TABLE products (id SERIAL PRIMARY KEY, name TEXT NOT NULL);`},
		{"002_create_categories.up.sql", `CREATE TABLE categories (id SERIAL PRIMARY KEY, label TEXT NOT NULL);`},
		{"003_create_orders.up.sql", `CREATE TABLE orders (id SERIAL PRIMARY KEY, total INT NOT NULL);`},
	})
	assert.Nil(t, err)

	var count int
	err = c.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_name IN ('products', 'categories', 'orders')`,
	).Scan(&count)
	assert.Nil(t, err)
	assert.Equal(t, count, 3)
}

func TestMultipleMigrations_AlterTableAddColumn(t *testing.T) {
	c := testCliWithDB(t)
	err := runMigrations(t, c, [][2]string{
		{"001_create_vendors.up.sql", `CREATE TABLE vendors (id SERIAL PRIMARY KEY, name TEXT NOT NULL);`},
		{"002_add_email.up.sql", `ALTER TABLE vendors ADD COLUMN email TEXT;`},
		{"003_add_phone.up.sql", `ALTER TABLE vendors ADD COLUMN phone TEXT;`},
	})
	assert.Nil(t, err)

	cols := tableColumns(t, c, "vendors")
	assert.Equal(t, cols, []string{"id", "name", "email", "phone"})
}

func TestMultipleMigrations_DropColumn(t *testing.T) {
	c := testCliWithDB(t)
	err := runMigrations(t, c, [][2]string{
		{"001_create_invoices.up.sql", `CREATE TABLE invoices (id SERIAL PRIMARY KEY, note TEXT, amount INT NOT NULL);`},
		{"002_drop_note.up.sql", `ALTER TABLE invoices DROP COLUMN note;`},
	})
	assert.Nil(t, err)

	cols := tableColumns(t, c, "invoices")
	assert.Equal(t, cols, []string{"id", "amount"})
}

func TestMultipleMigrations_FailedMigrationStopsChain(t *testing.T) {
	c := testCliWithDB(t)
	err := runMigrations(t, c, [][2]string{
		{"001_create_tickets.up.sql", `CREATE TABLE tickets (id SERIAL PRIMARY KEY);`},
		{"002_bad.up.sql", `THIS IS NOT SQL;`},
		{"003_create_comments.up.sql", `CREATE TABLE comments (id SERIAL PRIMARY KEY);`},
	})
	assert.NotNil(t, err)
	assert.MatchesRegexp(t, err.Error(), `002_bad\.up\.sql`)

	// 003 should not have been applied
	var count int
	_ = c.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_name = 'comments'`,
	).Scan(&count)
	assert.Equal(t, count, 0)
}

func TestRunMigration_MultipleStatements(t *testing.T) {
	c := testCliWithDB(t)
	err := c.runMigration(context.Background(), "001_create_tables.up.sql", strings.NewReader(`
		CREATE TABLE employees (id SERIAL PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE posts (id SERIAL PRIMARY KEY, title TEXT NOT NULL);
	`))
	assert.Nil(t, err)

	var count int
	err = c.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM information_schema.tables WHERE table_name IN ('employees', 'posts')`,
	).Scan(&count)
	assert.Nil(t, err)
	assert.Equal(t, count, 2)
}
