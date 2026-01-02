package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	// Querier

}

type PGXStore struct {
	// *Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PGXStore{
		// Queries: New(db),
		db: db,
	}
}
