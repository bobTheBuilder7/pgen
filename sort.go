package main

import (
	"slices"
	"strings"
)

func sortMigrations(files []string) {
	slices.SortFunc(files, func(i, j string) int {
		return strings.Compare(i, j)
	})
}
