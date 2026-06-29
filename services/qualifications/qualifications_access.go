package qualifications

import (
	"context"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var qualificationSubjectAccessOptions = access.SubjectAccessOptions{
	BlockedAccess: int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_BLOCKED),
	DeniedAccessLevels: []int32{
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_REQUEST),
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_TAKE),
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	},
}

func qualificationJobAccess(
	jobs []*qualificationsaccess.QualificationJobAccess,
) *resourcesaccess.Access {
	return &resourcesaccess.Access{Jobs: jobs}
}

func normalizeQualificationJobAccess(
	userInfo *userinfo.UserInfo,
	jobs []*qualificationsaccess.QualificationJobAccess,
) (*resourcesaccess.Access, error) {
	return access.NormalizeAccess(
		qualificationJobAccess(jobs),
		nil,
		&resourcesaccess.Access{
			Jobs: []*resourcesaccess.JobAccess{{
				Job:          userInfo.GetJob(),
				MinimumGrade: userInfo.GetJobGrade(),
				Access:       int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
			}},
		},
		15,
	)
}

func (s *Server) GetQualificationAccess(
	ctx context.Context,
	req *pbqualifications.GetQualificationAccessRequest,
) (*pbqualifications.GetQualificationAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		quali, err := s.store.GetQualification(
			ctx,
			req.GetQualificationId(),
			userInfo,
			false,
			false,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if quali == nil || !quali.GetPublic() {
			return nil, errorsqualifications.ErrFailedQuery
		}
	}

	access, err := s.access.ListTargetAccess(
		ctx,
		s.db,
		req.GetQualificationId(),
		qualificationSubjectAccessOptions,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	for i := range access.GetJobs() {
		s.enricher.EnrichJobInfo(access.GetJobs()[i])
	}

	resp := &pbqualifications.GetQualificationAccessResponse{
		Access: access,
	}

	return resp, nil
}

func (s *Server) SetQualificationAccess(
	ctx context.Context,
	req *pbqualifications.SetQualificationAccessRequest,
) (*pbqualifications.SetQualificationAccessResponse, error) {
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
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
		normalizedAccess, err := normalizeQualificationJobAccess(
			userInfo,
			req.GetAccess().GetJobs(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if _, err := s.access.ReplaceTargetAccess(
			ctx,
			tx,
			s.accessResolver,
			req.GetQualificationId(),
			normalizedAccess,
			qualificationSubjectAccessOptions,
		); err != nil {
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
