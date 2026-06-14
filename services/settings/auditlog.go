package settings

import (
	"context"

	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
)

func (s *Server) ViewAuditLog(
	ctx context.Context,
	req *pbsettings.ViewAuditLogRequest,
) (*pbsettings.ViewAuditLogResponse, error) {
	resp, err := s.store.ViewAuditLog(ctx, req)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}
