package jobs

import (
	"context"

	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
)

func (s *Server) ensureGroupAccess(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	groupID int64,
	level groupsaccess.AccessLevel,
) error {
	return s.ensureGroupAccessWithDeleted(ctx, userInfo, groupID, level, false)
}

func (s *Server) ensureGroupAccessWithDeleted(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	groupID int64,
	level groupsaccess.AccessLevel,
	includeDeleted bool,
) error {
	if s.groupAccess == nil {
		return errorsjobs.ErrFailedQuery
	}

	if includeDeleted {
		ok, err := s.groupAccess.CanUserAccessTargetIncludingDeleted(
			ctx,
			groupID,
			userInfo,
			int32(level),
		)
		if err != nil {
			return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if !ok {
			return errorsjobs.ErrNotFoundOrNoPerms
		}
		return nil
	}

	ok, err := s.groupAccess.CanUserAccessTarget(ctx, groupID, userInfo, int32(level))
	if err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if !ok {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	return nil
}
