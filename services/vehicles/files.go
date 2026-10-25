package vehicles

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	filepb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	vehiclesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/vehicles/activity"
	pbvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles"
	permsvehicles "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/vehicles/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsvehicles "github.com/fivenet-app/fivenet/v2026/services/vehicles/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc"
)

const (
	vehicleFileKindAttachment = "attachment"
	vehicleFilesPropertyKey   = "vehicle.files"
)

func (s *Server) canUploadVehicleFiles(userInfo *userinfo.UserInfo) error {
	propertyList, err := s.perms.AttrPropertyList(
		userInfo,
		permsvehicles.VehiclesService.SetVehicleProps.Properties,
	)
	if err != nil {
		return errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	if !propertyList.Contains(vehicleFilesPropertyKey) {
		return errorsvehicles.ErrPropsFilesDenied
	}
	return nil
}

func (s *Server) ListVehicleFiles(
	ctx context.Context,
	req *pbvehicles.ListVehicleFilesRequest,
) (*pbvehicles.ListVehicleFilesResponse, error) {
	propertyList, err := s.perms.AttrPropertyListContext(
		ctx,
		auth.MustGetUserInfoFromContext(ctx),
		permsvehicles.VehiclesService.ListVehicles.Properties,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	if !propertyList.Contains(vehicleFilesPropertyKey) {
		return nil, errorsvehicles.ErrPropsPropertiesDenied
	}
	files, err := s.listVehicleFiles(ctx, s.db, req.GetPlate())
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return &pbvehicles.ListVehicleFilesResponse{Files: files}, nil
}

func (s *Server) UploadVehicleFile(
	srv grpc.ClientStreamingServer[filepb.UploadFileRequest, filepb.UploadFileResponse],
) error {
	ctx := srv.Context()
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if err := s.canUploadVehicleFiles(userInfo); err != nil {
		return err
	}
	meta, err := s.filesHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	plate := meta.GetParentKey()
	if plate == "" || meta.GetReason() == "" {
		return errorsvehicles.ErrReasonRequired
	}
	meta.Namespace = "vehicle_files"
	resp, err := s.filesHandler.UploadFromMeta(ctx, meta, plate, srv)
	if err != nil {
		return err
	}
	if err := s.addFileActivity(
		ctx,
		userInfo,
		plate,
		true,
		vehicleFileKindAttachment,
		resp.GetFile(),
		meta.GetReason(),
	); err != nil {
		return errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	return nil
}

func (s *Server) DeleteVehicleFile(
	ctx context.Context,
	req *pbvehicles.DeleteVehicleFileRequest,
) (*pbvehicles.DeleteVehicleFileResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	if err := s.canUploadVehicleFiles(userInfo); err != nil {
		return nil, err
	}
	join := table.FivenetVehicleFiles
	target, err := filestore.FindFileReference(
		ctx,
		s.db,
		join,
		join.Plate.EQ(mysql.String(req.GetPlate())),
		join.FileID,
		join.Kind,
		req.GetFileId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	if target == nil {
		return nil, errorsvehicles.ErrFailedQuery
	}
	if err := s.filesHandler.Delete(ctx, req.GetPlate(), req.GetFileId()); err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	if err := s.addFileActivity(
		ctx,
		userInfo,
		req.GetPlate(),
		false,
		target.Kind,
		target.File,
		req.GetReason(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsvehicles.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbvehicles.DeleteVehicleFileResponse{}, nil
}

func (s *Server) listVehicleFiles(
	ctx context.Context,
	db qrm.Queryable,
	plate string,
) ([]*pbvehicles.VehicleFile, error) {
	join := table.FivenetVehicleFiles
	references, err := filestore.ListFileReferences(
		ctx,
		db,
		join,
		join.Plate.EQ(mysql.String(plate)),
		join.FileID,
		join.Kind,
		join.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	result := make([]*pbvehicles.VehicleFile, 0, len(references))
	for _, reference := range references {
		result = append(result, &pbvehicles.VehicleFile{Kind: reference.Kind, File: reference.File})
	}
	return result, nil
}

func (s *Server) addFileActivity(
	ctx context.Context,
	user *userinfo.UserInfo,
	plate string,
	added bool,
	kind string,
	file *filepb.File,
	reason string,
) error {
	creatorID := user.GetUserId()
	return s.store.AddActivity(
		ctx,
		&vehiclesactivity.VehicleActivity{
			Plate:        plate,
			CreatorId:    &creatorID,
			CreatorJob:   user.GetJob(),
			Reason:       &reason,
			ActivityType: vehiclesactivity.VehicleActivityType_VEHICLE_ACTIVITY_TYPE_FILE,
			Data: &vehiclesactivity.VehicleActivityData{
				Data: &vehiclesactivity.VehicleActivityData_FileChange{
					FileChange: &vehiclesactivity.FileChange{Added: added, Kind: kind, File: file},
				},
			},
		},
	)
}
