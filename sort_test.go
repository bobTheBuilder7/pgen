package main

import (
	"testing"

	"github.com/bobTheBuilder7/assert"
)

func TestMigrationSort_NumericPrefix(t *testing.T) {
	t.Parallel()
	files := []string{
		"010_add_column.up.sql",
		"001_create_table.up.sql",
		"002_add_index.up.sql",
	}
	sortMigrations(files)
	assert.Equal(t, files, []string{
		"001_create_table.up.sql",
		"002_add_index.up.sql",
		"010_add_column.up.sql",
	})
}

func TestMigrationSort_DatePrefix(t *testing.T) {
	t.Parallel()
	files := []string{
		"2024_01_03_add_column.up.sql",
		"2024_01_01_create_table.up.sql",
		"2024_01_02_add_index.up.sql",
	}
	sortMigrations(files)
	assert.Equal(t, files, []string{
		"2024_01_01_create_table.up.sql",
		"2024_01_02_add_index.up.sql",
		"2024_01_03_add_column.up.sql",
	})
}

func TestMigrationSort_AlphabeticalNames(t *testing.T) {
	t.Parallel()
	files := []string{
		"create_users.up.sql",
		"add_posts.up.sql",
		"add_comments.up.sql",
	}
	sortMigrations(files)
	assert.Equal(t, files, []string{
		"add_comments.up.sql",
		"add_posts.up.sql",
		"create_users.up.sql",
	})
}

func TestMigrationSort_AlreadySorted(t *testing.T) {
	t.Parallel()
	files := []string{
		"001_a.up.sql",
		"002_b.up.sql",
		"003_c.up.sql",
	}
	sortMigrations(files)
	assert.Equal(t, files, []string{
		"001_a.up.sql",
		"002_b.up.sql",
		"003_c.up.sql",
	})
}

func TestMigrationSort_SingleFile(t *testing.T) {
	t.Parallel()
	files := []string{"001_create_table.up.sql"}
	sortMigrations(files)
	assert.Equal(t, files, []string{"001_create_table.up.sql"})
}
