package docstore

import (
	context "context"
	"database/sql"
	"errors"
	"strings"

	htmldiff "github.com/documize/html-diff"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	permsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/htmlsanitizer"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DocsDefaultPageSize   = 16
	DocShortContentLength = 128
)

var (
	tUsers       = table.Users
	tUserProps   = table.FivenetUserProps
	tCreator     = tUsers.AS("creator")
	tDocument    = table.FivenetDocuments.AS("document")
	tDJobAccess  = table.FivenetDocumentsJobAccess.AS("job_access")
	tDUserAccess = table.FivenetDocumentsUserAccess.AS("user_access")
)

type Server struct {
	DocStoreServiceServer

	db       *sql.DB
	ps       perms.Permissions
	cache    *mstlystcdata.Cache
	enricher *mstlystcdata.UserAwareEnricher
	aud      audit.IAuditer
	ui       userinfo.UserInfoRetriever
	notif    notifi.INotifi

	htmlDiff *htmldiff.Config
}

func NewServer(db *sql.DB, ps perms.Permissions, cache *mstlystcdata.Cache, enricher *mstlystcdata.UserAwareEnricher, aud audit.IAuditer, ui userinfo.UserInfoRetriever, notif notifi.INotifi) *Server {
	return &Server{
		db:       db,
		ps:       ps,
		cache:    cache,
		enricher: enricher,
		aud:      aud,
		ui:       ui,
		notif:    notif,
		htmlDiff: &htmldiff.Config{
			Granularity:  5,
			InsertedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-success-600"}},
			DeletedSpan:  []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-error-600"}},
			ReplacedSpan: []htmldiff.Attribute{{Key: "class", Val: "htmldiff bg-info-600"}},
			CleanTags:    []string{""},
		},
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterDocStoreServiceServer(srv, s)
}

func (s *Server) ListDocuments(ctx context.Context, req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tDocument := tDocument.AS("documentshort")
	condition := jet.Bool(true)
	if req.Search != nil && *req.Search != "" {
		condition = jet.BoolExp(
			jet.Raw("MATCH(`title`) AGAINST ($search IN BOOLEAN MODE)",
				jet.RawArgs{"$search": *req.Search}),
		)
	}
	if len(req.CategoryIds) > 0 {
		ids := make([]jet.Expression, len(req.CategoryIds))
		for i := 0; i < len(req.CategoryIds); i++ {
			ids[i] = jet.Uint64(req.CategoryIds[i])
		}

		condition = condition.AND(
			tDocument.CategoryID.IN(ids...),
		)
	}
	if len(req.CreatorIds) > 0 {
		ids := make([]jet.Expression, len(req.CreatorIds))
		for i := 0; i < len(req.CreatorIds); i++ {
			ids[i] = jet.Int32(req.CreatorIds[i])
		}

		condition = condition.AND(
			tDocument.CreatorID.IN(ids...),
		)
	}
	if req.From != nil {
		condition = condition.AND(tDocument.CreatedAt.GT_EQ(
			jet.TimestampT(req.From.AsTime()),
		))
	}
	if req.To != nil {
		condition = condition.AND(tDocument.CreatedAt.LT_EQ(
			jet.TimestampT(req.To.AsTime()),
		))
	}
	if req.Closed != nil {
		condition = condition.AND(tDocument.Closed.EQ(
			jet.Bool(*req.Closed),
		))
	}
	if len(req.DocumentIds) > 0 {
		ids := make([]jet.Expression, len(req.DocumentIds))
		for i := 0; i < len(req.DocumentIds); i++ {
			ids[i] = jet.Uint64(req.DocumentIds[i])
		}

		condition = condition.AND(
			tDocument.ID.IN(ids...),
		)
	}

	countStmt := s.listDocumentsQuery(
		condition, jet.ProjectionList{jet.COUNT(jet.DISTINCT(tDocument.ID)).AS("datacount.totalcount")}, userInfo)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, DocsDefaultPageSize)
	resp := &ListDocumentsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.listDocumentsQuery(condition, nil, userInfo).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tDocument.ID).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Documents); i++ {
		if resp.Documents[i].Creator != nil {
			jobInfoFn(resp.Documents[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Documents))

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "GetDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrDocViewDenied
	}

	infoOnly := req.InfoOnly != nil && *req.InfoOnly
	withContent := req.InfoOnly == nil || !*req.InfoOnly

	resp := &GetDocumentResponse{}
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
		docAccess, err := s.GetDocumentAccess(ctx, &GetDocumentAccessRequest{
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

func (s *Server) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "CreateDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.TemplateId != nil {
		check, err := s.checkAccessAgainstTemplate(ctx, *req.TemplateId, req.Access)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
		if !check {
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
			utils.StringFirstN(htmlsanitizer.StripTags(req.Content), DocShortContentLength),
			req.Content,
			documents.DocContentType_DOC_CONTENT_TYPE_HTML,
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

	if _, err := s.addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   uint64(lastId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if err := s.handleDocumentAccessChanges(ctx, tx, documents.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, uint64(lastId), req.Access); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrDocAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateDocumentResponse{
		DocumentId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "UpdateDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrNotFoundOrNoPerms)
	}
	var onlyUpdateAccess bool
	if !check && !userInfo.SuperUser {
		onlyUpdateAccess, err = s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_ACCESS)
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
	if !s.checkIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocUpdateDenied
	}

	if doc.TemplateId != nil {
		check, err := s.checkAccessAgainstTemplate(ctx, *doc.TemplateId, req.Access)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
		if !check {
			return nil, errorsdocstore.ErrDocRequiredAccessTemplate
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if !onlyUpdateAccess {
		doc.Content = strings.TrimSuffix(doc.Content, "<br>")
		req.Content = strings.TrimSuffix(req.Content, "<br>")

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
				utils.StringFirstN(htmlsanitizer.StripTags(req.Content), DocShortContentLength),
				req.Content,
				jet.NULL,
				req.State,
				req.Closed,
				req.Public,
			).
			WHERE(
				tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
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

		if _, err := s.addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   req.DocumentId,
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
	}

	if err := s.handleDocumentAccessChanges(ctx, tx, documents.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE, req.DocumentId, req.Access); err != nil {
		if dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrDocAccessDuplicate)
		}
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateDocumentResponse{
		DocumentId: req.DocumentId,
	}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *DeleteDocumentRequest) (*DeleteDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "DeleteDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
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

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceDeleteDocumentPerm, permsdocstore.DocStoreServiceDeleteDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
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

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteDocumentResponse{}, nil
}

func (s *Server) ToggleDocument(ctx context.Context, req *ToggleDocumentRequest) (*ToggleDocumentResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "ToggleDocument",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_STATUS)
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

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsdocstore.DocStoreServicePerm, permsdocstore.DocStoreServiceToggleDocumentPerm, permsdocstore.DocStoreServiceToggleDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}
	if !s.checkIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocToggleDenied
	}

	activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED
	if !req.Closed {
		activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN
	}

	stmt := tDocument.
		UPDATE(
			tDocument.Closed,
		).
		SET(
			req.Closed,
		).
		WHERE(
			tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if _, err := s.addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.DocumentId,
		ActivityType: activityType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &ToggleDocumentResponse{}, nil
}

func (s *Server) ChangeDocumentOwner(ctx context.Context, req *ChangeDocumentOwnerRequest) (*ChangeDocumentOwnerResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "ChangeDocumentOwner",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
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
	if !s.checkIfHasAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocstore.ErrDocOwnerFailed
	}

	tUsers := tUsers.AS("user_short")
	stmtGetUser := tUsers.
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
	if err := stmtGetUser.QueryContext(ctx, s.db, &newOwner); err != nil {
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

	return &ChangeDocumentOwnerResponse{}, nil
}
