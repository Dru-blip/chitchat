package db

import (
	"chitchat/internal/db/sqlc"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Db      *pgxpool.Pool
	Queries *sqlc.Queries
}

func Connect(connString string) (*Store, error) {
	ctx := context.Background()

	dbConn, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	store := &Store{
		Db:      dbConn,
		Queries: sqlc.New(dbConn),
	}
	return store, nil
}
