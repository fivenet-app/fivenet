package documents

import (
	context "context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
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
	tDocumentShort = table.FivenetDocuments.AS("document_short")
	tDAccess       = table.FivenetDocumentsAccess.AS("job_access")
)

func (s *Server) ListDocuments(ctx context.Context, req *pbdocuments.ListDocumentsRequest) (*pbdocuments.ListDocumentsResponse, error) {
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
	if req.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(jet.Bool(*req.OnlyDrafts)))
	}

	if logRequest {
		defer s.aud.Log(&audit.AuditEntry{
			Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
			Method:  "ListDocuments",
			UserId:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   audit.EventType_EVENT_TYPE_VIEWED,
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
	resp := &pbdocuments.ListDocumentsResponse{
		Pagination: pag,
	}

	stmt := s.listDocumentsQuery(condition, nil, nil, userInfo).
		ORDER_BY(orderBys...).
		GROUP_BY(tDocumentShort.ID).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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

func (s *Server) GetDocument(ctx context.Context, req *pbdocuments.GetDocumentRequest) (*pbdocuments.GetDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "GetDocument",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	infoOnly := req.InfoOnly != nil && *req.InfoOnly
	withContent := req.InfoOnly == nil || !*req.InfoOnly

	resp := &pbdocuments.GetDocumentResponse{}
	resp.Document, err = s.getDocument(ctx,
		tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, withContent)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.Document == nil || resp.Document.Id <= 0 {
		return nil, errorsdocuments.ErrNotFoundOrNoPerms
	}

	if resp.Document.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.Document.Creator)
	}

	resp.Document.Pin, err = s.getDocumentPin(ctx, resp.Document.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if !infoOnly {
		docAccess, err := s.GetDocumentAccess(ctx, &pbdocuments.GetDocumentAccessRequest{
			DocumentId: resp.Document.Id,
		})
		if err != nil {
			if st, ok := status.FromError(err); !ok {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			} else {
				// Ignore permission denied as we are simply getting the document
				if st.Code() != codes.PermissionDenied {
					return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
				}
			}
		}
		if docAccess != nil {
			resp.Access = docAccess.Access
		}

		files, err := s.fHandler.ListFilesForParentID(ctx, resp.Document.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.Document.Files = files
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) getDocument(ctx context.Context, condition jet.BoolExpression, userInfo *userinfo.UserInfo, withContent bool) (*documents.Document, error) {
	var doc documents.Document

	stmt := s.getDocumentQuery(condition, nil, userInfo, withContent).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &doc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if doc.Creator != nil {
		s.enricher.EnrichJobInfo(doc.Creator)
	}

	return &doc, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *pbdocuments.CreateDocumentRequest) (*pbdocuments.CreateDocumentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "CreateDocument",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	var docContent string
	var docTitle string
	var docState string
	var categoryId *uint64
	docAccess := &documents.DocumentAccess{}
	docReferences := []*documents.DocumentReference{}
	docRelations := []*documents.DocumentRelation{}

	var tmpl *documents.Template
	if req.TemplateId != nil {
		var err error
		tmpl, err = s.getTemplate(ctx, *req.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		docTitle, docState, docContent, err = s.renderTemplate(tmpl, req.TemplateData)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		// Set access based on template
		docAccess = &documents.DocumentAccess{
			Jobs:  tmpl.ContentAccess.Jobs,
			Users: tmpl.ContentAccess.Users,
		}

		if tmpl.Category != nil {
			categoryId = &tmpl.Category.Id
		}

		// Add references from template data documents if not already present
		if tmpl != nil && req.TemplateData != nil {
			if len(req.TemplateData.Documents) > 0 {
				for _, doc := range req.TemplateData.Documents {
					exists := false
					for _, reference := range docReferences {
						if reference.TargetDocumentId == doc.Id {
							exists = true
							break
						}
					}

					if !exists {
						docReferences = append(docReferences, &documents.DocumentReference{
							// Id will be assigned by backend or can be zero for new
							SourceDocumentId: 0, // will be set after insert
							TargetDocumentId: doc.Id,
							// TargetDocument can be set if needed, or left nil
							CreatorId: &userInfo.UserId,
							// Creator can be set if needed, or left nil
							Reference: documents.DocReference_DOC_REFERENCE_SOLVES,
						})
					}
				}
			}

			// Add relations from template data users if not already present
			if len(req.TemplateData.Users) > 0 {
				for _, user := range req.TemplateData.Users {
					exists := false
					for _, relation := range docRelations {
						if relation.TargetUserId == user.UserId {
							exists = true
							break
						}
					}

					if !exists {
						docRelations = append(docRelations, &documents.DocumentRelation{
							// Id will be assigned by backend or can be zero for new
							DocumentId:   0, // will be set after insert
							TargetUserId: user.UserId,
							// TargetUser can be set if needed, or left nil
							SourceUserId: userInfo.UserId,
							// SourceUser can be set if needed, or left nil
							Relation: documents.DocRelation_DOC_RELATION_CAUSED,
						})
					}
				}
			}
		}
	} else {
		// Add minimum access for the creator's job
		docAccess.Jobs = append(docAccess.Jobs, &documents.DocumentJobAccess{
			Job:          userInfo.Job,
			MinimumGrade: userInfo.JobGrade,
			Access:       documents.AccessLevel_ACCESS_LEVEL_EDIT,
		})
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDocument := table.FivenetDocuments
	stmt := tDocument.
		INSERT(
			tDocument.Title,
			tDocument.Summary,
			tDocument.CategoryID,
			tDocument.Content,
			tDocument.ContentType,
			tDocument.State,
			tDocument.Data,
			tDocument.Closed,
			tDocument.Draft,
			tDocument.Public,
			tDocument.TemplateID,
			tDocument.CreatorID,
			tDocument.CreatorJob,
		).
		VALUES(
			docTitle,
			content.GetSummary(docContent, DocSummaryLength),
			categoryId,
			docContent,
			req.ContentType,
			docState,
			jet.NULL,
			false,
			true,
			false,
			req.TemplateId,
			userInfo.UserId,
			userInfo.Job,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   uint64(lastId),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.handleDocumentAccessChange(ctx, tx, uint64(lastId), userInfo, docAccess, false); err != nil {
		return nil, err
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, uint64(lastId), tmpl.Workflow); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	for _, ref := range docReferences {
		ref.SourceDocumentId = uint64(lastId)
		if _, err := s.addDocumentReference(ctx, tx, ref); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}
	for _, rel := range docRelations {
		rel.DocumentId = uint64(lastId)
		if _, err := s.addDocumentRelation(ctx, tx, userInfo, rel); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbdocuments.CreateDocumentResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *pbdocuments.UpdateDocumentRequest) (*pbdocuments.UpdateDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "UpdateDocument",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	var onlyUpdateAccess bool
	if !check && !userInfo.Superuser {
		onlyUpdateAccess, err = s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_ACCESS)
		if err != nil {
			return nil, errorsdocuments.ErrPermissionDenied
		}
		if !onlyUpdateAccess {
			return nil, errorsdocuments.ErrPermissionDenied
		}
	}

	oldDoc, err := s.getDocument(ctx,
		tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Either the document is closed and the update request isn't re-opening the document
	if oldDoc.Closed && req.Closed && !userInfo.Superuser {
		return nil, errorsdocuments.ErrClosedDoc
	}

	// A document can only be switched to published once
	if !oldDoc.Draft && oldDoc.Draft != req.Draft {
		// Allow a super user to change the draft state
		if !userInfo.Superuser {
			req.Draft = oldDoc.Draft
		}
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceUpdateDocumentPerm, permsdocuments.DocumentsServiceUpdateDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, oldDoc.CreatorJob, oldDoc.Creator) {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	var tmpl *documents.Template
	if oldDoc.TemplateId != nil {
		var err error
		tmpl, err = s.getTemplate(ctx, *oldDoc.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if !s.checkAccessAgainstTemplate(tmpl, req.Access) {
			return nil, errorsdocuments.ErrDocRequiredAccessTemplate
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
				tDocument.Draft,
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
				req.Draft,
				req.Public,
			).
			WHERE(
				tDocument.ID.EQ(jet.Uint64(oldDoc.Id)),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		diff, err := s.generateDocumentDiff(oldDoc, &documents.Document{
			Title:   req.Title,
			Content: req.Content,
			State:   req.State,
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		added, deleted, err := s.fHandler.HandleFileChangesForParent(ctx, tx, oldDoc.Id, req.Files)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if added > 0 || deleted > 0 {
			diff.FilesChange = &documents.DocFilesChange{
				Added:   added,
				Deleted: deleted,
			}
		}

		if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
			DocumentId:   oldDoc.Id,
			ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.Job,
			Data: &documents.DocActivityData{
				Data: &documents.DocActivityData_Updated{
					Updated: diff,
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if tmpl != nil {
			if err := s.createOrUpdateWorkflowState(ctx, tx, oldDoc.Id, tmpl.Workflow); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}

	}

	if err := s.handleDocumentAccessChange(ctx, tx, oldDoc.Id, userInfo, req.Access, true); err != nil {
		return nil, err
	}

	if !onlyUpdateAccess {
		if oldDoc.Draft != req.Draft {
			if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
				DocumentId:   oldDoc.Id,
				ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DRAFT_TOGGLED,
				CreatorId:    &userInfo.UserId,
				CreatorJob:   userInfo.Job,
			}); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	doc, err := s.getDocument(ctx,
		tDocument.ID.EQ(jet.Uint64(req.DocumentId)),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.collabServer.SendTargetSaved(ctx, doc.Id)

	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbdocuments.UpdateDocumentResponse{
		Document: doc,
	}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *pbdocuments.DeleteDocumentRequest) (*pbdocuments.DeleteDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "DeleteDocument",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		if !userInfo.Superuser {
			return nil, errorsdocuments.ErrDocDeleteDenied
		}
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Require a reason if the document is not already deleted
	if doc.DeletedAt == nil && req.Reason == nil {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceDeleteDocumentPerm, permsdocuments.DocumentsServiceDeleteDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	deletedAtTime := jet.CURRENT_TIMESTAMP()
	if doc.DeletedAt != nil && userInfo.Superuser {
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
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
		Reason:       req.Reason,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbdocuments.DeleteDocumentResponse{}, nil
}

func (s *Server) ToggleDocument(ctx context.Context, req *pbdocuments.ToggleDocumentRequest) (*pbdocuments.ToggleDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "ToggleDocument",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_STATUS)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		if !userInfo.Superuser {
			return nil, errorsdocuments.ErrDocToggleDenied
		}
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	var tmpl *documents.Template
	if !req.Closed && doc.TemplateId != nil { // If the document is opened, get template so we can update the reminder/auto close times
		tmpl, err = s.getTemplate(ctx, *doc.TemplateId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceToggleDocumentPerm, permsdocuments.DocumentsServiceToggleDocumentAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocuments.ErrDocToggleDenied
	}

	activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED
	if !req.Closed {
		activityType = documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_OPEN
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   doc.Id,
		ActivityType: activityType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, doc.Id, tmpl.Workflow); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.ToggleDocumentResponse{}, nil
}

func (s *Server) ChangeDocumentOwner(ctx context.Context, req *pbdocuments.ChangeDocumentOwnerRequest) (*pbdocuments.ChangeDocumentOwnerResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "ChangeDocumentOwner",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		return nil, errorsdocuments.ErrDocOwnerFailed
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(req.DocumentId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Document must be created by the same job
	if doc.CreatorJob != userInfo.Job {
		return nil, errorsdocuments.ErrDocOwnerWrongJob
	}

	// If user is not a super user make sure they can only change owner to themselves
	if req.NewUserId == nil || !userInfo.Superuser {
		req.NewUserId = &userInfo.UserId
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsdocuments.DocumentsServicePerm, permsdocuments.DocumentsServiceChangeDocumentOwnerPerm, permsdocuments.DocumentsServiceChangeDocumentOwnerAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.CreatorJob, doc.Creator) {
		return nil, errorsdocuments.ErrDocOwnerFailed
	}

	tUsers := tables.User().AS("user_short")

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
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if newOwner.UserId <= 0 {
		return nil, errorsdocuments.ErrFailedQuery
	}

	// Allow super users to transfer documents cross jobs
	if !userInfo.Superuser {
		if newOwner.Job != doc.CreatorJob {
			return nil, errorsdocuments.ErrDocOwnerWrongJob
		}

		if doc.CreatorId != nil && *doc.CreatorId == userInfo.UserId {
			return nil, errorsdocuments.ErrDocSameOwner
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.updateDocumentOwner(ctx, tx, req.DocumentId, userInfo, &newOwner); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbdocuments.ChangeDocumentOwnerResponse{}, nil
}
