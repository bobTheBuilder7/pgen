package db

import "context"

type GetUserByIDRow struct {
	UserId   int64
	UserName string
	Suren    string
}

const GetUserByIDSQL = "SELECT u.id as user_id, u.name as user_name, 'asdsda' as suren FROM users u WHERE u.id = $1;"

func (q *Queries) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, GetUserByIDSQL, id)
	var i GetUserByIDRow
	err := row.Scan(&i.UserId, &i.UserName, &i.Suren)
	return i, err
}
