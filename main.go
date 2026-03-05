package main

import (
	"context"
	"io"
	"os"

	"github.com/bobTheBuilder7/gen"
	"github.com/bobTheBuilder7/pgen/syncmap"
	"github.com/valkdb/postgresparser"
)

const dbDirectory = "db"
const queriesDirectory = "query"
const schemaDirectory = "schema"

type cli struct {
	tablesCol syncmap.Map[string, []postgresparser.DDLColumn]
}

func generateBaseFile(w io.Writer) error {
	f := gen.NewFile("db")

	f.AddBlock(gen.Import("", "context"))
	f.AddBlock(gen.Import("", "github.com/jackc/pgx/v5"))
	f.AddBlock(gen.Import("", "github.com/jackc/pgx/v5/pgconn"))

	f.AddBlock(gen.Interface("DBTX",
		gen.Method{Name: "Exec", Params: "context.Context, string, ...interface{}", Returns: "(pgconn.CommandTag, error)"},
		gen.Method{Name: "Query", Params: "context.Context, string, ...interface{}", Returns: "(pgx.Rows, error)"},
		gen.Method{Name: "QueryRow", Params: "context.Context, string, ...interface{}", Returns: "pgx.Row"},
		gen.Method{Name: "SendBatch", Params: "context.Context, *pgx.Batch", Returns: "pgx.BatchResults"},
	))

	f.AddBlock(gen.Func("New", "db DBTX", "*Queries",
		gen.Line("return &Queries{db: db}"),
	))

	f.AddBlock(gen.Struct("Queries",
		gen.Field{Name: "db", Type: "DBTX"},
	))

	f.AddBlock(gen.MethodFunc("q *Queries", "WithTx", "tx pgx.Tx", "*Queries",
		gen.Line("return &Queries{db: tx}"),
	))

	err := f.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := run(ctx)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Create("./db/db.go")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	err = generateBaseFile(file)
	if err != nil {
		panic(err.Error())
	}

}
