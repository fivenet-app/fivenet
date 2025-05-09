package internet

import (
	"context"
	"errors"

	internet "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/internet"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsinternet "github.com/fivenet-app/fivenet/v2025/services/internet/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getTLD(ctx context.Context, tx qrm.DB, id uint64) (*internet.TLD, error) {
	stmt := tTLDs.
		SELECT(
			tTLDs.ID,
			tTLDs.CreatedAt,
			tTLDs.Name,
			tTLDs.Internal,
			tTLDs.CreatorID,
		).
		FROM(tTLDs).
		WHERE(tTLDs.ID.EQ(jet.Uint64(id))).
		LIMIT(1)

	dest := &internet.TLD{}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsinternet.ErrFailedQuery)
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	return dest, nil
}
