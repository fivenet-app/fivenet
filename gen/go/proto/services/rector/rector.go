package rector

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/perms"
	"go.uber.org/zap"
)

type Server struct {
	RectorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	a      audit.IAuditer
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions, aud audit.IAuditer) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
		a:      aud,
	}
}
