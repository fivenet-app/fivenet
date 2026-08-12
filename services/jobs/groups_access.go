package jobs

import (
	"context"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	groupsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	"github.com/go-jet/jet/v2/qrm"
)

type groupAccessManager interface {
	CanUserAccessTarget(
		ctx context.Context,
		targetID int64,
		userInfo *userinfo.UserInfo,
		access int32,
	) (bool, error)
	ListTargetAccess(
		ctx context.Context,
		tx qrm.DB,
		targetID int64,
		opts access.SubjectAccessOptions,
	) (*resourcesaccess.Access, error)
	ReplaceTargetAccess(
		ctx context.Context,
		tx qrm.DB,
		resolver *access.SubjectResolver,
		targetID int64,
		in *resourcesaccess.Access,
		opts access.SubjectAccessOptions,
	) (*access.SubjectAccessChanges, error)
}

func (s *Server) ensureGroupAccess(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	groupID int64,
	level groupsaccess.AccessLevel,
) error {
	if s.groupAccess == nil {
		return errorsjobs.ErrFailedQuery
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
