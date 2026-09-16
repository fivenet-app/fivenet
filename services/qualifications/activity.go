package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) addQualificationActivity(
	ctx context.Context,
	tx qrm.DB,
	qualificationID int64,
	activityType qualificationsactivity.QualificationActivityType,
	actorUserID, targetUserID int32,
	data *qualificationsactivity.QualificationActivityData,
) error {
	return s.addQualificationActivityForAttempt(
		ctx, tx, qualificationID, activityType, actorUserID, targetUserID, "", data,
	)
}

func (s *Server) addQualificationActivityForAttempt(
	ctx context.Context,
	tx qrm.DB,
	qualificationID int64,
	activityType qualificationsactivity.QualificationActivityType,
	actorUserID, targetUserID int32,
	attemptID string,
	data *qualificationsactivity.QualificationActivityData,
) error {
	return s.store.CreateQualificationActivity(
		ctx,
		tx,
		&qualificationsactivity.QualificationActivity{
			QualificationId: qualificationID,
			Type:            activityType,
			ActorUserId:     int32PtrOrNil(actorUserID),
			TargetUserId:    int32PtrOrNil(targetUserID),
			AttemptId:       attemptID,
			Data:            data,
		},
	)
}

func int32PtrOrNil(value int32) *int32 {
	if value == 0 {
		return nil
	}
	return &value
}

func (s *Server) ListQualificationActivity(
	ctx context.Context,
	req *pbqualifications.ListQualificationActivityRequest,
) (*pbqualifications.ListQualificationActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	opts := qualificationsstore.ListQualificationActivityOptions{
		QualificationID: req.GetQualificationId(), Types: req.GetTypes(), UserID: req.GetUserId(),
		From: req.GetFrom(), To: req.GetTo(), Sort: req.GetSort(),
	}
	count, err := s.store.CountQualificationActivity(ctx, opts)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	pagination, limit := req.GetPagination().GetResponseWithPageSize(count, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationActivityResponse{
		Pagination: pagination,
		Activity:   []*qualificationsactivity.QualificationActivity{},
	}
	if count == 0 {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
		return resp, nil
	}
	opts.Offset, opts.Limit = pagination.GetOffset(), limit
	resp.Activity, err = s.store.ListQualificationActivity(ctx, opts)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetActivity())*2)
	for i, activity := range resp.GetActivity() {
		// Attempt IDs are internal correlation keys and must not be exposed to
		// clients, even though they are retained on the stored activity.
		activity.SetAttemptId("")
		if activity.GetActorUserId() > 0 {
			targets = append(
				targets,
				citizenshydrator.BasicTarget{
					UserID: activity.GetActorUserId(),
					Set:    func(user *usershort.UserShort) { resp.Activity[i].ActorUser = user },
				},
			)
		}
		if activity.GetTargetUserId() > 0 {
			targets = append(
				targets,
				citizenshydrator.BasicTarget{
					UserID: activity.GetTargetUserId(),
					Set:    func(user *usershort.UserShort) { resp.Activity[i].TargetUser = user },
				},
			)
		}
	}
	if len(targets) > 0 {
		if err := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)(ctx, nil, targets); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)
	return resp, nil
}
