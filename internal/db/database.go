package db

import (
	"chitchat/internal/db/sqlc"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Store struct {
	Db      *sql.DB
	Queries *sqlc.Queries
}

func Connect() (*Store, error) {
	db, err := sql.Open("sqlite", "data.db")
	if err != nil {
		return nil, err
	}

	store := &Store{
		Db:      db,
		Queries: sqlc.New(db),
	}
	return store, nil
}
