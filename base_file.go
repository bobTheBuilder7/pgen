package main

import (
	"io"

	"github.com/bobTheBuilder7/gen"
)

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

	_, err := f.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}
