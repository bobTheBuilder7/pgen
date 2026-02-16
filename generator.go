package main

import (
	"errors"
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

		if parsedSQL.Command == postgresparser.QueryCommandSelect {
			generatedFile.Const().Id(query.name + "SQL").Op("=").Lit(query.sql)
			generatedFile.Line()
			queryFunc := generatedFile.Func().Params(Id("q").Id("*Queries")).Id(query.name).Params(Id("ctx").Qual("context", "Context"), Id("id").Int64())
			generatedFile.Line()

			queryFunc.Block(
				Id("q").
					Dot("db").
					Dot("QueryRow").
					Call(
						Id("ctx"),
						Id(query.name+"SQL"),
						Id("id"),
					),
			)
		} else {
			return errors.New("not implemented")
		}

		// fmt.Printf("Tables:  %v\n", result.Tables)
		// fmt.Printf("Columns: %v\n", result.Columns)
		// fmt.Printf("Joins:   %v\n", result.JoinConditions)
		// fmt.Printf("Where:   %v\n", result.Where)
		// fmt.Printf("OrderBy: %v\n", result.OrderBy)
	}

	return generatedFile.Render(output)
}
