package internal

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func SupabaseConnect() (*pgx.Conn, error) {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	return conn, err

}
