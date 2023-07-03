package centrum

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/perms"
)

type Server struct {
	CentrumServiceServer

	db *sql.DB
	p  perms.Permissions
	a  audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, aud audit.IAuditer) *Server {
	return &Server{
		db: db,

		p: p,
		a: aud,
	}
}
