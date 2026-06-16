package settings

import (
	"context"

	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	settingsstore "github.com/fivenet-app/fivenet/v2026/stores/settings"
)

func (s *Server) ViewAuditLog(
	ctx context.Context,
	req *pbsettings.ViewAuditLogRequest,
) (*pbsettings.ViewAuditLogResponse, error) {
	resp, err := s.store.ViewAuditLog(ctx, settingsstore.ViewAuditLogOptions{
		Pagination: req.GetPagination(),
		Sort:       req.GetSort(),
		UserIDs:    req.GetUserIds(),
		From:       req.GetFrom(),
		To:         req.GetTo(),
		Services:   req.GetServices(),
		Methods:    req.GetMethods(),
		Actions:    req.GetActions(),
		Results:    req.GetResults(),
		Search:     req.GetSearch(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	return resp, nil
}
