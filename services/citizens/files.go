package citizens

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	filepb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbcitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorscitizens "github.com/fivenet-app/fivenet/v2026/services/citizens/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

const (
	citizenFileKindAttachment = "attachment"
	citizenFilesPropertyKey   = "user.files"
	citizenMugshotPropertyKey = "user.mugshot"
)

func (s *Server) ListCitizenFiles(
	ctx context.Context,
	req *pbcitizens.ListCitizenFilesRequest,
) (*pbcitizens.ListCitizenFilesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	propertyList, err := s.perms.AttrPropertyListContext(
		ctx,
		userInfo,
		permscitizens.CitizensService.ListCitizens.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if !propertyList.Contains(citizenFilesPropertyKey) &&
		!propertyList.Contains(citizenMugshotPropertyKey) {
		return nil, errorscitizens.ErrPropsPropertiesDenied
	}
	canView, err := s.canViewCitizenFiles(ctx, userInfo, req.GetUserId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if !canView {
		return nil, errorscitizens.ErrPropsPropertiesDenied
	}
	allowedKinds := citizenFileKinds(propertyList)
	files, err := s.listCitizenFiles(ctx, s.db, req.GetUserId(), allowedKinds)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return &pbcitizens.ListCitizenFilesResponse{Files: files}, nil
}

func citizenFileKinds(propertyList *permissionsattributes.PropertyList) []string {
	if propertyList == nil {
		return nil
	}
	kinds := make([]string, 0, 2)
	for _, key := range propertyList.GetKeys() {
		switch key {
		case citizenFilesPropertyKey:
			kinds = append(kinds, citizenFileKindAttachment)
		case citizenMugshotPropertyKey:
			kinds = append(kinds, "mugshot")
		}
	}
	return kinds
}

func (s *Server) UploadCitizenFile(
	srv grpc.ClientStreamingServer[filepb.UploadFileRequest, filepb.UploadFileResponse],
) error {
	ctx := srv.Context()
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if err := s.canUploadCitizenFiles(userInfo); err != nil {
		return err
	}
	meta, err := s.citizenFilesHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	targetID, ok := utils.ToInt32Checked(meta.GetParentId())
	if meta.GetReason() == "" || !ok || targetID <= 0 {
		return errorscitizens.ErrReasonRequired
	}
	if err := s.canAccessCitizenFilesTarget(ctx, userInfo, targetID); err != nil {
		return err
	}
	logging.InjectFields(ctx, logging.Fields{citizenIDLogFieldKey, targetID})
	meta.Namespace = "citizen_files"
	resp, err := s.citizenFilesHandler.UploadFromMeta(ctx, meta, targetID, srv)
	if err != nil {
		return err
	}
	if err := s.createCitizenFileActivity(
		ctx,
		userInfo.GetUserId(),
		targetID,
		true,
		citizenFileKindAttachment,
		resp.GetFile(),
		meta.GetReason(),
	); err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	return nil
}

func (s *Server) DeleteCitizenFile(
	ctx context.Context,
	req *pbcitizens.DeleteCitizenFileRequest,
) (*pbcitizens.DeleteCitizenFileResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if err := s.canUploadCitizenFiles(userInfo); err != nil {
		return nil, err
	}
	if err := s.canAccessCitizenFilesTarget(ctx, userInfo, req.GetUserId()); err != nil {
		return nil, err
	}
	tUserFiles := table.FivenetUserFiles
	target, err := filestore.FindFileReference(
		ctx,
		s.db,
		tUserFiles,
		tUserFiles.UserID.EQ(mysql.Int32(req.GetUserId())),
		tUserFiles.FileID,
		tUserFiles.Kind,
		req.GetFileId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if target == nil {
		return nil, errorscitizens.ErrCitizenNotFound
	}
	if target.Kind == "mugshot" {
		return nil, errorscitizens.ErrPropsMugshotDenied
	}
	if err := s.citizenFilesHandler.Delete(ctx, req.GetUserId(), req.GetFileId()); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if err := s.createCitizenFileActivity(
		ctx,
		userInfo.GetUserId(),
		req.GetUserId(),
		false,
		target.Kind,
		target.File,
		req.GetReason(),
	); err != nil {
		return nil, errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbcitizens.DeleteCitizenFileResponse{}, nil
}

func (s *Server) canAccessCitizenFilesTarget(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	targetID int32,
) error {
	u, err := s.store.GetUserAccess(ctx, targetID)
	if err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if u == nil || u.GetUserId() <= 0 {
		return errorscitizens.ErrJobGradeNoPermission
	}
	ok, err := s.checkIfUserCanAccess(userInfo, u.GetJob(), u.GetJobGrade())
	if err != nil {
		return err
	}
	if !ok {
		return errorscitizens.ErrJobGradeNoPermission
	}
	return nil
}

func (s *Server) canViewCitizenFiles(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	targetID int32,
) (bool, error) {
	if err := s.canAccessCitizenFilesTarget(ctx, userInfo, targetID); err != nil {
		if errors.Is(err, errorscitizens.ErrJobGradeNoPermission) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Server) canUploadCitizenFiles(userInfo *userinfo.UserInfo) error {
	propertyList, err := s.perms.AttrPropertyList(
		userInfo,
		permscitizens.CitizensService.SetUserProps.Properties,
	)
	if err != nil {
		return errswrap.NewError(err, errorscitizens.ErrFailedQuery)
	}
	if !propertyList.Contains(citizenFilesPropertyKey) {
		return errorscitizens.ErrPropsFilesDenied
	}
	return nil
}

func (s *Server) listCitizenFiles(
	ctx context.Context,
	db qrm.Queryable,
	userID int32,
	kinds []string,
) ([]*pbcitizens.CitizenFile, error) {
	if len(kinds) == 0 {
		return []*pbcitizens.CitizenFile{}, nil
	}
	kindExpressions := make([]mysql.Expression, 0, len(kinds))
	for _, kind := range kinds {
		kindExpressions = append(kindExpressions, mysql.String(kind))
	}
	tUserFiles := table.FivenetUserFiles
	references, err := filestore.ListFileReferences(
		ctx,
		db,
		tUserFiles,
		mysql.AND(
			tUserFiles.UserID.EQ(mysql.Int32(userID)),
			tUserFiles.Kind.IN(kindExpressions...),
		),
		tUserFiles.FileID,
		tUserFiles.Kind,
		tUserFiles.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result := make([]*pbcitizens.CitizenFile, 0, len(references))
	for _, reference := range references {
		result = append(result, &pbcitizens.CitizenFile{Kind: reference.Kind, File: reference.File})
	}
	return result, nil
}

func (s *Server) createCitizenFileActivity(
	ctx context.Context,
	sourceID, targetID int32,
	added bool,
	kind string,
	f *filepb.File,
	reason string,
) error {
	return usersactivity.CreateUserActivities(
		ctx,
		s.db,
		&usersactivity.UserActivity{
			SourceUserId: &sourceID,
			TargetUserId: targetID,
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_FILE,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_FileChange{
					FileChange: &usersactivity.FileChange{Added: added, Kind: kind, File: f},
				},
			},
		},
	)
}
