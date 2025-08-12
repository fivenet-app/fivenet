package citizens

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	users "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbcitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscitizens "github.com/fivenet-app/fivenet/v2025/services/citizens/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc "google.golang.org/grpc"
)

func (s *Server) UploadAvatar(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{"fivenet.citizens.user_id", userInfo.GetUserId()})

	auditEntry := &audit.AuditEntry{
		Service: pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:  "UploadAvatar",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}

	meta, err := s.avatarHandler.AwaitHandshake(srv)
	defer s.aud.Log(auditEntry, meta)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	meta.Namespace = "user_avatars"
	if _, err := s.avatarHandler.UploadFromMeta(ctx, meta, userInfo.GetUserId(), srv); err != nil {
		return err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return nil
}

func (s *Server) DeleteAvatar(
	ctx context.Context,
	req *pbcitizens.DeleteAvatarRequest,
) (*pbcitizens.DeleteAvatarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:  "DeleteAvatar",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, nil)

	stmt := tUserProps.
		SELECT(tUserProps.AvatarFileID.AS("avatar_file_id")).
		WHERE(tUserProps.UserID.EQ(jet.Int32(userInfo.GetUserId()))).
		LIMIT(1)

	var props struct {
		AvatarFileId *uint64
	}
	if err := stmt.QueryContext(ctx, s.db, &props); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	if props.AvatarFileId == nil || *props.AvatarFileId == 0 {
		return &pbcitizens.DeleteAvatarResponse{}, nil
	}

	if err := s.avatarHandler.Delete(ctx, userInfo.GetUserId(), *props.AvatarFileId); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.DeleteAvatarResponse{}, nil
}

func (s *Server) UploadMugshot(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	logging.InjectFields(ctx, logging.Fields{"fivenet.citizens.user_id", userInfo.GetUserId()})

	auditEntry := &audit.AuditEntry{
		Service: pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:  "UploadMugshot",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permscitizens.CitizensServicePerm,
		permscitizens.CitizensServiceSetUserPropsPerm,
		permscitizens.CitizensServiceSetUserPropsFieldsPermField,
	)
	if err != nil {
		s.aud.Log(auditEntry, nil)
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if !fields.Contains("Mugshot") {
		s.aud.Log(auditEntry, nil)
		return errorscitizens.ErrPropsMugshotDenied
	}

	meta, err := s.avatarHandler.AwaitHandshake(srv)
	defer s.aud.Log(auditEntry, meta)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	meta.Namespace = "mugshot"

	if meta.GetReason() == "" {
		return errorscitizens.ErrReasonRequired
	}

	targetUserId := int32(meta.GetParentId())
	if targetUserId <= 0 {
		return errorscitizens.ErrPropsMugshotDenied
	}

	tUser := tables.User().AS("user")

	u := &users.User{}
	stmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(tUser.ID.EQ(jet.Int32(targetUserId))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, u); err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if u.GetUserId() <= 0 {
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

	props, err := s.getUserProps(ctx, userInfo, targetUserId)
	if err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if props.MugshotFileId != nil && props.GetMugshotFileId() > 0 {
		if err := s.mugshotHandler.Delete(ctx, targetUserId, props.GetMugshotFileId()); err != nil {
			return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	resp, err := s.mugshotHandler.UploadFromMeta(ctx, meta, targetUserId, srv)
	if err != nil {
		return err
	}

	if props.MugshotFileId == nil || resp.GetId() != props.GetMugshotFileId() {
		if err := users.CreateUserActivities(ctx, s.db, &users.UserActivity{
			SourceUserId: &userInfo.UserId,
			TargetUserId: targetUserId,
			Type:         users.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       meta.GetReason(),
			Data: &users.UserActivityData{
				Data: &users.UserActivityData_MugshotChange{
					MugshotChange: &users.MugshotChange{},
				},
			},
		}); err != nil {
			return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return nil
}

func (s *Server) DeleteMugshot(
	ctx context.Context,
	req *pbcitizens.DeleteMugshotRequest,
) (*pbcitizens.DeleteMugshotResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcitizens.CitizensService_ServiceDesc.ServiceName,
		Method:  "DeleteMugshot",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, nil)

	if req.GetReason() == "" {
		return nil, errorscitizens.ErrReasonRequired
	}

	tUser := tables.User().AS("user")

	u := &users.User{}
	uStmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
		WHERE(tUser.ID.EQ(jet.Int32(req.GetUserId()))).
		LIMIT(1)

	if err := uStmt.QueryContext(ctx, s.db, u); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if u.GetUserId() <= 0 {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	check, err := s.checkIfUserCanAccess(userInfo, u.GetJob(), u.GetJobGrade())
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, errorscitizens.ErrJobGradeNoPermission
	}

	stmt := tUserProps.
		SELECT(tUserProps.MugshotFileID.AS("mugshot_file_id")).
		WHERE(tUserProps.UserID.EQ(jet.Int32(userInfo.GetUserId()))).
		LIMIT(1)

	var props struct {
		MugshotFileId *uint64
	}
	if err := stmt.QueryContext(ctx, s.db, &props); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
		}
	}

	if props.MugshotFileId == nil || *props.MugshotFileId == 0 {
		return &pbcitizens.DeleteMugshotResponse{}, nil
	}

	if err := s.mugshotHandler.Delete(ctx, userInfo.GetUserId(), *props.MugshotFileId); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	if err := users.CreateUserActivities(ctx, s.db, &users.UserActivity{
		SourceUserId: &userInfo.UserId,
		TargetUserId: req.GetUserId(),
		Type:         users.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
		Reason:       req.GetReason(),
		Data: &users.UserActivityData{
			Data: &users.UserActivityData_MugshotChange{
				MugshotChange: &users.MugshotChange{},
			},
		},
	}); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}

	return &pbcitizens.DeleteMugshotResponse{}, nil
}
