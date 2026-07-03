package documents

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	documentsstamps "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/stamps"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
)

const stampLimit = 5

var stampSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
	},
}

func (s *Server) ListUsableStamps(
	ctx context.Context,
	req *pbdocuments.ListUsableStampsRequest,
) (*pbdocuments.ListUsableStampsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	pag, stamps, err := s.store.ListUsableStamps(ctx, documentsstore.ListUsableStampsQuery{
		Pagination: req.GetPagination(),
		UserInfo:   userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.ListUsableStampsResponse{
		Pagination: pag,
		Stamps:     stamps,
	}, nil
}

func (s *Server) GetStamp(
	ctx context.Context,
	req *pbdocuments.GetStampRequest,
) (*pbdocuments.GetStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.signingStampAccess.CanUserAccessTarget(
		ctx,
		req.GetId(),
		userInfo,
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_USE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	stamp, err := s.getStamp(ctx, req.GetId(), true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.GetStampResponse{Stamp: stamp}, nil
}

func (s *Server) getStamp(
	ctx context.Context,
	stampID int64,
	withAccess bool,
) (*documentsstamps.Stamp, error) {
	stamp, err := s.store.GetStamp(ctx, stampID)
	if err != nil {
		return nil, err
	}
	if stamp == nil {
		return nil, nil
	}

	if withAccess {
		access, err := s.signingStampAccess.ListTargetAccess(
			ctx,
			s.db,
			stampID,
			stampSubjectAccessOptions,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		stamp.Access = access
	}

	return stamp, nil
}

func (s *Server) UpsertStamp(
	ctx context.Context,
	req *pbdocuments.UpsertStampRequest,
) (*pbdocuments.UpsertStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	st := req.GetStamp()

	// Stamps are job only and are currently limited to 5!
	if st.GetAccess() == nil {
		st.Access = &documentsstamps.StampAccess{}
	}
	if len(st.GetAccess().GetJobs()) == 0 {
		st.Access.Jobs = append(st.Access.Jobs, &documentsstamps.StampJobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
		})
	}

	fallbackAccess := &documentsstamps.StampAccess{
		Jobs: []*documentsstamps.StampJobAccess{{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
		}},
	}

	normalizedAccess, err := access.NormalizeAccess(st.GetAccess(), nil, fallbackAccess, 15)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	st.Access = normalizedAccess

	if count, err := s.store.CheckJobStampCount(ctx, userInfo.GetJob()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	} else if count >= stampLimit && st.GetId() == 0 {
		return nil, errorsdocuments.ErrStampLimitReached
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if st.GetId() != 0 {
		check, err := s.signingStampAccess.CanUserAccessTarget(
			ctx,
			st.GetId(),
			userInfo,
			int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
		)
		if err != nil {
			return nil, err
		}
		if !check {
			return nil, errorsdocuments.ErrPermissionDenied
		}

		if err := s.store.UpdateStamp(ctx, tx, st); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		lastID, err := s.store.CreateStamp(ctx, tx, st)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		st.SetId(lastID)
	}

	if _, err := s.signingStampAccess.ReplaceTargetAccess(
		ctx,
		tx,
		s.subjectResolver,
		st.GetId(),
		normalizedAccess,
		stampSubjectAccessOptions,
	); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	stamp, err := s.getStamp(ctx, st.GetId(), true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.UpsertStampResponse{Stamp: stamp}, nil
}

func (s *Server) DeleteStamp(
	ctx context.Context,
	req *pbdocuments.DeleteStampRequest,
) (*pbdocuments.DeleteStampResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.signingStampAccess.CanUserAccessTarget(
		ctx,
		req.GetStampId(),
		userInfo,
		int32(documentsstamps.StampAccessLevel_STAMP_ACCESS_LEVEL_MANAGE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	stamp, err := s.store.GetStamp(ctx, req.GetStampId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	var deletedAtTime *timestamp.Timestamp
	if stamp.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	if err := s.store.DeleteStamp(ctx, s.db, req.GetStampId(), deletedAtTime); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteStampResponse{}, nil
}
