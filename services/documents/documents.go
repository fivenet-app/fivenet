package documents

import (
	context "context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
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
	tDMeta         = table.FivenetDocumentsMeta.AS("meta")
	tDAccess       = table.FivenetDocumentsAccess.AS("job_access")
)

func (s *Server) ListDocuments(
	ctx context.Context,
	req *pbdocuments.ListDocumentsRequest,
) (*pbdocuments.ListDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := false

	condition := mysql.Bool(true)
	if req.Search != nil && req.GetSearch() != "" {
		logRequest = true
		condition = dbutils.MATCH(tDocumentShort.Title, mysql.String(req.GetSearch()))
	}
	if len(req.GetCategoryIds()) > 0 {
		ids := make([]mysql.Expression, len(req.GetCategoryIds()))
		for i := range req.GetCategoryIds() {
			ids[i] = mysql.Int64(req.GetCategoryIds()[i])
		}

		condition = condition.AND(
			tDocumentShort.CategoryID.IN(ids...),
		)
	}

	if len(req.GetCreatorIds()) > 0 {
		logRequest = true
		ids := make([]mysql.Expression, len(req.GetCreatorIds()))
		for i := range req.GetCreatorIds() {
			ids[i] = mysql.Int32(req.GetCreatorIds()[i])
		}

		condition = condition.AND(
			tDocumentShort.CreatorID.IN(ids...),
		)
	}

	if req.GetFrom() != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.GT_EQ(
			mysql.TimestampT(req.GetFrom().AsTime()),
		))
	}
	if req.GetTo() != nil {
		condition = condition.AND(tDocumentShort.CreatedAt.LT_EQ(
			mysql.TimestampT(req.GetTo().AsTime()),
		))
	}

	if req.Closed != nil {
		condition = condition.AND(tDocumentShort.Closed.EQ(
			mysql.Bool(req.GetClosed()),
		))
	}

	if len(req.GetDocumentIds()) > 0 {
		ids := make([]mysql.Expression, len(req.GetDocumentIds()))
		for i := range req.GetDocumentIds() {
			ids[i] = mysql.Int64(req.GetDocumentIds()[i])
		}

		condition = condition.AND(
			tDocumentShort.ID.IN(ids...),
		)
	}

	if req.OnlyDrafts != nil {
		condition = condition.AND(tDocumentShort.Draft.EQ(mysql.Bool(req.GetOnlyDrafts())))
	}

	if !logRequest {
		grpc_audit.Skip(ctx)
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "title":
			column = tDocumentShort.Title

		case "createdAt":
			fallthrough
		default:
			column = tDocumentShort.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
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

	pag, limit := req.GetPagination().
		GetResponseWithPageSize(database.NoTotalCount, DocsDefaultPageSize)
	resp := &pbdocuments.ListDocumentsResponse{
		Pagination: pag,
	}

	stmt := s.listDocumentsQuery(
		condition,
		nil,
		nil,
		userInfo,
		func(stmt mysql.SelectStatement) mysql.SelectStatement {
			return stmt.
				ORDER_BY(orderBys...).
				OFFSET(req.GetPagination().GetOffset()).
				LIMIT(limit)
		},
	)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetDocuments() {
		if resp.GetDocuments()[i].GetCreator() != nil {
			jobInfoFn(resp.GetDocuments()[i].GetCreator())
		}

		if job := s.enricher.GetJobByName(resp.GetDocuments()[i].GetCreatorJob()); job != nil {
			resp.Documents[i].CreatorJobLabel = &job.Label
		}
	}

	return resp, nil
}

