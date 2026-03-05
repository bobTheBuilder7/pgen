package db

import "context"

type GetUserByIDRow struct {
	Id   int64
	Name string
}

const GetUserByIDSQL = "SELECT users.id, users.name FROM users WHERE users.id = $1 and users.name = $2;"

func (q *Queries) GetUserByID(ctx context.Context, id int64, name string) (GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, GetUserByIDSQL, id, name)
	var i GetUserByIDRow
	err := row.Scan(&i.Id, &i.Name)
	return i, err
}

const DeleteUserByIDSQL = "DELETE FROM users WHERE users.id = $1;"

func (q *Queries) DeleteUserByID(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, DeleteUserByIDSQL, id)
	return err
}
