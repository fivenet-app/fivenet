package centrum

import (
	"context"
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

func (s *Server) CreateDispatch(ctx context.Context, req *CreateDispatchRequest) (*CreateDispatchResponse, error) {

	return nil, nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *UpdateDispatchRequest) (*UpdateDispatchResponse, error) {

	return nil, nil
}

func (s *Server) Stream(req *CentrumStreamRequest, srv CentrumService_StreamServer) error {

	return nil
}
