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
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const DocRequestMinimumWaitTime = 24 * time.Hour

var tDocRequest = table.FivenetDocumentsRequests.AS("doc_request")

func (s *Server) ListDocumentReqs(ctx context.Context, req *pbdocuments.ListDocumentReqsRequest) (*pbdocuments.ListDocumentReqsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	condition := tDocRequest.DocumentID.EQ(jet.Uint64(req.DocumentId))

	countStmt := tDocRequest.
		SELECT(
			jet.COUNT(tDocRequest.ID).AS("data_count.total"),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, ActivityDefaultPageSize)
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
			req.Pagination.Offset,
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

	resp.Pagination.Update(len(resp.Requests))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Requests {
		if resp.Requests[i].Creator != nil {
			jobInfoFn(resp.Requests[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) CreateDocumentReq(ctx context.Context, req *pbdocuments.CreateDocumentReqRequest) (*pbdocuments.CreateDocumentReqResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "CreateDocumentReq",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check && req.RequestType != documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if doc.Id <= 0 {
		doc.Id = req.DocumentId
	}

	// Owner override hatch for when a colleague isn't part of the job anymore and the document should be taken over
	if doc.CreatorJob == userInfo.Job && (doc.Creator == nil || doc.Creator.Job != doc.CreatorJob) &&
		req.RequestType == documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE {
		if err := s.updateDocumentOwner(ctx, s.db, doc.Id, userInfo, &users.UserShort{
			UserId: userInfo.UserId,
			Job:    userInfo.Job,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		accepted := true
		return &pbdocuments.CreateDocumentReqResponse{
			Request: &documents.DocRequest{
				DocumentId:  doc.Id,
				RequestType: documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
				CreatorId:   &userInfo.UserId,
				CreatorJob:  userInfo.Job,
				Reason:      req.Reason,
				Accepted:    &accepted,
			},
		}, nil
	}

	request, err := s.getDocumentReq(ctx, s.db,
		tDocRequest.DocumentID.EQ(jet.Uint64(doc.Id)).AND(
			tDocRequest.RequestType.EQ(jet.Int16(int16(req.RequestType))),
		),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if request != nil {
		// If a request of that type exists and is less than wait time old and by the same person,
		// make sure that we let the user know
		if request.CreatedAt != nil && time.Since(request.CreatedAt.AsTime()) <= DocRequestMinimumWaitTime {
			if request.CreatorId != nil && *request.CreatorId == userInfo.UserId {
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
		DocumentId:   doc.Id,
		ActivityType: req.RequestType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Reason:       req.Reason,
		Data:         req.Data,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// If no request of that type exists yet, create one, otherwise udpate the existing with the new requestors info
	if request == nil {
		request, err = s.addAndGetDocumentReq(ctx, tx, &documents.DocRequest{
			DocumentId:  doc.Id,
			CreatorId:   &userInfo.UserId,
			CreatorJob:  userInfo.Job,
			RequestType: req.RequestType,
			Reason:      req.Reason,
			Data:        req.Data,
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		accepted := false
		if err := s.updateDocumentReq(ctx, tx, request.Id, &documents.DocRequest{
			Id:          request.Id,
			DocumentId:  doc.Id,
			CreatorId:   &userInfo.UserId,
			CreatorJob:  userInfo.Job,
			RequestType: req.RequestType,
			Reason:      req.Reason,
			Data:        req.Data,
			Accepted:    &accepted,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	resp := &pbdocuments.CreateDocumentReqResponse{
		Request: request,
	}

	// If the document has no creator anymore, nothing we can do here
	if doc.CreatorId != nil {
		if err := s.notifyUserAboutRequest(ctx, doc, userInfo.UserId, int32(*doc.CreatorId)); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) UpdateDocumentReq(ctx context.Context, req *pbdocuments.UpdateDocumentReqRequest) (*pbdocuments.UpdateDocumentReqResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.id", int64(req.DocumentId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.request_id", int64(req.RequestId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "UpdateDocumentReq",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	request, err := s.getDocumentReq(ctx, s.db,
		tDocRequest.ID.EQ(jet.Uint64(req.RequestId)).
			AND(tDocRequest.DocumentID.EQ(jet.Uint64(req.DocumentId))),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if request == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	check, err := s.access.CanUserAccessTarget(ctx, request.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if (doc.CreatorId != nil && *doc.CreatorId != userInfo.UserId) && !userInfo.Superuser {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	// Skip already accepted or declined document requests
	if request.Accepted != nil && *request.Accepted {
		return nil, errorsdocuments.ErrDocReqAlreadyCompleted
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	accepted := req.Accepted
	request.Accepted = &accepted
	request.UpdatedAt = timestamp.Now()
	if err := s.updateDocumentReq(ctx, tx, request.Id, request); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Accepted the change
	if req.Accepted {
		activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED
		switch request.RequestType {
		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED

			if _, err := s.ToggleDocument(ctx, &pbdocuments.ToggleDocumentRequest{
				DocumentId: request.DocumentId,
				Closed:     true,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OPENING:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN

			if _, err := s.ToggleDocument(ctx, &pbdocuments.ToggleDocumentRequest{
				DocumentId: request.DocumentId,
				Closed:     false,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED
			// Nothing to do here, because the user is simply redirected to the editor on the frontend

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED

			if err := s.updateDocumentOwner(ctx, tx, request.DocumentId, userInfo, request.Creator); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_DELETION:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED

			if _, err := s.DeleteDocument(ctx, &pbdocuments.DeleteDocumentRequest{
				DocumentId: request.DocumentId,
				Reason:     request.Reason,
			}); err != nil {
				return nil, err
			}

		case documents.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS:
			activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED

			if request.CreatorId == nil || request.Data == nil || request.Data.GetAccessRequested() == nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			if err := s.access.Users.CreateEntry(ctx, tx, request.DocumentId, &documents.DocumentUserAccess{
				UserId:   *request.CreatorId,
				TargetId: request.DocumentId,
				Access:   request.Data.GetAccessRequested().Level,
			}); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}

		if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   request.DocumentId,
			ActivityType: activityType,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
			Reason:       req.Reason,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.UpdateDocumentReqResponse{
		Request: request,
	}, nil
}

func (s *Server) DeleteDocumentReq(ctx context.Context, req *pbdocuments.DeleteDocumentReqRequest) (*pbdocuments.DeleteDocumentReqResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.documents.request_id", int64(req.RequestId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteDocumentReq",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	request, err := s.getDocumentReq(ctx, s.db,
		tDocRequest.ID.EQ(jet.Uint64(req.RequestId)),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if request == nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	check, err := s.access.CanUserAccessTarget(ctx, request.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	if err := s.deleteDocumentReq(ctx, s.db, req.RequestId); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteDocumentReqResponse{}, nil
}

func (s *Server) notifyUserAboutRequest(ctx context.Context, doc *documents.Document, sourceUserId int32, targetUserId int32) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, targetUserId)
	if err != nil {
		return err
	}

	// Make sure target user has access to document
	check, err := s.access.CanUserAccessTarget(ctx, doc.Id, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
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
			Parameters: map[string]string{"title": doc.Title},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d#requests", doc.Id),
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

func (s *Server) addDocumentReq(ctx context.Context, tx qrm.DB, request *documents.DocRequest) (uint64, error) {
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
			request.DocumentId,
			request.RequestType,
			request.CreatorId,
			request.CreatorJob,
			request.Reason,
			request.Data,
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

	return uint64(lastId), nil
}

func (s *Server) addAndGetDocumentReq(ctx context.Context, tx qrm.DB, activitiy *documents.DocRequest) (*documents.DocRequest, error) {
	id, err := s.addDocumentReq(ctx, tx, activitiy)
	if err != nil {
		return nil, err
	}

	return s.getDocumentReqById(ctx, tx, id)
}

func (s *Server) updateDocumentReq(ctx context.Context, tx qrm.DB, id uint64, request *documents.DocRequest) error {
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
			request.DocumentId,
			request.RequestType,
			request.CreatorId,
			request.CreatorJob,
			request.Reason,
			request.Data,
			request.Accepted,
		).
		WHERE(
			tDocRequest.ID.EQ(jet.Uint64(id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}
	}

	return nil
}

func (s *Server) getDocumentReqById(ctx context.Context, tx qrm.DB, id uint64) (*documents.DocRequest, error) {
	return s.getDocumentReq(ctx, tx, tDocRequest.ID.EQ(jet.Uint64(id)))
}

func (s *Server) getDocumentReq(ctx context.Context, tx qrm.DB, condition jet.BoolExpression) (*documents.DocRequest, error) {
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

func (s *Server) deleteDocumentReq(ctx context.Context, tx qrm.DB, id uint64) error {
	tDocRequest := table.FivenetDocumentsRequests

	stmt := tDocRequest.
		DELETE().
		WHERE(tDocRequest.ID.EQ(jet.Uint64(id))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
