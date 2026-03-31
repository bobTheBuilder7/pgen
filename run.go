package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/bobTheBuilder7/gopglite"
	"github.com/bobTheBuilder7/pgen/syncmap"
)

type Column struct {
	Name     string
	Type     string
	Nullable bool
}

type cli struct {
	tablesCol syncmap.Map[string, []Column]
	db        *sql.DB
}

func run(ctx context.Context, _ bool) error {
	db, err := sql.Open("pglite", ":memory:")
	if err != nil {
		return errors.Join(err, errors.New("pglite db failed"))
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	db.SetConnMaxIdleTime(0)
	db.SetConnMaxLifetime(0)

	c := &cli{db: db}

	_, err = c.db.ExecContext(ctx, `CREATE TABLE users (
		id BIGSERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL,
		age SMALLINT,
		status SMALLINT NOT NULL,
		role_id INTEGER NOT NULL,
		login_count INTEGER,
		org_id BIGINT NOT NULL,
		referrer_id BIGINT,
		active BOOLEAN DEFAULT TRUE,
		verified BOOLEAN NOT NULL DEFAULT FALSE
);`)
	if err != nil {
		return err
	}

	err = c.loadSchemaFromDB(ctx)
	if err != nil {
		return err
	}

	fmt.Println(c.tablesCol.Load("users"))

	// files, err := os.ReadDir(filepath.Join(dbDirectory, migrationsDirectory))
	// if err != nil {
	// 	return err
	// }

	// var migrationNames []string
	// for _, file := range files {
	// 	if strings.HasSuffix(file.Name(), ".up.sql") {
	// 		migrationNames = append(migrationNames, file.Name())
	// 	}
	// }
	// sortMigrations(migrationNames)

	// for _, name := range migrationNames {
	// 	f, err := os.Open(filepath.Join(dbDirectory, migrationsDirectory, name))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = c.runMigration(ctx, name, f)
	// 	f.Close()
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if err := c.loadSchemaFromDB(ctx); err != nil {
	// 	return err
	// }

	// queryFiles, err := os.ReadDir(filepath.Join(dbDirectory, queriesDirectory))
	// if err != nil {
	// 	return err
	// }

	// for _, file := range queryFiles {
	// 	filename := file.Name()

	// 	if !strings.HasSuffix(filename, ".sql") {
	// 		continue
	// 	}

	// 	f, err := os.Open(filepath.Join(dbDirectory, queriesDirectory, filename))
	// 	if err != nil {
	// 		return err
	// 	}

	// 	queries, err := parseFileToQueries(ctx, f)
	// 	f.Close()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	buf := bytesbufferpool.Get()
	// 	err = c.generateCode(ctx, queries, buf, std)
	// 	if err != nil {
	// 		bytesbufferpool.Put(buf)
	// 		return err
	// 	}

	// 	err = os.WriteFile(filepath.Join(dbDirectory, strings.Replace(filename, ".sql", ".go", 1)), buf.Bytes(), 0644)
	// 	bytesbufferpool.Put(buf)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// baseFile, err := os.Create("./db/db.go")
	// if err != nil {
	// 	return err
	// }
	// defer baseFile.Close()

	// err = generateBaseFile(baseFile, std)
	// if err != nil {
	// 	return errors.Join(err, errors.New("generateBaseFile failed"))
	// }

	return nil
}
