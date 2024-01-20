package rector

import (
	"database/sql"

	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

type Server struct {
	RectorServiceServer

	logger   *zap.Logger
	db       *sql.DB
	ps       perms.Permissions
	aud      audit.IAuditer
	enricher *mstlystcdata.Enricher
	cache    *mstlystcdata.Cache
}

func NewServer(logger *zap.Logger, db *sql.DB, ps perms.Permissions, aud audit.IAuditer, enricher *mstlystcdata.Enricher, cache *mstlystcdata.Cache) *Server {
	return &Server{
		logger:   logger,
		db:       db,
		ps:       ps,
		aud:      aud,
		enricher: enricher,
		cache:    cache,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterRectorServiceServer(srv, s)
}
