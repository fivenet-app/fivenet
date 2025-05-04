package docstore

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore"
	permsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/perms"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorsdocstore "github.com/fivenet-app/fivenet/services/docstore/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DocsDefaultPageSize = 16
	DocSummaryLength    = 128
)

var (
	tUserProps     = table.FivenetUserProps
	tDocument      = table.FivenetDocuments.AS("document")
	tDocumentShort = table.FivenetDocuments.AS("documentshort")
	tDAccess       = table.FivenetDocumentsAccess.AS("job_access")
)

func (s *Server) ListDocuments(ctx context.Context, req *pbdocstore.ListDocumentsRequest) (*pbdocstore.ListDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := false

	condition := jet.Bool(true)
	if req.Search != nil && *req.Search != "" {
		logRequest = true
		condition = jet.BoolExp(
			jet.Raw(
				"MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{
					"$search": *req.Search,
				},
			),
		)
	}
	if len(req.CategoryIds) > 0 {
		ids := make([]jet.Expression, len(req.CategoryIds))
		for i := range req.CategoryIds {
			ids[i] = jet.Uint64(req.CategoryIds[i])
		}

		condition = condition.AND(
			tDocumentShort.CategoryID.IN(ids...),
		)
	}
	if len(req.CreatorIds) > 0 {
		logRequest = true
		ids := make([]jet.Expression, len(req.CreatorIds))
		for i := range req.CreatorIds {
			ids[i] = jet.Int32(req.CreatorIds[i])
		}

		condition = condition.AND(
			tDocumentShort.CreatorID.IN(ids...),
		)
	}
	if req.From != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.GT_EQ(
			jet.TimestampT(req.From.AsTime()),
		))
	}
	if req.To != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.LT_EQ(
			jet.TimestampT(req.To.AsTime()),
		))
	}
	if req.Closed != nil {
		condition = condition.AND(tDocumentShort.Closed.EQ(
			jet.Bool(*req.Closed),
		))
	}
	if len(req.DocumentIds) > 0 {
		ids := make([]jet.Expression, len(req.DocumentIds))
		for i := range req.DocumentIds {
			ids[i] = jet.Uint64(req.DocumentIds[i])
		}

		condition = condition.AND(
			tDocumentShort.ID.IN(ids...),
		)
	}

	if logRequest {
		defer s.aud.Log(&model.FivenetAuditLog{
			Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
			Method:  "ListDocuments",
			UserID:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
		}, req)
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "title":
			column = tDocumentShort.Title
		case "createdAt":
			fallthrough
		default:
			column = tDocumentShort.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys,
				column.ASC(),
				tDocumentShort.UpdatedAt.DESC(),
			)
		} else {
			orderBys = append(orderBys,
				column.DESC(),
				tDocumentShort.UpdatedAt.DESC(),
			)
		}
	} else {
		orderBys = append(orderBys, tDocumentShort.UpdatedAt.DESC())
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(database.NoTotalCount, DocsDefaultPageSize)
	resp := &pbdocstore.ListDocumentsResponse{
		Pagination: pag,
	}

	stmt := s.listDocumentsQuery(condition, nil, userInfo).
		ORDER_BY(orderBys...).
		GROUP_BY(tDocumentShort.ID).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Documents {
		if resp.Documents[i].Creator != nil {
			jobInfoFn(resp.Documents[i].Creator)
		}

		if job := s.enricher.GetJobByName(resp.Documents[i].CreatorJob); job != nil {
			resp.Documents[i].CreatorJobLabel = &job.Label
		}
	}

	resp.Pagination.Update(len(resp.Documents))

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *pbdocstore.GetDocumentRequest) (*pbdocstore.GetDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "GetDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrDocViewDenied
	}

	infoOnly := req.InfoOnly != nil && *req.InfoOnly
	withContent := req.InfoOnly == nil || !*req.InfoOnly

	resp := &pbdocstore.GetDocumentResponse{}
	resp.Document, err = s.getDocument(ctx,
		tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, withContent)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if resp.Document == nil || resp.Document.Id <= 0 {
		return nil, errorsdocstore.ErrNotFoundOrNoPerms
	}

	if resp.Document.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.Document.Creator)
	}

	if !infoOnly {
		docAccess, err := s.GetDocumentAccess(ctx, &pbdocstore.GetDocumentAccessRequest{
			DocumentId: resp.Document.Id,
		})
		if err != nil {
			if st, ok := status.FromError(err); !ok {
				return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
			} else {
				// Ignore permission denied as we are simply getting the document
				if st.Code() != codes.PermissionDenied {
					return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
				}
			}
		}
		if docAccess != nil {
			resp.Access = docAccess.Access
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) getDocument(ctx context.Context, condition jet.BoolExpression, userInfo *userinfo.UserInfo, withContent bool) (*documents.Document, error) {
	var doc documents.Document

	stmt := s.getDocumentQuery(condition, nil, userInfo, withContent).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &doc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	if doc.Creator != nil {
		s.enricher.EnrichJobInfo(doc.Creator)
	}

	return &doc, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *pbdocstore.CreateDocumentRequest) (*pbdocstore.CreateDocumentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	var tmpl *documents.Template
	if req.TemplateId != nil {
		var err error
		tmpl, err = s.getTemplate(ctx, *req.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		if !s.checkAccessAgainstTemplate(tmpl, req.Access) {
			return nil, errorsdocstore.ErrDocRequiredAccessTemplate
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDocument := table.FivenetDocuments
	stmt := tDocument.
		INSERT(
			tDocument.CategoryID,
			tDocument.Title,
			tDocument.Summary,
			tDocument.Content,
			tDocument.ContentType,
			tDocument.Data,
			tDocument.CreatorID,
			tDocument.CreatorJob,
			tDocument.State,
			tDocument.Closed,
			tDocument.Public,
			tDocument.TemplateID,
		).
		VALUES(
			req.CategoryId,
			req.Title,
			req.Content.GetSummary(DocSummaryLength),
			req.Content,
			content.ContentType_CONTENT_TYPE_HTML,
			req.Data,
			userInfo.UserId,
			userInfo.Job,
			req.State,
			req.Closed,
			req.Public,
			req.TemplateId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   uint64(lastId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if err := s.handleDocumentAccessChange(ctx, tx, uint64(lastId), userInfo, req.Access, false); err != nil {
		return nil, err
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, uint64(lastId), tmpl.Workflow); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &pbdocstore.CreateDocumentResponse{
		DocumentId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *pbdocstore.UpdateDocumentRequest) (*pbdocstore.UpdateDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	var onlyUpdateAccess bool
	if !check && !userInfo.SuperUser {
		onlyUpdateAccess, err = s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_ACCESS)
		if err != nil {
			return nil, errorsdocstore.ErrPermissionDenied
		}
		if !onlyUpdateAccess {
			return nil, errorsdocstore.ErrPermissionDenied
		}
	}

	doc, err := s.getDocument(ctx,
		tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Either the document is closed and the update request isn't re-opening the document
	if doc.Closed && req.Closed && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrClosedDoc
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceUpdateDocumentPerm, permsdocstore.DocStoreServiceUpdateDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !access.CheckIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocUpdateDenied
	}

	var tmpl *documents.Template
	if doc.TemplateId != nil {
		var err error
		tmpl, err = s.getTemplate(ctx, *doc.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		if !s.checkAccessAgainstTemplate(tmpl, req.Access) {
			return nil, errorsdocstore.ErrDocRequiredAccessTemplate
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if !onlyUpdateAccess {
		tDocument := table.FivenetDocuments
		stmt := tDocument.
			UPDATE(
				tDocument.CategoryID,
				tDocument.Title,
				tDocument.Summary,
				tDocument.Content,
				tDocument.Data,
				tDocument.State,
				tDocument.Closed,
				tDocument.Public,
			).
			SET(
				req.CategoryId,
				req.Title,
				req.Content.GetSummary(DocSummaryLength),
				req.Content,
				jet.NULL,
				req.State,
				req.Closed,
				req.Public,
			).
			WHERE(
				tDocument.ID.EQ(jet.Uint64(doc.Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		diff, err := s.generateDocumentDiff(doc, &documents.Document{
			Title:   req.Title,
			Content: req.Content,
			State:   req.State,
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   doc.Id,
			ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
			Data: &documents.DocActivityData{
				Data: &documents.DocActivityData_Updated{
					Updated: diff,
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		if tmpl != nil {
			if err := s.createOrUpdateWorkflowState(ctx, tx, doc.Id, tmpl.Workflow); err != nil {
				return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
			}
		}
	}

	if err := s.handleDocumentAccessChange(ctx, tx, doc.Id, userInfo, req.Access, true); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbdocstore.UpdateDocumentResponse{
		DocumentId: doc.Id,
	}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *pbdocstore.DeleteDocumentRequest) (*pbdocstore.DeleteDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsdocstore.ErrDocDeleteDenied
		}
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Require a reason if the document is not already deleted
	if doc.DeletedAt == nil && req.Reason == nil {
		return nil, errorsdocstore.ErrDocDeleteDenied
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceDeleteDocumentPerm, permsdocstore.DocStoreServiceDeleteDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !access.CheckIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocDeleteDenied
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if doc.DeletedAt != nil && userInfo.SuperUser {
		deletedAtTime = jet.TimestampExp(jet.NULL)
	}

	stmt := tDocument.
		UPDATE(
			tDocument.DeletedAt,
		).
		SET(
			tDocument.DeletedAt.SET(deletedAtTime),
		).
		WHERE(
			tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Reason:       req.Reason,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbdocstore.DeleteDocumentResponse{}, nil
}

func (s *Server) ToggleDocument(ctx context.Context, req *pbdocstore.ToggleDocumentRequest) (*pbdocstore.ToggleDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "ToggleDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_STATUS)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		if !userInfo.SuperUser {
			return nil, errorsdocstore.ErrDocToggleDenied
		}
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	var tmpl *documents.Template
	if !req.Closed && doc.TemplateId != nil { // If the document is opened, get template so we can update the reminder/auto close times
		tmpl, err = s.getTemplate(ctx, *doc.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceToggleDocumentPerm, permsdocstore.DocStoreServiceToggleDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !access.CheckIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocToggleDenied
	}

	activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED
	if !req.Closed {
		activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	stmt := tDocument.
		UPDATE(
			tDocument.Closed,
		).
		SET(
			req.Closed,
		).
		WHERE(
			tDocument.ID.EQ(jet.Uint64(doc.Id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   doc.Id,
		ActivityType: activityType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, doc.Id, tmpl.Workflow); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbdocstore.ToggleDocumentResponse{}, nil
}

func (s *Server) ChangeDocumentOwner(ctx context.Context, req *pbdocstore.ChangeDocumentOwnerRequest) (*pbdocstore.ChangeDocumentOwnerResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbdocstore.DocStoreService_ServiceDesc.ServiceName,
		Method:  "ChangeDocumentOwner",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrDocOwnerFailed
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Document must be created by the same job
	if doc.CreatorJob != userInfo.Job {
		return nil, errorsdocstore.ErrDocOwnerWrongJob
	}

	// If user is not a super user make sure they can only change owner to themselves
	if req.NewUserId == nil || !userInfo.SuperUser {
		req.NewUserId = &userInfo.UserId
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceChangeDocumentOwnerPerm, permsdocstore.DocStoreServiceChangeDocumentOwnerAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !access.CheckIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocOwnerFailed
	}

	tUsers := tables.Users().AS("user_short")

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.Dateofbirth,
		).
		FROM(tUsers).
		WHERE(tUsers.ID.EQ(jet.Int32(*req.NewUserId))).
		LIMIT(1)

	var newOwner users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &newOwner); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if newOwner.UserId <= 0 {
		return nil, errorsdocstore.ErrFailedQuery
	}

	// Allow super users to transfer documents cross jobs
	if !userInfo.SuperUser {
		if newOwner.Job != doc.CreatorJob {
			return nil, errorsdocstore.ErrDocOwnerWrongJob
		}

		if doc.CreatorId != nil && *doc.CreatorId == userInfo.UserId {
			return nil, errorsdocstore.ErrDocSameOwner
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.updateDocumentOwner(ctx, tx, req.DocumentId, userInfo, &newOwner); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbdocstore.ChangeDocumentOwnerResponse{}, nil
}
