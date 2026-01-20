package generic

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func PSQLConnect() (*pgx.Conn, error) {

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		return nil, err
	}

	return conn, nil

}
