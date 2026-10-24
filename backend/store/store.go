package store

import "database/sql"

type Store struct {
	dbConn *sql.DB
}

func NewStore(parDBConn *sql.DB) *Store {
	return &Store{
		dbConn: parDBConn,
	}
}
