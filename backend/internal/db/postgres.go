package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dbSource string) (*pgxpool.Pool, error) {
	return pgxpool.New(
		context.Background(),
		dbSource,
	)
}