func (s *Server) GetDocument(
	ctx context.Context,
	req *pbdocuments.GetDocumentRequest,
) (*pbdocuments.GetDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	infoOnly := req.InfoOnly != nil && req.GetInfoOnly()
	withContent := req.InfoOnly == nil || !req.GetInfoOnly()

	resp := &pbdocuments.GetDocumentResponse{}
	resp.Document, err = s.getDocument(ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())), userInfo, withContent)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.GetDocument() == nil || resp.GetDocument().GetId() <= 0 {
		return nil, errorsdocuments.ErrNotFoundOrNoPerms
	}

	if resp.GetDocument().GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.GetDocument().GetCreator())
	}

	resp.Document.Pin, err = s.getDocumentPin(ctx, resp.GetDocument().GetId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if !infoOnly {
		docAccess, err := s.GetDocumentAccess(ctx, &pbdocuments.GetDocumentAccessRequest{
			DocumentId: resp.GetDocument().GetId(),
		})
		if err != nil {
			if st, ok := status.FromError(err); !ok {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			} else if st.Code() != codes.PermissionDenied {
				// Ignore permission denied as we are simply getting the document
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
		if docAccess != nil {
			resp.Access = docAccess.GetAccess()
		}

		files, err := s.fHandler.ListFilesForParentID(ctx, resp.GetDocument().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		resp.Document.Files = files
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return resp, nil
}

func (s *Server) getDocument(
	ctx context.Context,
	condition mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	withContent bool,
) (*documents.Document, error) {
	var doc documents.Document

	stmt := s.getDocumentQuery(condition, nil, userInfo, withContent).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &doc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if doc.GetCreator() != nil {
		s.enricher.EnrichJobInfo(doc.GetCreator())
	}
	if doc.Meta == nil {
		doc.Meta = &documents.DocumentMeta{
			DocumentId: doc.GetId(),
		}
	}

	return &doc, nil
}

func (s *Server) CreateDocument(
	ctx context.Context,
	req *pbdocuments.CreateDocumentRequest,
) (*pbdocuments.CreateDocumentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	var docContent string
	var docTitle string
	var docState string
	var categoryId *int64
	docAccess := &documents.DocumentAccess{}
	docReferences := []*documents.DocumentReference{}
	docRelations := []*documents.DocumentRelation{}

	var tmpl *documents.Template
	if req.GetTemplateId() > 0 {
		var err error
		tmpl, err = s.getTemplate(ctx, req.GetTemplateId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		docTitle, docState, docContent, err = s.renderTemplate(tmpl, req.GetTemplateData())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		// Set access based on template
		docAccess = &documents.DocumentAccess{
			Jobs:  tmpl.GetContentAccess().GetJobs(),
			Users: tmpl.GetContentAccess().GetUsers(),
		}

		if tmpl.GetCategory() != nil {
			categoryId = &tmpl.Category.Id
		}

		// Add references from template data documents if not already present
		if tmpl != nil && req.GetTemplateData() != nil {
			if len(req.GetTemplateData().GetDocuments()) > 0 {
				for _, doc := range req.GetTemplateData().GetDocuments() {
					exists := false
					for _, reference := range docReferences {
						if reference.GetTargetDocumentId() == doc.GetId() {
							exists = true
							break
						}
					}

					if !exists {
						docReferences = append(docReferences, &documents.DocumentReference{
							// Id will be assigned by backend or can be zero for new
							SourceDocumentId: 0, // will be set after insert
							TargetDocumentId: doc.GetId(),
							// TargetDocument can be set if needed, or left nil
							CreatorId: &userInfo.UserId,
							// Creator can be set if needed, or left nil
							Reference: documents.DocReference_DOC_REFERENCE_SOLVES,
						})
					}
				}
			}

			// Add relations from template data users if not already present
			if len(req.GetTemplateData().GetUsers()) > 0 {
				for _, user := range req.GetTemplateData().GetUsers() {
					exists := false
					for _, relation := range docRelations {
						if relation.GetTargetUserId() == user.GetUserId() {
							exists = true
							break
						}
					}

					if !exists {
						docRelations = append(docRelations, &documents.DocumentRelation{
							// Id will be assigned by backend or can be zero for new
							DocumentId:   0, // will be set after insert
							TargetUserId: user.GetUserId(),
							// TargetUser can be set if needed, or left nil
							SourceUserId: userInfo.GetUserId(),
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
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
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
			req.GetContentType(),
			docState,
			mysql.NULL,
			false,
			true,
			false,
			req.GetTemplateId(),
			userInfo.GetUserId(),
			userInfo.GetJob(),
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
		DocumentId:   lastId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.handleDocumentAccessChange(ctx, tx, lastId, userInfo, docAccess, false); err != nil {
		return nil, err
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, lastId, tmpl.GetWorkflow()); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	for _, ref := range docReferences {
		ref.SourceDocumentId = lastId
		if _, err := s.addDocumentReference(ctx, tx, ref); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}
	for _, rel := range docRelations {
		rel.DocumentId = lastId
		if _, err := s.addDocumentRelation(ctx, tx, userInfo, rel); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return &pbdocuments.CreateDocumentResponse{
		Id: lastId,
	}, nil
}

func (s *Server) UpdateDocument(
	ctx context.Context,
	req *pbdocuments.UpdateDocumentRequest,
) (*pbdocuments.UpdateDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	var onlyUpdateAccess bool
	if !check && !userInfo.GetSuperuser() {
		onlyUpdateAccess, err = s.access.CanUserAccessTarget(
			ctx,
			req.GetDocumentId(),
			userInfo,
			documents.AccessLevel_ACCESS_LEVEL_ACCESS,
		)
		if err != nil {
			return nil, errorsdocuments.ErrPermissionDenied
		}
		if !onlyUpdateAccess {
			return nil, errorsdocuments.ErrPermissionDenied
		}
	}

	oldDoc, err := s.getDocument(ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Document is closed and the update request isn't re-opening the document
	if oldDoc.GetMeta().GetClosed() && req.GetMeta().GetClosed() && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrClosedDoc
	}

	// A document can only be switched to published once
	if !oldDoc.GetMeta().GetDraft() && oldDoc.GetMeta().GetDraft() != req.GetMeta().GetDraft() {
		// Allow a super user to change the draft state
		if !userInfo.GetSuperuser() {
			req.GetMeta().Draft = oldDoc.GetMeta().GetDraft()
		}
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceUpdateDocumentPerm,
		permsdocuments.DocumentsServiceUpdateDocumentAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields,
		userInfo,
		oldDoc.GetCreatorJob(),
		oldDoc.GetCreator(),
	) {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	var tmpl *documents.Template
	if oldDoc.GetTemplateId() > 0 {
		var err error
		tmpl, err = s.getTemplate(ctx, oldDoc.GetTemplateId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if !s.checkAccessAgainstTemplate(tmpl, req.GetAccess()) {
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
				req.GetTitle(),
				req.GetContent().GetSummary(DocSummaryLength),
				req.GetContent(),
				mysql.NULL,
				req.GetMeta().GetState(),
				req.GetMeta().GetClosed(),
				req.GetMeta().GetDraft(),
				req.GetMeta().GetPublic(),
			).
			WHERE(
				tDocument.ID.EQ(mysql.Int64(oldDoc.GetId())),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		diff, err := s.generateDocumentDiff(oldDoc, &documents.Document{
			Title:   req.GetTitle(),
			Content: req.GetContent(),
			Meta: &documents.DocumentMeta{
				State: req.GetMeta().GetState(),
			},
		})
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		added, deleted, err := s.fHandler.HandleFileChangesForParent(
			ctx,
			tx,
			oldDoc.GetId(),
			req.GetFiles(),
		)
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
			DocumentId:   oldDoc.GetId(),
			ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED,
			CreatorId:    &userInfo.UserId,
			CreatorJob:   userInfo.GetJob(),
			Data: &documents.DocActivityData{
				Data: &documents.DocActivityData_Updated{
					Updated: diff,
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if tmpl != nil && tmpl.GetWorkflow() != nil {
			if err := s.createOrUpdateWorkflowState(ctx, tx, oldDoc.GetId(), tmpl.GetWorkflow()); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	}

	if err := s.handleDocumentAccessChange(ctx, tx, oldDoc.GetId(), userInfo, req.GetAccess(), true); err != nil {
		return nil, err
	}

	if !onlyUpdateAccess {
		if oldDoc.GetMeta().GetDraft() != req.GetMeta().GetDraft() {
			if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
				DocumentId:   oldDoc.GetId(),
				ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DRAFT_TOGGLED,
				CreatorId:    &userInfo.UserId,
				CreatorJob:   userInfo.GetJob(),
			}); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			if tmpl != nil && !req.GetMeta().GetDraft() {
				if err := s.handleDocumentPublish(ctx, tx, userInfo, oldDoc, tmpl); err != nil {
					return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
				}
			}
		}

		if err := s.handleApprovalOnEditBehaviors(ctx, tx, oldDoc); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	doc, err := s.getDocument(ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.collabServer.SendTargetSaved(ctx, doc.GetId())

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

func (s *Server) handleDocumentPublish(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	doc *documents.Document,
	tmpl *documents.Template,
) error {
	apr := tmpl.GetApproval()
	if apr == nil || !apr.GetEnabled() {
		return nil
	}

	pol, err := s.getApprovalPolicy(
		ctx,
		tx,
		table.FivenetDocumentsApprovalPolicies.AS(
			"approval_policy",
		).DocumentID.EQ(
			mysql.Int64(doc.GetId()),
		),
	)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if pol != nil {
		// A policy already exists, don't update the existing one
		return nil
	}

	now := time.Now()

	// Fill in and create new approval policy
	newPol := &documents.ApprovalPolicy{
		SnapshotDate: timestamp.New(now),
	}
	if apr.GetPolicy() != nil {
		newPol.RuleKind = apr.GetPolicy().GetRuleKind()
		newPol.OnEditBehavior = apr.GetPolicy().GetOnEditBehavior()
		requiredCount := apr.GetPolicy().GetRequiredCount()

		switch newPol.RuleKind {
		case documents.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL:
			requiredCount = 0
		case documents.ApprovalRuleKind_APPROVAL_RULE_KIND_QUORUM_ANY:
			if requiredCount <= 0 {
				requiredCount = int32(len(apr.GetTasks()))
				if requiredCount == 0 {
					requiredCount = 1
				}
			}
		}
		if requiredCount > 0 {
			newPol.RequiredCount = &requiredCount
		}

		newPol.SignatureRequired = apr.GetPolicy().GetSignatureRequired()
	}
	newPol.Default()

	if err = s.createApprovalPolicy(ctx, tx, doc.GetId(), newPol); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	seeds := []*pbdocuments.ApprovalTaskSeed{}
	for _, task := range apr.GetTasks() {
		var dueAt *timestamp.Timestamp
		if task.GetDueInDays() > 0 {
			dueTime := now.AddDate(0, 0, int(task.GetDueInDays()))
			dueAt = timestamp.New(dueTime)
		}

		slots := task.GetSlots()
		if slots <= 0 {
			slots = 1
		}

		seeds = append(seeds, &pbdocuments.ApprovalTaskSeed{
			UserId:            task.GetUserId(),
			Job:               task.GetJob(),
			MinimumGrade:      task.GetMinimumGrade(),
			DueAt:             dueAt,
			SignatureRequired: task.GetSignatureRequired(),
			Slots:             slots,
			Label:             task.Label,
			Comment:           task.Comment,
		})
	}

	if _, _, err := s.createApprovalTasks(
		ctx,
		tx,
		userInfo,
		doc.GetId(),
		now,
		seeds,
	); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return nil
}

func (s *Server) DeleteDocument(
	ctx context.Context,
	req *pbdocuments.DeleteDocumentRequest,
) (*pbdocuments.DeleteDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		if !userInfo.GetSuperuser() {
			return nil, errorsdocuments.ErrDocDeleteDenied
		}
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

	// Require a reason if the document is not already deleted
	if doc.GetDeletedAt() == nil && req.Reason == nil {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceDeleteDocumentPerm,
		permsdocuments.DocumentsServiceDeleteDocumentAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.GetCreatorJob(), doc.GetCreator()) {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	deletedAtTime := mysql.CURRENT_TIMESTAMP()
	if doc.GetDeletedAt() != nil && userInfo.GetSuperuser() {
		deletedAtTime = mysql.TimestampExp(mysql.NULL)
	}

	stmt := tDocument.
		UPDATE(
			tDocument.DeletedAt,
		).
		SET(
			tDocument.DeletedAt.SET(deletedAtTime),
		).
		WHERE(
			tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, s.db, &documents.DocActivity{
		DocumentId:   req.GetDocumentId(),
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Reason:       req.Reason,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	return &pbdocuments.DeleteDocumentResponse{}, nil
}

func (s *Server) ToggleDocument(
	ctx context.Context,
	req *pbdocuments.ToggleDocumentRequest,
) (*pbdocuments.ToggleDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_STATUS,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		if !userInfo.GetSuperuser() {
			return nil, errorsdocuments.ErrDocToggleDenied
		}
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

	var tmpl *documents.Template
	if !req.GetClosed() && doc.GetTemplateId() > 0 {
		// If the document is opened, get template so we can update the reminder/auto close times
		tmpl, err = s.getTemplate(ctx, doc.GetTemplateId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceToggleDocumentPerm,
		permsdocuments.DocumentsServiceToggleDocumentAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.GetCreatorJob(), doc.GetCreator()) {
		return nil, errorsdocuments.ErrDocToggleDenied
	}

	activityType := documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED
	if !req.GetClosed() {
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
			req.GetClosed(),
		).
		WHERE(
			tDocument.ID.EQ(mysql.Int64(doc.GetId())),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   doc.GetId(),
		ActivityType: activityType,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if tmpl != nil {
		if err := s.createOrUpdateWorkflowState(ctx, tx, doc.GetId(), tmpl.GetWorkflow()); err != nil {
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.ToggleDocumentResponse{}, nil
}

func (s *Server) ChangeDocumentOwner(
	ctx context.Context,
	req *pbdocuments.ChangeDocumentOwnerRequest,
) (*pbdocuments.ChangeDocumentOwnerResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		return nil, errorsdocuments.ErrDocOwnerFailed
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

	// Document must be created by the same job
	if doc.GetCreatorJob() != userInfo.GetJob() {
		return nil, errorsdocuments.ErrDocOwnerWrongJob
	}

	// If user is not a super user make sure they can only change owner to themselves
	if req.NewUserId == nil || !userInfo.GetSuperuser() {
		req.NewUserId = &userInfo.UserId
	}

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceChangeDocumentOwnerPerm,
		permsdocuments.DocumentsServiceChangeDocumentOwnerAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.GetCreatorJob(), doc.GetCreator()) {
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
		WHERE(tUsers.ID.EQ(mysql.Int32(req.GetNewUserId()))).
		LIMIT(1)

	var newOwner users.UserShort
	if err := stmt.QueryContext(ctx, s.db, &newOwner); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if newOwner.GetUserId() <= 0 {
		return nil, errorsdocuments.ErrFailedQuery
	}

	// Allow super users to transfer documents cross jobs
	if !userInfo.GetSuperuser() {
		if newOwner.GetJob() != doc.GetCreatorJob() {
			return nil, errorsdocuments.ErrDocOwnerWrongJob
		}

		if doc.CreatorId != nil && doc.GetCreatorId() == userInfo.GetUserId() {
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

	if err := s.updateDocumentOwner(ctx, tx, req.GetDocumentId(), userInfo, &newOwner); err != nil {
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.ChangeDocumentOwnerResponse{}, nil
}
