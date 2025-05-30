package settings

import (
	"context"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/filestore"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const listFilesPageSize = 50

var tUserProps = table.FivenetUserProps

func (s *Server) ListFiles(ctx context.Context, req *pbsettings.ListFilesRequest) (*pbsettings.ListFilesResponse, error) {
	if req.Path != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.filestore.path", *req.Path))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbsettings.FilestoreService_ServiceDesc.ServiceName,
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
	resp := &pbsettings.ListFilesResponse{
		Pagination: pag,
	}

	files, err := s.st.List(ctx, filePath, int(req.Pagination.Offset), listFilesPageSize)
	if err != nil {
		return nil, err
	}

	fs := make([]*filestore.FileInfo, len(files))
	for i := range files {
		fs[i] = &filestore.FileInfo{
			Name:         files[i].Name,
			Size:         files[i].Size,
			LastModified: timestamp.New(files[i].LastModified),
			ContentType:  files[i].ContentType,
		}
	}
	resp.Files = fs

	resp.Pagination.Update(len(resp.Files))

	return resp, nil
}

func (s *Server) UploadFile(ctx context.Context, req *pbsettings.UploadFileRequest) (*pbsettings.UploadFileResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.filestore.prefix", req.Prefix))
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.filestore.name", req.Name))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.FilestoreService_ServiceDesc.ServiceName,
		Method:  "UploadFile",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	ctx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	if err := req.File.Upload(ctx, s.st, filestore.FilePrefix(req.Prefix), req.Name); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	objInfo, err := s.st.Stat(ctx, filestore.StripURLPrefix(*req.File.Url))
	if err != nil {
		return nil, err
	}

	return &pbsettings.UploadFileResponse{
		File: &filestore.FileInfo{
			Name:         objInfo.GetName(),
			LastModified: timestamp.New(objInfo.GetLastModified()),
			Size:         objInfo.GetSize(),
			ContentType:  objInfo.GetContentType(),
		},
	}, nil
}

func (s *Server) DeleteFile(ctx context.Context, req *pbsettings.DeleteFileRequest) (*pbsettings.DeleteFileResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.filestore.path", req.Path))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.FilestoreService_ServiceDesc.ServiceName,
		Method:  "DeleteFile",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	fullFilePath := filepath.Clean(req.Path)

	if err := s.st.Delete(ctx, fullFilePath); err != nil {
		return nil, err
	}

	resp := &pbsettings.DeleteFileResponse{}

	prefixSplit := strings.Split(fullFilePath, "/")
	if len(prefixSplit) <= 1 {
		return resp, nil
	}

	fileUrlPath := path.Join(filestore.FilestoreURLPrefix, fullFilePath)
	// Remove reference to file(s) from database for our "known file prefixes"
	switch prefixSplit[0] {
	case filestore.Avatars:
		stmt := tUserProps.
			UPDATE(
				tUserProps.Avatar,
			).
			SET(
				tUserProps.Avatar.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tUserProps.Avatar.EQ(jet.String(fileUrlPath)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}

	case filestore.MugShots:
		stmt := tUserProps.
			UPDATE(
				tUserProps.MugShot,
			).
			SET(
				tUserProps.MugShot.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tUserProps.MugShot.EQ(jet.String(fileUrlPath)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}

	case filestore.JobLogos:
		stmt := tJobProps.
			UPDATE(
				tJobProps.LogoURL,
			).
			SET(
				tJobProps.LogoURL.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tJobProps.LogoURL.EQ(jet.String(fileUrlPath)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return resp, nil
}
