package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bobTheBuilder7/gen"
	"github.com/valkdb/postgresparser"
)

// resolvedQuery holds all pre-resolved information needed to emit Go code for a single query.
type resolvedQuery struct {
	name       string
	sqlConst   string // SQL with $N params (after named param conversion)
	paramNames []string
	paramTypes []string
	t          string      // "one", "many", "exec", "execresult"
	rowFields  []gen.Field // set for one/many queries
	scanFields []string    // set for one/many queries
}

// resolveQuery validates and resolves all type information for a query.
// It returns a resolvedQuery ready to be handed to emitQuery.
func (c *cli) resolveQuery(q query) (resolvedQuery, error) {
	switch q.t {
	case "one", "many", "exec", "execresult":
	default:
		return resolvedQuery{}, fmt.Errorf("query %q: unknown query type %q, must be one of: one, many, exec, execresult", q.name, q.t)
	}

	sqlForParsing, namedParams, err := convertNamedParams(q.sql)
	if err != nil {
		return resolvedQuery{}, err
	}

	parsedSQL, err := postgresparser.ParseSQLStrict(sqlForParsing)
	if err != nil {
		return resolvedQuery{}, err
	}

	// Validate that positional parameters are sequential with no gaps.
	seen := make(map[int]bool)
	for _, p := range parsedSQL.Parameters {
		seen[p.Position] = true
	}
	for i := 1; i <= len(seen); i++ {
		if !seen[i] {
			return resolvedQuery{}, fmt.Errorf("query %q: parameter positions are not sequential (missing $%d)", q.name, i)
		}
	}

	// UPDATE/DELETE must have a WHERE clause.
	if parsedSQL.Command == postgresparser.QueryCommandUpdate || parsedSQL.Command == postgresparser.QueryCommandDelete {
		hasFilter := false
		for _, cu := range parsedSQL.ColumnUsage {
			if cu.UsageType == postgresparser.ColumnUsageTypeFilter {
				hasFilter = true
				break
			}
		}
		if !hasFilter {
			return resolvedQuery{}, fmt.Errorf("query %q: %s without WHERE clause is not allowed", q.name, parsedSQL.Command)
		}
	}

	hasReturning := len(parsedSQL.Returning) > 0

	// Validate type+command combinations for mutating statements.
	switch parsedSQL.Command {
	case postgresparser.QueryCommandInsert, postgresparser.QueryCommandUpdate, postgresparser.QueryCommandDelete:
		if !hasReturning && (q.t == "one" || q.t == "many") {
			return resolvedQuery{}, fmt.Errorf("query type %s not supported for %s without RETURNING", q.t, parsedSQL.Command)
		}
		if hasReturning && (q.t == "exec" || q.t == "execresult") {
			return resolvedQuery{}, fmt.Errorf("query type %s not supported for %s with RETURNING (use :one or :many)", q.t, parsedSQL.Command)
		}
	}

	paramNames, paramTypes, err := c.resolveParams(parsedSQL)
	if err != nil {
		return resolvedQuery{}, err
	}
	if namedParams != nil {
		paramNames = namedParams
	}

	var rowFields []gen.Field
	var scanFields []string

	switch {
	case parsedSQL.Command == postgresparser.QueryCommandSelect && (q.t == "one" || q.t == "many"):
		rowFields, scanFields, err = c.resolveProjections(parsedSQL.Columns, parsedSQL.Tables, parsedSQL.CTEs)
	case hasReturning:
		rowFields, scanFields, err = c.resolveReturning(parsedSQL)
	}
	if err != nil {
		return resolvedQuery{}, err
	}

	return resolvedQuery{
		name:       q.name,
		sqlConst:   sqlForParsing,
		paramNames: paramNames,
		paramTypes: paramTypes,
		t:          q.t,
		rowFields:  rowFields,
		scanFields: scanFields,
	}, nil
}

var namedParamRegex = regexp.MustCompile(`@(\w+)`)
var positionalParamRegex = regexp.MustCompile(`\$\d+`)

// convertNamedParams detects @name style parameters in SQL and converts them
// to positional $N parameters. Returns the converted SQL and the ordered list
// of parameter names. If no named params are found, returns the original SQL
// with nil paramNames. Errors if both @name and $N styles are mixed.
func convertNamedParams(sql string) (string, []string, error) {
	hasNamed := namedParamRegex.MatchString(sql)
	hasPositional := positionalParamRegex.MatchString(sql)

	if hasNamed && hasPositional {
		return "", nil, fmt.Errorf("cannot mix named (@name) and positional ($N) parameters in the same query")
	}

	if !hasNamed {
		return sql, nil, nil
	}

	nameToPos := make(map[string]int)
	var paramNames []string
	converted := namedParamRegex.ReplaceAllStringFunc(sql, func(match string) string {
		name := match[1:] // strip @
		if pos, seen := nameToPos[name]; seen {
			return fmt.Sprintf("$%d", pos)
		}
		pos := len(paramNames) + 1
		nameToPos[name] = pos
		paramNames = append(paramNames, name)
		return fmt.Sprintf("$%d", pos)
	})

	return converted, paramNames, nil
}

func pgTypeToGoType(pgType string, nullable bool) string {
	if strings.HasSuffix(pgType, "[]") {
		elemGoType := pgTypeToGoType(pgType[:len(pgType)-2], false)
		if nullable {
			return "pgtype.Array[" + elemGoType + "]"
		}
		return "[]" + elemGoType
	}

	switch strings.ToLower(pgType) {
	case "bigserial", "bigint", "int8":
		if nullable {
			return "pgtype.Int8"
		}
		return "int64"
	case "serial", "integer", "int", "int4":
		if nullable {
			return "pgtype.Int4"
		}
		return "int32"
	case "smallserial", "smallint", "int2":
		if nullable {
			return "pgtype.Int2"
		}
		return "int16"
	case "boolean", "bool":
		if nullable {
			return "pgtype.Bool"
		}
		return "bool"
	case "real", "float4":
		if nullable {
			return "pgtype.Float4"
		}
		return "float32"
	case "double precision", "float8":
		if nullable {
			return "pgtype.Float8"
		}
		return "float64"
	case "text", "varchar", "character varying", "char", "character":
		if nullable {
			return "pgtype.Text"
		}
		return "string"
	case "uuid":
		if nullable {
			return "pgtype.UUID"
		}
		return "string"
	case "date":
		if nullable {
			return "pgtype.Date"
		}
		return "pgtype.Date"
	case "timestamp", "timestamp without time zone":
		if nullable {
			return "pgtype.Timestamp"
		}
		return "time.Time"
	case "timestamptz", "timestamp with time zone":
		if nullable {
			return "pgtype.Timestamptz"
		}
		return "time.Time"
	case "time", "time without time zone":
		return "pgtype.Time"
	case "timetz", "time with time zone":
		return "pgtype.Time"
	case "interval":
		return "pgtype.Interval"
	case "json":
		if nullable {
			return "pgtype.JSON"
		}
		return "[]byte"
	case "jsonb":
		if nullable {
			return "pgtype.JSONB"
		}
		return "[]byte"
	case "bytea":
		if nullable {
			return "pgtype.Bytea"
		}
		return "[]byte"
	case "numeric", "decimal":
		return "pgtype.Numeric"
	default:
		return "string"
	}
}
