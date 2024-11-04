package wiki

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	jet "github.com/go-jet/jet/v2/mysql"
)

type Access[U any, T protoutils.ProtoMessage[U]] struct {
	db *sql.DB

	table jet.Table
}

func NewAccess[U any, T protoutils.ProtoMessage[U]](db *sql.DB) *Access[U, T] {
	return &Access[U, T]{
		db: db,
	}
}

func (a *Access[U, T]) List() []T {
	return nil
}

// TODO create access manage logic like docstore
// TODO is this the point of making a generic access control system per table type?
