package completor

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
)

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) *Store {
	return &Store{db: db, customDB: customDB}
}
