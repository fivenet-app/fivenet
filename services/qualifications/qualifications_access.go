package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) GetQualificationAccess(
	ctx context.Context,
	req *pbqualifications.GetQualificationAccessRequest,
) (*pbqualifications.GetQualificationAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		quali, err := s.getQualification(ctx, req.GetQualificationId(), nil, userInfo, false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if quali == nil || !quali.Public {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	access, err := s.access.Jobs.List(ctx, s.db, req.GetQualificationId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	for i := range access {
		s.enricher.EnrichJobInfo(access[i])
	}

	resp := &pbqualifications.GetQualificationAccessResponse{
		Access: &qualifications.QualificationAccess{
			Jobs: access,
		},
	}

	return resp, nil
}

func (s *Server) SetQualificationAccess(
	ctx context.Context,
	req *pbqualifications.SetQualificationAccessRequest,
) (*pbqualifications.SetQualificationAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if req.GetAccess() != nil {
		if _, err := s.access.HandleAccessChanges(ctx, tx, req.GetQualificationId(), req.GetAccess().GetJobs(), nil, nil); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbqualifications.SetQualificationAccessResponse{}, nil
}
