package rector

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/filestore"
	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

const listFilesPageSize = 50

var (
	tUserProps = table.FivenetUserProps
)

func (s *Server) ListFiles(ctx context.Context, req *ListFilesRequest) (*ListFilesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Pagination.Offset <= 0 {
		defer s.aud.Log(&model.FivenetAuditLog{
			Service: RectorService_ServiceDesc.ServiceName,
			Method:  "ViewAuditLog",
			UserID:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
		}, req)
	}

	filePath := ""
	if req.Path != nil {
		filePath = filepath.Clean(*req.Path)
	}
	if filePath == "" {
		filePath = "/"
	}

	pag, _ := req.Pagination.GetResponseWithPageSize(listFilesPageSize)
	resp := &ListFilesResponse{
		Pagination: pag,
	}

	files, err := s.st.List(ctx, filePath, int(req.Pagination.Offset), listFilesPageSize)
	if err != nil {
		return nil, err
	}

	fs := make([]*filestore.FileInfo, len(files))
	for i := 0; i < len(files); i++ {
		fs[i] = &filestore.FileInfo{
			Name:         files[i].Name,
			Size:         files[i].Size,
			LastModified: timestamp.New(files[i].LastModified),
			ContentType:  files[i].ContentType,
		}
	}
	resp.Files = fs

	resp.Pagination.Update(-1, len(resp.Files))

	return resp, nil
}
func (s *Server) DeleteFile(ctx context.Context, req *DeleteFileRequest) (*DeleteFileResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorService_ServiceDesc.ServiceName,
		Method:  "DeleteFile",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	filePath := filepath.Clean(req.Path)

	if err := s.st.Delete(ctx, filePath); err != nil {
		return nil, err
	}

	resp := &DeleteFileResponse{}

	prefixSplit := strings.Split(filePath, "/")
	if len(prefixSplit) <= 1 {
		return resp, nil
	}

	prefix := prefixSplit[0]
	fileName := strings.Join(prefixSplit[1:], "/")

	// Remove reference to file(s) from database for our "known file prefixes"
	switch prefix {
	case filestore.Avatars:
		file := getFileName(fileName)
		id, err := getIdFromFile(file)
		if err != nil {
			s.logger.Error("failed to get user id from file name", zap.Error(err))
			break
		}
		if id <= 0 {
			break
		}

		stmt := tUserProps.
			UPDATE(
				tUserProps.Avatar,
			).
			SET(
				tUserProps.Avatar.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tUserProps.UserID.EQ(jet.Int32(int32(id))),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}

	case filestore.MugShots:
		file := getFileName(fileName)
		id, err := getIdFromFile(file)
		if err != nil {
			s.logger.Error("failed to get user id from file name", zap.Error(err))
			break
		}
		if id <= 0 {
			break
		}

		stmt := tUserProps.
			UPDATE(
				tUserProps.Avatar,
			).
			SET(
				tUserProps.Avatar.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tUserProps.UserID.EQ(jet.Int32(int32(id))),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}

	case filestore.JobLogos:
		job := getFileName(fileName)
		stmt := tJobProps.
			UPDATE(
				tJobProps.LogoURL,
			).
			SET(
				tJobProps.LogoURL.SET(jet.StringExp(jet.NULL)),
			).
			WHERE(
				tJobProps.Job.EQ(jet.String(job)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, err
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return resp, nil
}

func getIdFromFile(in string) (int, error) {
	parts := strings.Split(in, "-")
	return strconv.Atoi(parts[len(parts)-1])
}

func getFileName(in string) string {
	file := filepath.Base(in)
	return strings.TrimSuffix(file, filepath.Ext(file))
}
