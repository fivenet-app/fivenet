package internet

import (
	"context"
	"errors"

	internet "github.com/fivenet-app/fivenet/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorsinternet "github.com/fivenet-app/fivenet/services/internet/errors"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListTLDs(ctx context.Context, req *pbinternet.ListTLDsRequest) (*pbinternet.ListTLDsResponse, error) {
	tTLDs := table.FivenetInternetTlds.AS("tld")

	stmt := tTLDs.
		SELECT(
			tTLDs.ID,
			tTLDs.Name,
		).
		FROM(tTLDs)

	resp := &pbinternet.ListTLDsResponse{
		Tlds: []*internet.TLD{},
	}
	if err := stmt.QueryContext(ctx, s.db, &resp.Tlds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}
	}

	return resp, nil
}
