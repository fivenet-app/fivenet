package documents

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const DocRequestMinimumWaitTime = 24 * time.Hour

var tDocRequest = table.FivenetDocumentsRequests.AS("doc_request")

func (s *Server) ListDocumentReqs(
	ctx context.Context,
	req *pbdocuments.ListDocumentReqsRequest,
) (*pbdocuments.ListDocumentReqsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	condition := tDocRequest.DocumentID.EQ(mysql.Int64(req.GetDocumentId()))

	countStmt := tDocRequest.
		SELECT(
			mysql.COUNT(tDocRequest.ID).AS("data_count.total"),
		).
		FROM(
			tDocRequest,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, ActivityDefaultPageSize)
	resp := &pbdocuments.ListDocumentReqsResponse{
		Pagination: pag,
		Requests:   []*documents.DocRequest{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tCreator := tables.User().AS("creator")

	stmt := tDocRequest.
		SELECT(
			tDocRequest.ID,
			tDocRequest.CreatedAt,
			tDocRequest.UpdatedAt,
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(
			tDocRequest.
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tDocRequest.CreatorID),
				),
		).
		WHERE(condition).
		OFFSET(
			req.GetPagination().GetOffset(),
		).
		ORDER_BY(
			tDocRequest.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetRequests() {
		if resp.GetRequests()[i].GetCreator() != nil {
			jobInfoFn(resp.GetRequests()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) CreateDocumentReq(
	ctx context.Context,
	req *pbdocuments.CreateDocumentReqRequest,
) (*pbdocuments.CreateDocumentReqResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check &&
		req.GetRequestType() != documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	doc, err := s.getDocument(
		ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if doc.GetId() <= 0 {
		doc.Id = req.GetDocumentId()
	}

	// Owner override hatch for when a colleague isn't part of the job anymore and the document should be taken over
	if doc.GetCreatorJob() == userInfo.GetJob() &&
		(doc.GetCreator() == nil || doc.GetCreator().GetJob() != doc.GetCreatorJob()) &&
		req.GetRequestType() == documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE {
		if err := s.updateDocumentOwner(ctx, s.db, doc.GetId(), userInfo, &users.UserShort{
			UserId: userInfo.GetUserId(),
			Job:    userInfo.GetJob(),
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		accepted := true
		return &pbdocuments.CreateDocumentReqResponse{
			Request: &documents.DocRequest{
				DocumentId:  doc.GetId(),
				RequestType: documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
				CreatorId:   &userInfo.UserId,
				CreatorJob:  userInfo.GetJob(),
				Reason:      req.Reason,
				Accepted:    &accepted,
			},
		}, nil
	}

	request, err := s.getDocumentReq(ctx, s.db,
		tDocRequest.DocumentID.EQ(mysql.Int64(doc.GetId())).AND(
			tDocRequest.RequestType.EQ(mysql.Int32(int32(req.GetRequestType()))),
		),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if request != nil {
		// If a request of that type exists and is less than wait time old and by the same person,
		// make sure that we let the user know
		if request.GetCreatedAt() != nil &&
			time.Since(request.GetCreatedAt().AsTime()) <= DocRequestMinimumWaitTime {
			if request.CreatorId != nil && request.GetCreatorId() == userInfo.GetUserId() {
				return nil, errorsdocuments.ErrDocReqAlreadyCreated
			}
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   doc.GetId(),
		ActivityType: req.GetRequestType(),
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Reason:       req.Reason,
		Data:         req.GetData(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// If no request of that type exists yet, create one, otherwise udpate the existing with the new requestors info
	if request == nil {
		request, err = s.addAndGetDocumentReq(ctx, tx, &documents.DocRequest{
			DocumentId:  doc.GetId(),
			CreatorId:   &userInfo.UserId,
			CreatorJob:  userInfo.GetJob(),
			RequestType: req.GetRequestType(),
			Reason:      req.Reason,
			Data:        req.GetData(),
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		accepted := false
		if err := s.updateDocumentReq(ctx, tx, request.GetId(), &documents.DocRequest{
			Id:          request.GetId(),
			DocumentId:  doc.GetId(),
			CreatorId:   &userInfo.UserId,
			CreatorJob:  userInfo.GetJob(),
			RequestType: req.GetRequestType(),
			Reason:      req.Reason,
			Data:        req.GetData(),
			Accepted:    &accepted,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	resp := &pbdocuments.CreateDocumentReqResponse{
		Request: request,
	}

	// If the document has no creator anymore, nothing we can do here
	if doc.CreatorId != nil {
		if err := s.notifyUserAboutRequest(ctx, doc, userInfo.GetUserId(), doc.GetCreatorId()); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) UpdateDocumentReq(
	ctx context.Context,
	req *pbdocuments.UpdateDocumentReqRequest,
) (*pbdocuments.UpdateDocumentReqResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.documents.id", req.GetDocumentId(),
		"fivenet.documents.request_id", req.GetRequestId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	request, err := s.getDocumentReq(ctx, s.db,
		mysql.AND(
			tDocRequest.ID.EQ(mysql.Int64(req.GetRequestId())),
			tDocRequest.DocumentID.EQ(mysql.Int64(req.GetDocumentId())),
		),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if request == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		request.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	doc, err := s.getDocument(
		ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if (doc.CreatorId != nil && doc.GetCreatorId() != userInfo.GetUserId()) &&
		!userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	// Skip already accepted or declined document requests
	if request.Accepted != nil && request.GetAccepted() {
		return nil, errorsdocuments.ErrDocReqAlreadyCompleted
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	accepted := req.GetAccepted()
	request.Accepted = &accepted
	request.UpdatedAt = timestamp.Now()
	if err := s.updateDocumentReq(ctx, tx, request.GetId(), request); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Accepted the change
	if req.GetAccepted() {
		activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED
		switch request.GetRequestType() {
		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED

			if _, err := s.ToggleDocument(ctx, &pbdocuments.ToggleDocumentRequest{
				DocumentId: request.GetDocumentId(),
				Closed:     true,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OPENING:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN

			if _, err := s.ToggleDocument(ctx, &pbdocuments.ToggleDocumentRequest{
				DocumentId: request.GetDocumentId(),
				Closed:     false,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED
			// Nothing to do here, because the user is simply redirected to the editor on the frontend

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED

			if err := s.updateDocumentOwner(ctx, tx, request.GetDocumentId(), userInfo, request.GetCreator()); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_DELETION:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED

			if _, err := s.DeleteDocument(ctx, &pbdocuments.DeleteDocumentRequest{
				DocumentId: request.GetDocumentId(),
				Reason:     request.Reason,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED

			if request.CreatorId == nil || request.GetData() == nil ||
				request.GetData().GetAccessRequested() == nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			if err := s.access.Users.CreateEntry(ctx, tx, request.GetDocumentId(), &documents.DocumentUserAccess{
				UserId:   request.GetCreatorId(),
				TargetId: request.GetDocumentId(),
				Access:   request.GetData().GetAccessRequested().GetLevel(),
			}); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}

		if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   request.GetDocumentId(),
			ActivityType: activityType,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Reason:       req.Reason,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.UpdateDocumentReqResponse{
		Request: request,
	}, nil
}

func (s *Server) DeleteDocumentReq(
	ctx context.Context,
	req *pbdocuments.DeleteDocumentReqRequest,
) (*pbdocuments.DeleteDocumentReqResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.request_id", req.GetRequestId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	request, err := s.getDocumentReq(ctx, s.db,
		tDocRequest.ID.EQ(mysql.Int64(req.GetRequestId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if request == nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		request.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	if err := s.deleteDocumentReq(ctx, s.db, req.GetRequestId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbdocuments.DeleteDocumentReqResponse{}, nil
}

func (s *Server) notifyUserAboutRequest(
	ctx context.Context,
	doc *documents.Document,
	sourceUserId int32,
	targetUserId int32,
) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, targetUserId)
	if err != nil {
		return err
	}

	// Make sure target user has access to document
	check, err := s.access.CanUserAccessTarget(
		ctx,
		doc.GetId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	not := &notifications.Notification{
		UserId: targetUserId,
		Title: &common.I18NItem{
			Key: "notifications.documents.document_request_added.title",
		},
		Content: &common.I18NItem{
			Key:        "notifications.documents.document_request_added.content",
			Parameters: map[string]string{"title": doc.GetTitle()},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d#requests", doc.GetId()),
			},
			CausedBy: &users.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	if err := s.notifi.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}

func (s *Server) addDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	request *documents.DocRequest,
) (int64, error) {
	tDocRequest := table.FivenetDocumentsRequests
	stmt := tDocRequest.
		INSERT(
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
		).
		VALUES(
			request.GetDocumentId(),
			request.GetRequestType(),
			request.GetCreatorId(),
			request.GetCreatorJob(),
			request.Reason,
			request.GetData(),
			request.Accepted,
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return 0, err
		}
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Server) addAndGetDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	activitiy *documents.DocRequest,
) (*documents.DocRequest, error) {
	id, err := s.addDocumentReq(ctx, tx, activitiy)
	if err != nil {
		return nil, err
	}

	return s.getDocumentReqById(ctx, tx, id)
}

func (s *Server) updateDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	id int64,
	request *documents.DocRequest,
) error {
	tDocRequest := table.FivenetDocumentsRequests
	stmt := tDocRequest.
		UPDATE(
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tDocRequest.Accepted,
		).
		SET(
			request.GetDocumentId(),
			request.GetRequestType(),
			request.GetCreatorId(),
			request.GetCreatorJob(),
			request.Reason,
			request.GetData(),
			request.Accepted,
		).
		WHERE(
			tDocRequest.ID.EQ(mysql.Int64(id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Server) getDocumentReqById(
	ctx context.Context,
	tx qrm.DB,
	id int64,
) (*documents.DocRequest, error) {
	return s.getDocumentReq(ctx, tx, tDocRequest.ID.EQ(mysql.Int64(id)))
}

func (s *Server) getDocumentReq(
	ctx context.Context,
	tx qrm.DB,
	condition mysql.BoolExpression,
) (*documents.DocRequest, error) {
	tCreator := tables.User().AS("creator")

	stmt := tDocRequest.
		SELECT(
			tDocRequest.ID,
			tDocRequest.CreatedAt,
			tDocRequest.UpdatedAt,
			tDocRequest.DocumentID,
			tDocRequest.RequestType,
			tDocRequest.CreatorID,
			tDocRequest.CreatorJob,
			tDocRequest.Reason,
			tDocRequest.Data,
			tCreator.ID,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Job,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(
			tDocRequest.
				INNER_JOIN(tCreator,
					tCreator.ID.EQ(tDocRequest.CreatorID),
				),
		).
		WHERE(condition).
		LIMIT(1)

	var dest documents.DocRequest
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &dest, nil
}

func (s *Server) deleteDocumentReq(ctx context.Context, tx qrm.DB, id int64) error {
	tDocRequest := table.FivenetDocumentsRequests

	stmt := tDocRequest.
		DELETE().
		WHERE(tDocRequest.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
