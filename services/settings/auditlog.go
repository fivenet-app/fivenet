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

	// Force the job filter to the user's job unless they are a job admin, in which case they can view all jobs.
	job := userInfo.GetJob()
	if userInfo.GetJobAdmin() {
		job = ""
	}

	resp, err := s.store.ViewAuditLog(ctx, settingsstore.ViewAuditLogOptions{
		Pagination: req.GetPagination(),
		Sort:       req.GetSort(),
		Job:        job,
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

	hydrateShort := s.hydrator.HydrateShortTargetsSafeFunc(userInfo)
	targets := make([]citizenshydrator.ShortTarget, 0, len(resp.GetLogs())*2)
	for i, logEntry := range resp.GetLogs() {
		if logEntry.GetUserId() > 0 {
			targets = append(targets, citizenshydrator.ShortTarget{
				UserID: logEntry.GetUserId(),
				Set: func(user *usershort.UserShort) {
					resp.Logs[i].User = user
				},
			})
		}
		if logEntry.GetTargetUserId() > 0 {
			targets = append(targets, citizenshydrator.ShortTarget{
				UserID: logEntry.GetTargetUserId(),
				Set: func(user *usershort.UserShort) {
					resp.Logs[i].TargetUser = user
				},
			})
		}
	}
	if len(targets) > 0 {
		if err := hydrateShort(ctx, nil, targets); err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
	}

	return resp, nil
}
