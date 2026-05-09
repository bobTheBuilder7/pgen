package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bobTheBuilder7/assert"
	"github.com/bobTheBuilder7/pgen/syncmap"
)

// --- enumDefsFromMap ---

func TestEnumDefsFromMap_Empty(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	defs := enumDefsFromMap(&m)
	assert.Equal(t, len(defs), 0)
}

func TestEnumDefsFromMap_SingleEnum(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser", "trailer", "clip"})
	defs := enumDefsFromMap(&m)
	assert.Equal(t, len(defs), 1)
	assert.Equal(t, defs[0].PgName, "trailer_type")
	assert.Equal(t, defs[0].GoName, "TrailerType")
	assert.Equal(t, len(defs[0].Values), 3)
	assert.Equal(t, defs[0].Values[0].Label, "teaser")
	assert.Equal(t, defs[0].Values[0].GoValue, "Teaser")
	assert.Equal(t, defs[0].Values[1].Label, "trailer")
	assert.Equal(t, defs[0].Values[1].GoValue, "Trailer")
	assert.Equal(t, defs[0].Values[2].Label, "clip")
	assert.Equal(t, defs[0].Values[2].GoValue, "Clip")
}

func TestEnumDefsFromMap_SortedAlphabetically(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser", "trailer", "clip"})
	m.Store("movie_status", []string{"draft", "released", "archived"})
	defs := enumDefsFromMap(&m)
	assert.Equal(t, len(defs), 2)
	assert.Equal(t, defs[0].PgName, "movie_status")
	assert.Equal(t, defs[1].PgName, "trailer_type")
}

func TestEnumDefsFromMap_SnakeCaseGoName(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("movie_status", []string{"draft", "in_progress", "not_found"})
	defs := enumDefsFromMap(&m)
	assert.Equal(t, defs[0].GoName, "MovieStatus")
	assert.Equal(t, defs[0].Values[0].GoValue, "Draft")
	assert.Equal(t, defs[0].Values[1].GoValue, "InProgress")
	assert.Equal(t, defs[0].Values[2].GoValue, "NotFound")
}

// --- generateModelsFile ---

func TestGenerateModelsFile_EmptyEnums(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `package db`)
	assert.False(t, strings.Contains(out, "import"))
}

func TestGenerateModelsFile_TypeDefinition(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser", "trailer", "clip"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `type TrailerType string`)
}

func TestGenerateModelsFile_ConstBlock(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser", "trailer", "clip"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `TrailerTypeTeaser\s+TrailerType = "teaser"`)
	assert.MatchesRegexp(t, out, `TrailerTypeTrailer\s+TrailerType = "trailer"`)
	assert.MatchesRegexp(t, out, `TrailerTypeClip\s+TrailerType = "clip"`)
}

func TestGenerateModelsFile_ScanMethod(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `func \(e \*TrailerType\) Scan`)
	assert.MatchesRegexp(t, out, `\*e = TrailerType\(s\)`)
	assert.MatchesRegexp(t, out, `unsupported scan type for TrailerType`)
}

func TestGenerateModelsFile_NullStruct(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `type NullTrailerType struct`)
	assert.MatchesRegexp(t, out, `TrailerType TrailerType`)
	assert.MatchesRegexp(t, out, `json:"trailer_type"`)
	assert.MatchesRegexp(t, out, `Valid.*bool`)
}

func TestGenerateModelsFile_NullScanMethod(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `func \(ns \*NullTrailerType\) Scan`)
	assert.MatchesRegexp(t, out, `ns\.TrailerType, ns\.Valid = "", false`)
	assert.MatchesRegexp(t, out, `ns\.TrailerType\.Scan\(value\)`)
}

func TestGenerateModelsFile_NullValueMethod(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `func \(ns NullTrailerType\) Value\(\) \(driver\.Value, error\)`)
	assert.MatchesRegexp(t, out, `string\(ns\.TrailerType\)`)
}

func TestGenerateModelsFile_MultipleEnums(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	m.Store("trailer_type", []string{"teaser", "trailer", "clip"})
	m.Store("movie_status", []string{"draft", "released", "archived"})
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `type TrailerType string`)
	assert.MatchesRegexp(t, out, `type MovieStatus string`)
	assert.MatchesRegexp(t, out, `MovieStatusDraft\s+MovieStatus = "draft"`)
	assert.MatchesRegexp(t, out, `MovieStatusReleased\s+MovieStatus = "released"`)
	assert.MatchesRegexp(t, out, `MovieStatusArchived\s+MovieStatus = "archived"`)
}

func TestGenerateModelsFile_HasGeneratedComment(t *testing.T) {
	t.Parallel()
	var m syncmap.Map[string, []string]
	var buf bytes.Buffer
	err := generateModelsFile(&buf, &m)
	assert.Nil(t, err)
	out := buf.String()
	assert.MatchesRegexp(t, out, `Code generated by pgen. DO NOT EDIT.`)
}

// --- loadEnumsFromDB (schema-backed) ---

func TestLoadEnumsFromDB_LoadsTrailerType(t *testing.T) {
	t.Parallel()
	vals, ok := testSharedCli.enums.Load("trailer_type")
	assert.Equal(t, ok, true)
	assert.Equal(t, vals, []string{"teaser", "trailer", "clip"})
}

func TestLoadEnumsFromDB_LoadsMovieStatus(t *testing.T) {
	t.Parallel()
	vals, ok := testSharedCli.enums.Load("movie_status")
	assert.Equal(t, ok, true)
	assert.Equal(t, vals, []string{"draft", "released", "archived"})
}

func TestLoadEnumsFromDB_UnknownTypeNotPresent(t *testing.T) {
	t.Parallel()
	_, ok := testSharedCli.enums.Load("nonexistent_type")
	assert.Equal(t, ok, false)
}
