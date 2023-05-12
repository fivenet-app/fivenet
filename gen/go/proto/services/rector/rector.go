package rector

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"go.uber.org/zap"
)

type Server struct {
	RectorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	a      audit.IAuditer
	c      *mstlystcdata.Enricher
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions, aud audit.IAuditer, c *mstlystcdata.Enricher) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
		a:      aud,
		c:      c,
	}
}
