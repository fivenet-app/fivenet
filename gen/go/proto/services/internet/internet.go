package internet

import (
	"context"
	"database/sql"

	errorsinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type Server struct {
	InternetServiceServer
	AdsServiceServer

	db *sql.DB
}

type Params struct {
	fx.In

	DB *sql.DB
}

func NewServer(p Params) *Server {
	s := &Server{
		db: p.DB,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterInternetServiceServer(srv, s)
	RegisterAdsServiceServer(srv, s)
}

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	domain, err := s.getDomainByName(ctx, req.Domain)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}
	resp := &GetPageResponse{}

	if domain == nil {
		return resp, nil
	}

	page, err := s.getPageByDomainAndPath(ctx, domain.Id, req.Path)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}
	resp.Page = page

	if page != nil {
		page.CreatorJob = nil
		page.CreatorId = nil
	}

	return resp, nil
}
