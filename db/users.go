package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetUserByIDRow struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
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
	Id   int64  `json:"id"`
	Name string `json:"name"`
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

const DeleteUserByNameSQL = "DELETE FROM users WHERE users.name = $1;"

func (q *Queries) DeleteUserByName(ctx context.Context, name string) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, DeleteUserByNameSQL, name)
}

const UpdateUserAgeSQL = "UPDATE users SET age = $1 WHERE users.id = $2;"

func (q *Queries) UpdateUserAge(ctx context.Context, age pgtype.Int4, id int64) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, UpdateUserAgeSQL, age, id)
}

const InsertUserSQL = "INSERT INTO users (name, age) VALUES ($1, $2);"

func (q *Queries) InsertUser(ctx context.Context, name string, age pgtype.Int4) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, InsertUserSQL, name, age)
}

const InsertUserReturningSQL = "INSERT INTO users (name, age) VALUES ($1, $2) RETURNING id, name;"

type InsertUserReturningRow struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) InsertUserReturning(ctx context.Context, name string, age pgtype.Int4) (InsertUserReturningRow, error) {
	row := q.db.QueryRow(ctx, InsertUserReturningSQL, name, age)
	var i InsertUserReturningRow
	err := row.Scan(&i.Id, &i.Name)
	return i, err
}

type GetUserPostsRow struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	PostId   int64  `json:"post_id"`
	PostName string `json:"post_name"`
}

const GetUserPostsSQL = "SELECT u.id, u.name, p.id as post_id, p.name as post_name FROM users u JOIN posts p ON u.id = p.user_id WHERE u.id = $1;"

func (q *Queries) GetUserPosts(ctx context.Context, id int64) ([]GetUserPostsRow, error) {
	rows, err := q.db.Query(ctx, GetUserPostsSQL, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserPostsRow
	for rows.Next() {
		var i GetUserPostsRow
		if err := rows.Scan(&i.Id, &i.Name, &i.PostId, &i.PostName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

type GetUserByNameRow struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

const GetUserByNameSQL = "SELECT users.id, users.name FROM users WHERE users.id = $1 AND users.name = $2;"

func (q *Queries) GetUserByName(ctx context.Context, user_id int64, user_name string) (GetUserByNameRow, error) {
	row := q.db.QueryRow(ctx, GetUserByNameSQL, user_id, user_name)
	var i GetUserByNameRow
	err := row.Scan(&i.Id, &i.Name)
	return i, err
}
