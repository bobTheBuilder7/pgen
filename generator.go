package main

import (
	"errors"
	"fmt"
	"io"

	. "github.com/dave/jennifer/jen"
	"github.com/valkdb/postgresparser"
)

func generateCode(queries []Query, output io.Writer) error {
	generatedFile := NewFile("db")

	for _, query := range queries {
		parsedSQL, err := postgresparser.ParseSQL(query.sql)
		if err != nil {
			return err
		}

		for _, table := range parsedSQL.Tables {
			_, ok := tablesCol.Load(table.Name)
			if !ok {
				return fmt.Errorf("query: %s access invalid table %s", query.name, table.Name)
			}
		}

		switch command := parsedSQL.Command; command {
		case postgresparser.QueryCommandSelect:
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

			// generatedFile.Const().Id(query.name + "SQL").Op("=").Lit(query.sql)
			// generatedFile.Line()
			// queryFunc := generatedFile.Func().Params(Id("q").Id("*Queries")).Id(query.name).Params(Id("ctx").Qual("context", "Context"), Id("id").Int64())
			// generatedFile.Line()

			// queryFunc.Block(
			// 	Id("q").
			// 		Dot("db").
			// 		Dot("QueryRow").
			// 		Call(
			// 			Id("ctx"),
			// 			Id(query.name+"SQL"),
			// 			Id("id"),
			// 		),
			// )
		case postgresparser.QueryCommandInsert:
		case postgresparser.QueryCommandUpdate:
		case postgresparser.QueryCommandDelete:
		default:
			return errors.New("not implemented")
		}
	}
	return generatedFile.Render(output)

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
