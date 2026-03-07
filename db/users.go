package db

import "context"
import "github.com/jackc/pgx/v5/pgtype"

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

const UpdateUserNameSQL = "UPDATE users SET name = $1 WHERE users.id = $2;"

func (q *Queries) UpdateUserName(ctx context.Context, name string, id int64) error {
	_, err := q.db.Exec(ctx, UpdateUserNameSQL, name, id)
	return err
}

const CreateUserSQL = "INSERT INTO users (name, age) VALUES ($1, $2);"

func (q *Queries) CreateUser(ctx context.Context, name string, age pgtype.Int4) error {
	_, err := q.db.Exec(ctx, CreateUserSQL, name, age)
	return err
}

type ListUsersRow struct {
	Id   int64
	Name string
}

const ListUsersSQL = "SELECT users.id, users.name FROM users;"

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.Query(ctx, ListUsersSQL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUsersRow
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(&i.Id, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}
