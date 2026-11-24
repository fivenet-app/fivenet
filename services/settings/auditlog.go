package settings

import (
	"context"

	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	settingsstore "github.com/fivenet-app/fivenet/v2026/stores/settings"
)

func (s *Server) ViewAuditLog(
	ctx context.Context,
	req *pbsettings.ViewAuditLogRequest,
) (*pbsettings.ViewAuditLogResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Only job and config admins can choose which jobs to view. Everyone else is
	// always restricted to their active job, regardless of the request contents.
	jobs := req.GetJobs()
	if !userInfo.GetJobAdmin() && !userInfo.GetCanBeConfigAdmin() {
		jobs = []string{userInfo.GetJob()}
	}

	resp, err := s.store.ViewAuditLog(ctx, settingsstore.ViewAuditLogOptions{
		Pagination: req.GetPagination(),
		Sort:       req.GetSort(),
		Jobs:       jobs,
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

	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetLogs())*2)
	for i, logEntry := range resp.GetLogs() {
		if logEntry.GetUserId() > 0 {
			targets = append(targets, citizenshydrator.BasicTarget{
				UserID: logEntry.GetUserId(),
				Set: func(user *usershort.UserShort) {
					resp.Logs[i].User = user
				},
			})
		}
		if logEntry.GetTargetUserId() > 0 {
			targets = append(targets, citizenshydrator.BasicTarget{
				UserID: logEntry.GetTargetUserId(),
				Set: func(user *usershort.UserShort) {
					resp.Logs[i].TargetUser = user
				},
			})
		}
	}
	if len(targets) > 0 {
		if err := s.hydrator.HydrateBasicTargets(ctx, s.db, userInfo, targets); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	return resp, nil
}
