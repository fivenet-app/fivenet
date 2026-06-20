package qualifications

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) ListQualificationRequests(
	ctx context.Context,
	req *pbqualifications.ListQualificationRequestsRequest,
) (*pbqualifications.ListQualificationRequestsResponse, error) {
	if req.GetQualificationId() > 0 {
		logging.InjectFields(
			ctx,
			logging.Fields{qualificationIDLogFieldKey, req.GetQualificationId()},
		)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.GetQualificationId() > 0 {
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
	}

	includePhoneNumber := false
	if fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(
		s.perms,
		userInfo,
	); err == nil {
		includePhoneNumber = fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		)
	}

	resp, err := s.store.ListQualificationRequests(
		ctx,
		qualificationsstore.ListQualificationRequestsOptions{
			Pagination:      req.GetPagination(),
			Sort:            req.GetSort(),
			QualificationID: req.GetQualificationId(),
			Status:          req.GetStatus(),
			UserIDs:         req.GetUserIds(),
		},
		userInfo,
		includePhoneNumber,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetRequests() {
		if resp.GetRequests()[i].GetUser() != nil {
			jobInfoFn(resp.GetRequests()[i].GetUser())
		}

		if resp.GetRequests()[i].GetApprover() != nil {
			jobInfoFn(resp.GetRequests()[i].GetApprover())
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationRequest(
	ctx context.Context,
	req *pbqualifications.CreateOrUpdateQualificationRequestRequest,
) (*pbqualifications.CreateOrUpdateQualificationRequestResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{
			qualificationIDLogFieldKey, req.GetRequest().GetQualificationId(),
			qualificationIDLogFieldKey, req.GetRequest().GetUserId(),
		},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	canGrade, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetRequest().GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	quali, err := s.store.GetQualification(
		ctx,
		req.GetRequest().GetQualificationId(),
		userInfo,
		false,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// If the qualification is closed and user is not a grade tutor
	if !canGrade && quali.GetClosed() {
		return nil, errorsqualifications.ErrQualificationClosed
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	defer tx.Rollback()

	// If user can grade a qualification, they are treated as an "approver" of requests
	if canGrade && req.GetRequest().GetUserId() > 0 {
		if err := s.store.ApproveQualificationRequest(
			ctx,
			tx,
			req.GetRequest(),
			userInfo,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		request, err := s.getQualificationRequest(
			ctx,
			req.GetRequest().GetQualificationId(),
			req.GetRequest().GetUserId(),
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		// Only send notification when it wasn't already in the same status
		if request == nil || request.Status == nil ||
			request.GetStatus().Enum() != req.GetRequest().GetStatus().Enum() {
			if err := s.notif.NotifyUser(ctx, &notifications.Notification{
				UserId: request.GetUserId(),
				Title: &common.I18NItem{
					Key: "notifications.qualifications.request_updated.title",
				},
				Content: &common.I18NItem{
					Key: "notifications.qualifications.request_updated.content",
					Parameters: map[string]string{
						"abbreviation": quali.GetAbbreviation(),
						"title":        quali.GetTitle(),
					},
				},
				Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
				Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
				Data: &notifications.Data{
					Link: &notifications.Link{
						To: fmt.Sprintf("/qualifications/%d", request.GetQualificationId()),
					},
				},
			}); err != nil {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	} else {
		canRequest, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetRequest().GetQualificationId(),
			userInfo,
			int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_REQUEST),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !canRequest {
			return nil, errorsqualifications.ErrFailedQuery
		}

		// Make sure the requirements of the qualification are fullfiled by the user, ErrRequirementsMissing
		reqsFullfilled, err := s.store.CheckRequirementsMetForQualification(
			ctx,
			req.GetRequest().GetQualificationId(),
			userInfo.GetUserId(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !reqsFullfilled {
			return nil, errorsqualifications.ErrRequirementsMissing
		}

		request, err := s.getQualificationRequest(
			ctx,
			req.GetRequest().GetQualificationId(),
			userInfo.GetUserId(),
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if request != nil &&
			(request.Status == nil || (request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_PENDING &&
				request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)) {
			return nil, errorsqualifications.ErrFailedQuery
		}
		if err := s.store.UpsertQualificationRequest(ctx, tx, req.GetRequest()); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	request, err := s.getQualificationRequest(
		ctx,
		req.GetRequest().GetQualificationId(),
		userInfo.GetUserId(),
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func (s *Server) getQualificationRequest(
	ctx context.Context,
	qualificationId int64,
	userId int32,
	userInfo *userinfo.UserInfo,
) (*qualifications.QualificationRequest, error) {
	request, err := s.store.GetQualificationRequest(ctx, qualificationId, userId, userInfo, false)
	if err != nil {
		return nil, err
	}

	if request.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.GetUser())
	}

	if request.GetApprover() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.GetApprover())
	}

	if !userInfo.GetSuperuser() && request.GetDeletedAt() != nil {
		return nil, errorsqualifications.ErrQualiViewDenied
	}

	return request, nil
}

func (s *Server) DeleteQualificationReq(
	ctx context.Context,
	req *pbqualifications.DeleteQualificationReqRequest,
) (*pbqualifications.DeleteQualificationReqResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	re, err := s.getQualificationRequest(ctx, req.GetQualificationId(), req.GetUserId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if re == nil {
		return &pbqualifications.DeleteQualificationReqResponse{}, nil
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		re.GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if err := s.deleteQualificationRequest(
		ctx,
		s.db,
		re.GetQualificationId(),
		re.GetUserId(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbqualifications.DeleteQualificationReqResponse{}, nil
}

func (s *Server) deleteQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	if err := s.store.DeleteQualificationRequest(ctx, tx, qualificationId, userId); err != nil {
		return err
	}

	if err := s.store.DeleteExamUser(ctx, tx, qualificationId, userId); err != nil {
		return err
	}

	return nil
}
