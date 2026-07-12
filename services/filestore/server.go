package filestore

import (
	"database/sql"

	pbfilestore "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/storage"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	pbfilestore.FilestoreServiceServer

	st       storage.IStorage
	fHandler *filestore.Handler[int64]
}

type Params struct {
	fx.In

	DB      *sql.DB
	Storage storage.IStorage
}

func NewServer(p Params) *Server {
	fHandler := filestore.NewHandler[int64](
		p.Storage,
		p.DB,
		nil,
		nil,
		nil,
		-1,
		-1,
		nil,
		filestore.InsertJoinRow,
		false,
	)

	return &Server{
		st:       p.Storage,
		fHandler: fHandler,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	pbfilestore.RegisterFilestoreServiceServer(srv, s)
}
