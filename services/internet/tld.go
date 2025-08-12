package internet

import (
	"context"
	"errors"

	internet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	pbinternet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsinternet "github.com/fivenet-app/fivenet/v2025/services/internet/errors"
	"github.com/go-jet/jet/v2/qrm"
)

var tTLDs = table.FivenetInternetTlds.AS("tld")

func (s *Server) ListTLDs(
	ctx context.Context,
	req *pbinternet.ListTLDsRequest,
) (*pbinternet.ListTLDsResponse, error) {
	stmt := tTLDs.
		SELECT(
			tTLDs.ID,
			tTLDs.CreatedAt,
			tTLDs.Name,
			tTLDs.Internal,
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
