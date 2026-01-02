package db

import (
	sqlc "circa/internal/db/sqlc/generated"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	sqlc.Querier
}

type PGXStore struct {
	*sqlc.Queries
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) Store {
	return &PGXStore{
		Queries: sqlc.New(db),
		db:      db,
	}
}

// GetDB returns the underlying database pool
func (s *PGXStore) GetDB() *pgxpool.Pool {
	return s.db
}
