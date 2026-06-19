package statsstore

import (
	"context"
	"database/sql"
)

type IStore interface {
	LoadPublicStats(ctx context.Context) (Stats, error)
}

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) IStore {
	return &Store{db: db}
}
