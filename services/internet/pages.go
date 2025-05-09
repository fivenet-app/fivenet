package internet

import (
	"context"

	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsinternet "github.com/fivenet-app/fivenet/v2025/services/internet/errors"
)

func (s *Server) GetPage(ctx context.Context, req *pbinternet.GetPageRequest) (*pbinternet.GetPageResponse, error) {
	domain, err := s.getDomainByName(ctx, s.db, req.Domain)
	if err != nil {
		return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
	}
	resp := &pbinternet.GetPageResponse{}

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
