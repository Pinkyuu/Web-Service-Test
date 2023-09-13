package postgresql

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func getDBConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://user:password@localhost/dbname")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeDBConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}
