package rector

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/perms"
	"go.uber.org/zap"
)

var (
	ignoredGuardPermissions = []string{
		"authservice-setjob",
		"superuser-override",
	}
)

type Server struct {
	RectorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
	}
}
