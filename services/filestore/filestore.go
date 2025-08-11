package filestore

import (
	"context"
	"path/filepath"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	pbfilestore "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

const listFilesPageSize = 50

func (s *Server) ListFiles(ctx context.Context, req *pbfilestore.ListFilesRequest) (*pbfilestore.ListFilesResponse, error) {
	if req.Path != nil {
		logging.InjectFields(ctx, logging.Fields{"fivenet.file.path", *req.Path})
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbfilestore.FilestoreService_ServiceDesc.ServiceName,
		Method:  "ListFiles",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_VIEWED,
	}, req)

	filePath := ""
	if req.Path != nil {
		filePath = filepath.Clean(*req.Path)
	}
	if filePath == "" {
		filePath = "/"
	}

	pag, _ := req.Pagination.GetResponseWithPageSize(database.NoTotalCount, listFilesPageSize)
	resp := &pbfilestore.ListFilesResponse{
		Pagination: pag,
	}

	files, err := s.st.List(ctx, filePath, int(req.Pagination.Offset), listFilesPageSize)
	if err != nil {
		return nil, err
	}

	fs := make([]*file.File, len(files))
	for i := range files {
		fs[i] = &file.File{
			FilePath:    files[i].Name,
			ByteSize:    files[i].Size,
			ContentType: files[i].ContentType,
		}
	}
	resp.Files = fs

	resp.Pagination.Update(len(resp.Files))

	return resp, nil
}

func (s *Server) Upload(srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse]) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbfilestore.FilestoreService_ServiceDesc.ServiceName,
		Method:  "Upload",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}

	meta, err := s.fHandler.AwaitHandshake(srv)
	defer s.aud.Log(auditEntry, meta)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	_, err = s.fHandler.UploadFromMeta(ctx, meta, meta.ParentId, srv)
	defer s.aud.Log(auditEntry, meta)
	if err != nil {
		return err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.file.namespace", meta.Namespace,
		"fivenet.file.name", meta.OriginalName,
	})

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return nil
}

func (s *Server) DeleteFile(ctx context.Context, req *file.DeleteFileRequest) (*file.DeleteFileResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.file.parent_id", req.ParentId,
		"fivenet.file.file_id", req.FileId,
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbfilestore.FilestoreService_ServiceDesc.ServiceName,
		Method:  "Delete",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.fHandler.Delete(ctx, req.ParentId, req.FileId); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &file.DeleteFileResponse{}, nil
}

func (s *Server) DeleteFileByPath(ctx context.Context, req *pbfilestore.DeleteFileByPathRequest) (*pbfilestore.DeleteFileByPathResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.file.path", req.Path})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbfilestore.FilestoreService_ServiceDesc.ServiceName,
		Method:  "DeleteFileByPath",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.fHandler.DeleteFileByPath(ctx, 0, req.Path); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbfilestore.DeleteFileByPathResponse{}, nil
}
