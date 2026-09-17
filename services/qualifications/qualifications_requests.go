package qualifications

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	citizenshydrator "github.com/fivenet-app/fivenet/v2026/stores/citizens/hydrator"
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

	resp, err := s.store.ListQualificationRequests(
		ctx,
		qualificationsstore.ListQualificationRequestsOptions{
			Pagination:      req.GetPagination(),
			Sort:            req.GetSort(),
			QualificationID: req.GetQualificationId(),
			Status:          req.GetStatus(),
			UserIDs:         req.GetUserIds(),
			Search:          req.GetSearch(),
		},
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	targets := make([]citizenshydrator.BasicTarget, 0, len(resp.GetRequests())*2)
	for i, request := range resp.GetRequests() {
		if request.GetUserId() > 0 {
			targets = append(targets, citizenshydrator.BasicTarget{
				UserID: request.GetUserId(),
				Set: func(user *usershort.UserShort) {
					resp.Requests[i].User = user
				},
			})
		}
		if request.GetApproverId() > 0 {
			targets = append(targets, citizenshydrator.BasicTarget{
				UserID: request.GetApproverId(),
				Set: func(user *usershort.UserShort) {
					resp.Requests[i].Approver = user
				},
			})
		}
	}
	if len(targets) > 0 {
		hydrateShort := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)
		if err := hydrateShort(ctx, nil, targets); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
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
			userIDLogFieldKey, req.GetRequest().GetUserId(),
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
	var publishNotification func(context.Context) error
	responseUserID := userInfo.GetUserId()

	// If user can grade a qualification, they are treated as an "approver" of requests
	if canGrade && req.GetRequest().GetUserId() > 0 {
		responseUserID = req.GetRequest().GetUserId()
		previousRequest, err := s.store.GetQualificationRequestForUpdate(
			ctx,
			tx,
			req.GetRequest().GetQualificationId(),
			req.GetRequest().GetUserId(),
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if previousRequest == nil {
			return nil, errorsqualifications.ErrQualificationRequestInvalidTransition
		}
		if !isValidQualificationRequestStatusTransition(
			previousRequest.GetStatus(),
			req.GetRequest().GetStatus(),
		) {
			return nil, errorsqualifications.ErrQualificationRequestInvalidTransition
		}
		previousStatus := previousRequest.GetStatus()
		if err := s.store.ApproveQualificationRequest(
			ctx,
			tx,
			req.GetRequest(),
			userInfo,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		requestStatus := req.GetRequest().GetStatus()
		if err := s.addQualificationActivity(
			ctx,
			tx,
			previousRequest.GetQualificationId(),
			qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_REQUEST_UPDATED,
			userInfo.GetUserId(),
			previousRequest.GetUserId(),
			qualificationsactivity.QualificationActivityData_builder{
				RequestStatus: &requestStatus,
			}.Build(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		// Only send notification when the status actually changed.
		if !req.GetSkipNotification() && previousStatus != req.GetRequest().GetStatus() &&
			previousRequest.GetUserId() != userInfo.GetUserId() {
			requestID := previousRequest.GetQualificationId()
			actorID := userInfo.GetUserId()
			entityType := "qualifications.request"
			notificationKey := "request_updated"
			switch req.GetRequest().GetStatus() {
			case qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED:
				notificationKey = "request_accepted"
			case qualifications.RequestStatus_REQUEST_STATUS_DENIED:
				notificationKey = "request_denied"
			}
			notificationData, err := s.qualificationNotificationData(
				ctx,
				previousRequest.GetUserId(),
				requestID,
			)
			if err != nil {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
			publishNotification, err = s.notif.PrepareUserNotification(
				ctx,
				tx,
				notifi.NewUserNotification(notifi.UserNotificationParams{
					UserID: previousRequest.GetUserId(),
					Title: &common.I18NItem{
						Key: "notifications.qualifications." + notificationKey + ".title",
					},
					Content: &common.I18NItem{
						Key: "notifications.qualifications." + notificationKey + ".content",
						Parameters: map[string]string{
							"abbreviation": quali.GetAbbreviation(),
							"title":        quali.GetTitle(),
						},
					},
					Category:    notifications.NotificationCategory_NOTIFICATION_CATEGORY_QUALIFICATIONS,
					Type:        notifications.NotificationType_NOTIFICATION_TYPE_INFO,
					ActorUserID: &actorID,
					EntityType:  &entityType,
					EntityID:    &requestID,
					Data:        notificationData,
				}),
			)
			if err != nil {
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

		request, err := s.store.GetQualificationRequestForUpdate(
			ctx,
			tx,
			req.GetRequest().GetQualificationId(),
			userInfo.GetUserId(),
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		// A user may submit a new request after a denial or completed attempt,
		// but an active request must not be overwritten by another submission.
		if request != nil && (!request.HasStatus() ||
			(request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_DENIED &&
				request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)) {
			return nil, errorsqualifications.ErrQualificationRequestActive
		}
		if request != nil &&
			request.GetStatus() == qualifications.RequestStatus_REQUEST_STATUS_COMPLETED &&
			quali.GetResult().GetStatus() == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL {
			return nil, errorsqualifications.ErrQualificationAlreadySuccessful
		}
		req.Request.UserId = userInfo.GetUserId()
		if err := s.store.UpsertQualificationRequest(ctx, tx, req.GetRequest()); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		requestStatus := qualifications.RequestStatus_REQUEST_STATUS_PENDING

		if err := s.addQualificationActivity(
			ctx,
			tx,
			req.GetRequest().GetQualificationId(),
			qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_REQUEST_CREATED,
			userInfo.GetUserId(),
			userInfo.GetUserId(),
			&qualificationsactivity.QualificationActivityData{
				RequestStatus: &requestStatus,
			},
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	notifi.PublishAfterCommit(ctx, s.logger, "qualification_request_updated", publishNotification)

	request, err := s.getQualificationRequest(
		ctx,
		req.GetRequest().GetQualificationId(),
		responseUserID,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func isValidQualificationRequestStatusTransition(
	previous qualifications.RequestStatus,
	current qualifications.RequestStatus,
) bool {
	switch previous {
	case qualifications.RequestStatus_REQUEST_STATUS_PENDING:
		return current == qualifications.RequestStatus_REQUEST_STATUS_PENDING ||
			current == qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED ||
			current == qualifications.RequestStatus_REQUEST_STATUS_DENIED
	case qualifications.RequestStatus_REQUEST_STATUS_DENIED:
		return current == qualifications.RequestStatus_REQUEST_STATUS_PENDING ||
			current == qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED ||
			current == qualifications.RequestStatus_REQUEST_STATUS_DENIED
	case qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED:
		return current == qualifications.RequestStatus_REQUEST_STATUS_ACCEPTED ||
			current == qualifications.RequestStatus_REQUEST_STATUS_DENIED
	default:
		return false
	}
}

func (s *Server) getQualificationRequest(
	ctx context.Context,
	qualificationId int64,
	userId int32,
	userInfo *userinfo.UserInfo,
) (*qualifications.QualificationRequest, error) {
	request, err := s.store.GetQualificationRequest(ctx, qualificationId, userId, userInfo)
	if err != nil {
		return nil, err
	}

	targets := make([]citizenshydrator.BasicTarget, 0, 2)
	if request.GetUserId() > 0 {
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: request.GetUserId(),
			Set: func(user *usershort.UserShort) {
				request.User = user
			},
		})
	}
	if request.GetApproverId() > 0 {
		targets = append(targets, citizenshydrator.BasicTarget{
			UserID: request.GetApproverId(),
			Set: func(user *usershort.UserShort) {
				request.Approver = user
			},
		})
	}
	if len(targets) > 0 {
		hydrateShort := s.hydrator.HydrateBasicTargetsSafeFunc(userInfo)
		if err := hydrateShort(ctx, nil, targets); err != nil {
			return nil, err
		}
	}

	if !userInfo.GetJobAdmin() && request.GetDeletedAt() != nil {
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.deleteQualificationRequest(
		ctx,
		tx,
		re.GetQualificationId(),
		re.GetUserId(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if err := s.addQualificationActivity(
		ctx,
		tx,
		re.GetQualificationId(),
		qualificationsactivity.QualificationActivityType_QUALIFICATION_ACTIVITY_TYPE_REQUEST_DELETED,
		userInfo.GetUserId(),
		re.GetUserId(),
		nil,
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	var publishNotification func(context.Context) error
	if !req.GetSkipNotification() && re.GetUserId() != userInfo.GetUserId() {
		qualificationID := re.GetQualificationId()
		actorID := userInfo.GetUserId()
		entityType := "qualifications.request"
		notificationData, err := s.qualificationNotificationData(
			ctx,
			re.GetUserId(),
			qualificationID,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		publishNotification, err = s.notif.PrepareUserNotification(
			ctx,
			tx,
			notifi.NewUserNotification(notifi.UserNotificationParams{
				UserID: re.GetUserId(),
				Title: &common.I18NItem{
					Key: "notifications.qualifications.request_deleted.title",
				},
				Content: &common.I18NItem{
					Key: "notifications.qualifications.request_deleted.content",
				},
				Category:    notifications.NotificationCategory_NOTIFICATION_CATEGORY_QUALIFICATIONS,
				Type:        notifications.NotificationType_NOTIFICATION_TYPE_INFO,
				ActorUserID: &actorID,
				EntityType:  &entityType,
				EntityID:    &qualificationID,
				Data:        notificationData,
			}),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	notifi.PublishAfterCommit(ctx, s.logger, "qualification_request_deleted", publishNotification)

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbqualifications.DeleteQualificationReqResponse{}, nil
}

func (s *Server) deleteQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	if err := s.softDeleteQualificationRequest(ctx, tx, qualificationId, userId); err != nil {
		return err
	}
	examUser, err := s.store.GetExamUser(ctx, tx, qualificationId, userId)
	if err != nil {
		return err
	}
	if examUser == nil {
		return nil
	}
	if err := s.deleteExamAttempt(ctx, tx, examUser.GetAttemptId()); err != nil {
		return err
	}

	return nil
}

func (s *Server) softDeleteQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	return s.store.DeleteQualificationRequest(ctx, tx, qualificationId, userId)
}
