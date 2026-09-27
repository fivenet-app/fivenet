package jobs

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	file "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	jobAssetsNamespace = "jobassets"
	jobAssetsPageSize  = 20
)

func (s *Server) UploadJobAsset(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()
	user := auth.MustGetUserInfoFromContext(ctx)

	meta, err := s.jobAssetsFileHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	if meta.GetNamespace() != jobAssetsNamespace {
		return filestore.ErrInvalidUploadMeta
	}
	if meta.GetOriginalName() == "" {
		return status.Error(codes.InvalidArgument, "missing file name")
	}

	name := filestore.SanitizeFileName(filepath.Base(meta.GetOriginalName()))
	key := fmt.Sprintf("%s/%s/%s", jobAssetsNamespace, user.GetJob(), name)
	// Add a timestamp to avoid replacing a previous upload with the same name.
	key = fmt.Sprintf(
		"%s/%d-%s",
		filepath.Dir(key),
		time.Now().UTC().UnixNano(),
		filepath.Base(key),
	)

	resp, err := s.jobAssetsFileHandler.UploadFile(
		ctx,
		user.GetJob(),
		key,
		meta.GetSize(),
		meta.GetContentType(),
		srv,
	)
	if err != nil {
		return err
	}

	if err := s.store.UpdateJobAssetMetadata(
		ctx,
		s.db,
		user.GetJob(),
		resp.GetId(),
		user.GetUserId(),
		name,
	); err != nil {
		// Remove the committed file and association so a failed metadata update does
		// not leave an asset that cannot be managed correctly.
		if cleanupErr := s.jobAssetsFileHandler.Delete(
			ctx,
			user.GetJob(),
			resp.GetId(),
		); cleanupErr != nil {
			return errors.Join(err, cleanupErr)
		}
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	return nil
}

func (s *Server) ListJobAssets(
	ctx context.Context,
	req *pbsettings.ListJobAssetsRequest,
) (*pbsettings.ListJobAssetsResponse, error) {
	user := auth.MustGetUserInfoFromContext(ctx)

	count, err := s.store.CountJobAssets(ctx, s.db, user.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, jobAssetsPageSize)
	assets, err := s.store.ListJobAssets(ctx, s.db, user.GetJob(), pag.GetOffset(), limit)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return &pbsettings.ListJobAssetsResponse{
		Pagination: pag,
		Assets:     assets,
	}, nil
}

func (s *Server) DeleteJobAsset(
	ctx context.Context,
	req *pbsettings.DeleteJobAssetRequest,
) (*pbsettings.DeleteJobAssetResponse, error) {
	user := auth.MustGetUserInfoFromContext(ctx)

	asset, err := s.store.GetJobAsset(ctx, s.db, user.GetJob(), req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if asset == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if err := s.jobAssetsFileHandler.Delete(ctx, user.GetJob(), req.GetId()); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	return &pbsettings.DeleteJobAssetResponse{}, nil
}

func (s *Server) UpdateJobAsset(
	ctx context.Context,
	req *pbsettings.UpdateJobAssetRequest,
) (*pbsettings.UpdateJobAssetResponse, error) {
	user := auth.MustGetUserInfoFromContext(ctx)

	if asset, err := s.store.GetJobAsset(ctx, s.db, user.GetJob(), req.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	} else if asset == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	displayName := strings.TrimSpace(req.GetDisplayName())
	if displayName == "" {
		return nil, status.Error(codes.InvalidArgument, "display name must not be empty")
	}

	if err := s.store.UpdateJobAsset(
		ctx,
		s.db,
		user.GetJob(),
		req.GetId(),
		displayName,
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	asset, err := s.store.GetJobAsset(ctx, s.db, user.GetJob(), req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return &pbsettings.UpdateJobAssetResponse{Asset: asset}, nil
}

func (s *Server) ReplaceJobAsset(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()
	user := auth.MustGetUserInfoFromContext(ctx)

	meta, err := s.jobAssetsFileHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	if meta.GetNamespace() != jobAssetsNamespace || meta.GetParentId() <= 0 {
		return filestore.ErrInvalidUploadMeta
	}

	asset, err := s.store.GetJobAsset(ctx, s.db, user.GetJob(), meta.GetParentId())
	if err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if asset == nil {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	_, err = s.jobAssetsFileHandler.ReplaceFile(
		ctx,
		asset.GetId(),
		asset.GetFile().GetFilePath(),
		meta.GetSize(),
		meta.GetContentType(),
		srv,
	)
	if err != nil {
		return err
	}
	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	return nil
}
