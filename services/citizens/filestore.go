package citizens

import (
	context "context"
	"math"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc "google.golang.org/grpc"
)

func (s *Server) UploadAvatar(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, userInfo.GetUserId()})

	meta, err := s.profilePictureHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	meta.Namespace = "user_profile_pictures"
	if _, err := s.profilePictureHandler.UploadFromMeta(
		ctx,
		meta,
		userInfo.GetUserId(),
		srv,
	); err != nil {
		return err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return nil
}

func (s *Server) DeleteAvatar(
	ctx context.Context,
	req *pbcitizens.DeleteAvatarRequest,
) (*pbcitizens.DeleteAvatarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	avatarFileID, err := s.store.GetAvatarFileID(ctx, userInfo.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if avatarFileID == nil || *avatarFileID == 0 {
		return &pbcitizens.DeleteAvatarResponse{}, nil
	}

	if err := s.profilePictureHandler.Delete(
		ctx,
		userInfo.GetUserId(),
		*avatarFileID,
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.DeleteAvatarResponse{}, nil
}

func (s *Server) UploadMugshot(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, userInfo.GetUserId()})

	// Field Permission Check
	propertyList, err := s.perms.AttrPropertyList(
		userInfo,
		permscitizens.CitizensService.SetUserProps.Properties,
	)
	if err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if !propertyList.Contains("user.mugshot") {
		return errorscitizens.ErrPropsMugshotDenied
	}

	meta, err := s.mugshotHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	if meta.GetReason() == "" {
		return errorscitizens.ErrReasonRequired
	}

	parentId := meta.GetParentId()
	if parentId <= 0 || parentId > math.MaxInt32 {
		return errorscitizens.ErrPropsMugshotDenied
	}
	targetUserId := int32(parentId)

	u, err := s.store.GetUserAccess(ctx, targetUserId)
	if err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if u == nil || u.GetUserId() <= 0 {
		return errorscitizens.ErrJobGradeNoPermission
	}

	check, err := s.checkIfUserCanAccess(userInfo, u.GetJob(), u.GetJobGrade())
	if err != nil {
		return err
	}
	if !check {
		return errorscitizens.ErrJobGradeNoPermission
	}

	meta.Namespace = "user_mugshots"

	resp, oldMugshotFileID, err := s.mugshotHandler.UploadFromMetaWithReplacement(
		ctx,
		meta,
		targetUserId,
		srv,
		replaceUserMugshot,
	)
	if err != nil {
		return err
	}
	if oldMugshotFileID > 0 && oldMugshotFileID != resp.GetId() {
		if err := s.mugshotHandler.Delete(ctx, targetUserId, oldMugshotFileID); err != nil {
			return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	if oldMugshotFileID == 0 || resp.GetId() != oldMugshotFileID {
		if err := usersactivity.CreateUserActivities(ctx, s.db, &usersactivity.UserActivity{
			SourceUserId: &userInfo.UserId,
			TargetUserId: targetUserId,
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       meta.GetReason(),
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_MugshotChange{
					MugshotChange: &usersactivity.MugshotChange{},
				},
			},
		}); err != nil {
			return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return nil
}

func (s *Server) DeleteMugshot(
	ctx context.Context,
	req *pbcitizens.DeleteMugshotRequest,
) (*pbcitizens.DeleteMugshotResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetReason() == "" {
		return nil, errorscitizens.ErrReasonRequired
	}

	u, err := s.store.GetUserAccess(ctx, req.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if u == nil || u.GetUserId() <= 0 {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	check, err := s.checkIfUserCanAccess(userInfo, u.GetJob(), u.GetJobGrade())
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	mugshotFileID, err := s.mugshotHandler.DeleteWithParentLock(
		ctx,
		u.GetUserId(),
		lockCitizenFileLimitParent,
		findUserMugshot,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if mugshotFileID == 0 {
		return &pbcitizens.DeleteMugshotResponse{}, nil
	}

	if err := usersactivity.CreateUserActivities(ctx, s.db, &usersactivity.UserActivity{
		SourceUserId: new(userInfo.UserId),
		TargetUserId: u.GetUserId(),
		Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
		Reason:       req.GetReason(),
		Data: &usersactivity.UserActivityData{
			Data: &usersactivity.UserActivityData_MugshotChange{
				MugshotChange: &usersactivity.MugshotChange{},
			},
		},
	}); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.DeleteMugshotResponse{}, nil
}
