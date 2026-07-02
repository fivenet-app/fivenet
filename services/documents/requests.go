package documents

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsrequests "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/requests"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const DocRequestMinimumWaitTime = 24 * time.Hour

var tDocRequest = table.FivenetDocumentsRequests.AS("doc_request")

func (s *Server) ListDocumentReqs(
	ctx context.Context,
	req *pbdocuments.ListDocumentReqsRequest,
) (*pbdocuments.ListDocumentReqsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	count, reqs, err := s.store.ListDocumentReqs(ctx, documentsstore.ListDocumentReqsQuery{
		DocumentID: req.GetDocumentId(),
		Pagination: req.GetPagination(),
		UserInfo:   userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pag, _ := req.GetPagination().GetResponseWithPageSize(count.Total, ActivityDefaultPageSize)
	resp := &pbdocuments.ListDocumentReqsResponse{
		Pagination: pag,
		Requests:   []*documentsrequests.DocRequest{},
	}
	if count.Total <= 0 {
		return resp, nil
	}
	resp.Requests = reqs

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
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check &&
		req.GetRequestType() != documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS {
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
	if doc == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}
	if doc.GetId() <= 0 {
		doc.Id = req.GetDocumentId()
	}

	// Owner override hatch for when a colleague isn't part of the job anymore and the document should be taken over
	if doc.GetCreatorJob() == userInfo.GetJob() &&
		(doc.GetCreator() == nil || doc.GetCreator().GetJob() != doc.GetCreatorJob()) &&
		req.GetRequestType() == documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE {
		if err := s.store.UpdateDocumentOwner(
			ctx,
			s.db,
			doc.GetId(),
			userInfo,
			&usershort.UserShort{
				UserId: userInfo.GetUserId(),
				Job:    userInfo.GetJob(),
			},
		); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		accepted := true
		return &pbdocuments.CreateDocumentReqResponse{
			Request: &documentsrequests.DocRequest{
				DocumentId:  doc.GetId(),
				RequestType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED,
				CreatorId:   &userInfo.UserId,
				CreatorJob:  userInfo.GetJob(),
				Reason:      req.Reason,
				Accepted:    &accepted,
			},
		}, nil
	}

	request, err := s.store.GetDocumentReq(ctx, s.db,
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

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
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
		requestId, err := s.store.AddDocumentReq(ctx, tx, &documentsrequests.DocRequest{
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
		request, err = s.store.GetDocumentReq(ctx, tx, tDocRequest.ID.EQ(mysql.Int64(requestId)))
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		accepted := false
		if err := s.store.UpdateDocumentReq(ctx, tx, request.GetId(), &documentsrequests.DocRequest{
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
		if err := s.notifyUserAboutRequest(
			ctx,
			doc,
			userInfo.GetUserId(),
			doc.GetCreatorId(),
		); err != nil {
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
		documentIDLogFieldKey, req.GetDocumentId(),
		"fivenet.documents.request_id", req.GetRequestId(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	request, err := s.store.GetDocumentReq(ctx, s.db,
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

	check, err := s.canUserAccessDocument(
		ctx,
		request.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
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
	if doc == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	if (doc.GetCreatorId() > 0 && doc.GetCreatorId() != userInfo.GetUserId()) &&
		!userInfo.GetJobAdmin() {
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
	if err := s.store.UpdateDocumentReq(ctx, tx, request.GetId(), request); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Accepted the change
	if req.GetAccepted() {
		activityType := documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED
		switch request.GetRequestType() {
		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_CLOSURE:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED

			if err := s.store.ToggleDocument(
				ctx,
				tx,
				request.GetDocumentId(),
				doc.GetTemplateId(),
				true,
				userInfo,
			); err != nil {
				return nil, err
			}

		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OPENING:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN

			if err := s.store.ToggleDocument(
				ctx,
				tx,
				request.GetDocumentId(),
				doc.GetTemplateId(),
				false,
				userInfo,
			); err != nil {
				return nil, err
			}

		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_UPDATE:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED
			// Nothing to do here, because the user is simply redirected to the editor on the frontend

		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_OWNER_CHANGE:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_OWNER_CHANGED

			if err := s.store.UpdateDocumentOwner(
				ctx,
				tx,
				request.GetDocumentId(),
				userInfo,
				request.GetCreator(),
			); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_DELETION:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_DELETED

			if _, err := s.DeleteDocument(ctx, &pbdocuments.DeleteDocumentRequest{
				DocumentId: request.GetDocumentId(),
				Reason:     request.Reason,
			}); err != nil {
				return nil, err
			}

		case documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_REQUESTED_ACCESS:
			activityType = documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_ACCESS_UPDATED

			if request.CreatorId == nil || request.GetData() == nil ||
				request.GetData().GetAccessRequested() == nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			if err := s.subjectAccess.CreateUserAccess(
				ctx,
				tx,
				s.subjectResolver,
				request.GetDocumentId(),
				request.GetCreatorId(),
				int32(request.GetData().GetAccessRequested().GetLevel()),
			); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}

		if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
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

	request, err := s.store.GetDocumentReq(ctx, s.db,
		tDocRequest.ID.EQ(mysql.Int64(req.GetRequestId())),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if request == nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrDocReqAlreadyCompleted)
	}

	check, err := s.canUserAccessDocument(
		ctx,
		request.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	if err := s.store.DeleteDocumentReq(ctx, s.db, req.GetRequestId()); err != nil {
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
	userInfo, err := s.ui.GetUserInfo(ctx, targetUserId)
	if err != nil {
		return err
	}

	// Make sure target user has access to document
	check, err := s.canUserAccessDocument(
		ctx,
		doc.GetId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
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
			CausedBy: &usershort.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	if err := s.notifi.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
