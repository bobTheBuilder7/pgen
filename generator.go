package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/postgresparser"
)

func generateCode(queries []Query, output io.Writer) (int64, error) {
	generatedFile := gen.NewFile("db")

	for _, query := range queries {
		parsedSQL, err := postgresparser.ParseSQLStrict(query.sql)
		if err != nil {
			return 0, err
		}

		if len(parsedSQL.Tables) != 1 {
			panic("complex query")
		}

		fmt.Println(parsedSQL.Tables[0].Type)

		// for _, table := range parsedSQL.Tables {
		// 	_, ok := tablesCol.Load(table.Name)
		// 	if !ok {
		// 		return fmt.Errorf("query: %s access invalid table %s", query.name, table.Name)
		// 	}
		// }

		switch command := parsedSQL.Command; command {
		case postgresparser.QueryCommandSelect:

			// for _, table := range parsedSQL.Tables {
			// 	fmt.Println(table.JoinType)
			// }

			// fmt.Println(parsedSQL.ColumnUsage[0])

			// projections := filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeProjection)
			// joins := filterColumns(parsedSQL.ColumnUsage, postgresparser.ColumnUsageTypeJoin)
			// projection := projections[0]
			// _, err := findTable(parsedSQL.Tables, projection.TableAlias)
			// if err != nil {
			// 	return err
			// }

			// fmt.Printf("%+v\n", projection)
			// fmt.Printf("%+v\n", joins)
			// for _, i := range parsedSQL.ColumnUsage {
			// 	fmt.Println(i)
			// }

			generatedFile.AddBlock(
				gen.Const(query.name+"SQL", gen.String(query.sql)),
			)

			generatedFile.AddBlock(
				gen.MethodFunc("q *Queries", query.name, "ctx context.Context, id int64", "(string, error)",
					gen.Call("rows, err", "q.db.Query", gen.Arg("ctx"), gen.Arg(query.name+"SQL")),
					gen.ErrCheck("\"empty\""),
					gen.Line("return \"asdds\", nil"),
				),
			)

		case postgresparser.QueryCommandInsert:
		case postgresparser.QueryCommandUpdate:
		case postgresparser.QueryCommandDelete:
		default:
			return 0, errors.New("not implemented")
		}
	}
	return generatedFile.WriteTo(output)

}

func filterColumns(columns []postgresparser.ColumnUsage, usage postgresparser.ColumnUsageType) []postgresparser.ColumnUsage {
	var cols []postgresparser.ColumnUsage

	for _, col := range columns {
		if col.UsageType == usage {
			cols = append(cols, col)
		}
	}

	return cols
}

func findTable(tables []postgresparser.TableRef, search string) (postgresparser.TableRef, error) {
	for _, table := range tables {
		if table.Name == search || table.Alias == search {
			return table, nil
		}
	}

	return postgresparser.TableRef{}, fmt.Errorf("table %s not found", search)
}
