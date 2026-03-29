package main

import (
	"io"

	"github.com/bobTheBuilder7/gen"
)

func generateBaseFile(w io.Writer, std bool) error {
	if std {
		return generateBaseFileStd(w)
	}
	return generateBaseFilePgx(w)
}

func generateBaseFilePgx(w io.Writer) error {
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
	return err
}

func generateBaseFileStd(w io.Writer) error {
	f := gen.NewFile("db")

	f.AddBlock(gen.Import("", "context"))
	f.AddBlock(gen.Import("", "database/sql"))

	f.AddBlock(gen.Interface("DBTX",
		gen.Method{Name: "ExecContext", Params: "context.Context, string, ...interface{}", Returns: "(sql.Result, error)"},
		gen.Method{Name: "QueryContext", Params: "context.Context, string, ...interface{}", Returns: "(*sql.Rows, error)"},
		gen.Method{Name: "QueryRowContext", Params: "context.Context, string, ...interface{}", Returns: "*sql.Row"},
		gen.Method{Name: "PrepareContext", Params: "context.Context, string", Returns: "(*sql.Stmt, error)"},
	))

	f.AddBlock(gen.Func("New", "db DBTX", "*Queries",
		gen.Line("return &Queries{db: db}"),
	))

	f.AddBlock(gen.Struct("Queries",
		gen.Field{Name: "db", Type: "DBTX"},
	))

	f.AddBlock(gen.MethodFunc("q *Queries", "WithTx", "tx *sql.Tx", "*Queries",
		gen.Line("return &Queries{db: tx}"),
	))

	_, err := f.WriteTo(w)
	return err
}
