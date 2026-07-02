package documents

import (
	context "context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsapproval "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/approval"
	documentsreferences "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/references"
	documentsrelations "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/relations"
	documentstemplates "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/templates"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	documentsstore "github.com/fivenet-app/fivenet/v2026/stores/documents"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

const (
	DocsDefaultPageSize = 20
	DocSummaryLength    = 128
)

var (
	tDocument      = table.FivenetDocuments.AS("document")
	tDocumentShort = table.FivenetDocuments.AS("document_short")
	tDCategory     = table.FivenetDocumentsCategories.AS("category")
	tDMeta         = table.FivenetDocumentsMeta.AS("meta")
	tDPins         = table.FivenetDocumentsPins.AS("pin")
)

func documentFilesEqual(left, right []*file.File) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if !proto.Equal(left[i], right[i]) {
			return false
		}
	}

	return true
}

func documentUpdateContentChanged(
	oldDoc *documents.Document,
	req *pbdocuments.UpdateDocumentRequest,
) bool {
	if oldDoc.GetTitle() != req.GetTitle() {
		return true
	}

	if oldDoc.GetCategoryId() != req.GetCategoryId() {
		return true
	}

	if !proto.Equal(oldDoc.GetContent(), req.GetContent()) {
		return true
	}

	if !proto.Equal(oldDoc.GetData(), req.GetData()) {
		return true
	}

	return !documentFilesEqual(oldDoc.GetFiles(), req.GetFiles())
}

func documentUpdateStatusChanged(
	oldDoc *documents.Document,
	req *pbdocuments.UpdateDocumentRequest,
) bool {
	oldMeta := oldDoc.GetMeta()
	newMeta := req.GetMeta()

	return oldMeta.GetState() != newMeta.GetState() ||
		oldMeta.GetClosed() != newMeta.GetClosed() ||
		oldMeta.GetDraft() != newMeta.GetDraft() ||
		oldMeta.GetPublic() != newMeta.GetPublic()
}

func documentUpdateAccessChanged(
	oldAccess *documentsaccess.DocumentAccess,
	req *pbdocuments.UpdateDocumentRequest,
) bool {
	if req.GetAccess() == nil {
		return false
	}

	return !proto.Equal(oldAccess, req.GetAccess())
}

func (s *Server) ListDocuments(
	ctx context.Context,
	req *pbdocuments.ListDocumentsRequest,
) (*pbdocuments.ListDocumentsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	logRequest := req.Search != nil && req.GetSearch() != ""

	if len(req.GetCreatorIds()) > 0 {
		logRequest = true
	}

	if !logRequest {
		grpc_audit.Skip(ctx)
	}

	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().
		GetResponseWithPageSize(database.NoTotalCount, DocsDefaultPageSize)
	resp := &pbdocuments.ListDocumentsResponse{
		Pagination: pag,
	}

	docs, err := s.store.List(ctx, documentsstore.ListQuery{
		Search:      req.GetSearch(),
		CategoryIDs: req.GetCategoryIds(),
		CreatorIDs:  req.GetCreatorIds(),
		From:        req.GetFrom(),
		To:          req.GetTo(),
		Closed:      req.Closed,
		DocumentIDs: req.GetDocumentIds(),
		OnlyDrafts:  req.OnlyDrafts,
		Sort:        req.GetSort(),
		Offset:      req.GetPagination().GetOffset(),
		Limit:       limit,
		IncludePhoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		UserInfo: userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	resp.Documents = docs

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
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)
	fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	infoOnly := req.InfoOnly != nil && req.GetInfoOnly()
	withContent := req.InfoOnly == nil || !req.GetInfoOnly()

	resp := &pbdocuments.GetDocumentResponse{}
	resp.Document, err = s.store.Get(ctx, documentsstore.GetQuery{
		DocumentID: req.GetDocumentId(),
		IncludePhoneNumber: fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		),
		WithContent: withContent,
		UserInfo:    userInfo,
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if resp.GetDocument() == nil || resp.GetDocument().GetId() <= 0 {
		return nil, errorsdocuments.ErrNotFoundOrNoPerms
	}

	if resp.GetDocument().GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, resp.GetDocument().GetCreator())
	}

	resp.Document.Pin, err = s.store.GetDocumentPin(ctx, resp.GetDocument().GetId(), userInfo)
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
	if doc.GetId() == 0 {
		return nil, nil
	}

	if doc.GetCreator() != nil {
		s.enricher.EnrichJobInfo(doc.GetCreator())
	}
	if doc.GetMeta() == nil {
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

	docContent := &content.Content{
		Version:     content.ContentVersionTiptapV1,
		ContentType: content.ContentType_CONTENT_TYPE_TIPTAP_JSON,
		TiptapJson:  nil,
	}
	var docTitle string
	var docState string
	var categoryId *int64
	docAccess := &documentsaccess.DocumentAccess{}
	docReferences := []*documentsreferences.DocumentReference{}
	docRelations := []*documentsrelations.DocumentRelation{}

	var tmpl *documentstemplates.Template
	if req.GetTemplateId() > 0 {
		var err error
		tmpl, err = s.store.GetTemplate(ctx, req.GetTemplateId(), false)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		var tplContent string
		docTitle, docState, tplContent, err = s.renderTemplate(tmpl, req.GetTemplateData())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		// Build Content object
		htmlNode, err := content.FromHTML(tplContent)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		docContent = &content.Content{
			Version:     content.ContentVersionLegacyJSONV1,
			ContentType: content.ContentType_CONTENT_TYPE_HTML,
			Content:     htmlNode,
		}

		// Set access based on template
		docAccess, err = access.SanitizeAccessJobs(s.jobs, tmpl.GetContentAccess())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		docAccess = access.NormalizeRequiredAccessFloors(docAccess)

		if tmpl.GetCategory() != nil {
			categoryId = &tmpl.Category.Id
		}

		// Add references from template data documents if not already present
		if tmpl != nil && req.GetTemplateData() != nil {
			if len(req.GetTemplateData().GetDocuments()) > 0 {
				for _, doc := range req.GetTemplateData().GetDocuments() {
					if doc == nil {
						continue
					}

					exists := false
					for _, reference := range docReferences {
						if reference.GetTargetDocumentId() == doc.GetId() {
							exists = true
							break
						}
					}

					if !exists {
						docReferences = append(
							docReferences,
							&documentsreferences.DocumentReference{
								// Id will be assigned by backend or can be zero for new
								SourceDocumentId: 0, // will be set after insert
								TargetDocumentId: doc.GetId(),
								// TargetDocument can be set if needed, or left nil
								CreatorId: &userInfo.UserId,
								// Creator can be set if needed, or left nil
								Reference: documentsreferences.DocReference_DOC_REFERENCE_SOLVES,
							},
						)
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
						docRelations = append(docRelations, &documentsrelations.DocumentRelation{
							// Id will be assigned by backend or can be zero for new
							DocumentId:   0, // will be set after insert
							TargetUserId: user.GetUserId(),
							// TargetUser can be set if needed, or left nil
							SourceUserId: userInfo.GetUserId(),
							// SourceUser can be set if needed, or left nil
							Relation: documentsrelations.DocRelation_DOC_RELATION_CAUSED,
						})
					}
				}
			}
		}
	} else {
		// Add minimum access for the creator's job
		docAccess.Jobs = append(docAccess.Jobs, &documentsaccess.DocumentJobAccess{
			Job:          userInfo.GetJob(),
			MinimumGrade: userInfo.GetJobGrade(),
			Access:       int32(documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
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
			tDocument.CategoryID,
			tDocument.Title,
			tDocument.Summary,
			tDocument.ContentJSON,
			tDocument.ContentText,
			tDocument.ContentType,
			tDocument.FirstHeading,
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
			categoryId,
			docTitle,
			"",
			docContent,
			"",
			int32(docContent.GetContentType()),
			"",
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

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   lastId,
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_CREATED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.handleDocumentAccessChange(
		ctx,
		tx,
		lastId,
		userInfo,
		docAccess,
		false,
	); err != nil {
		return nil, err
	}

	if tmpl != nil {
		if err := s.store.UpsertWorkflowState(ctx, tx, lastId, tmpl.GetWorkflow()); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	for _, ref := range docReferences {
		ref.SourceDocumentId = lastId
		if _, err := s.store.CreateDocumentReference(ctx, tx, ref); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}
	for _, rel := range docRelations {
		rel.DocumentId = lastId
		_, created, err := s.store.CreateDocumentRelation(ctx, tx, rel)
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
		if created {
			if err := s.addUserActivity(
				ctx,
				tx,
				userInfo.GetUserId(),
				rel.GetTargetUserId(),
				usersactivity.UserActivityType_USER_ACTIVITY_TYPE_DOCUMENT,
				"",
				&usersactivity.UserActivityData{
					Data: &usersactivity.UserActivityData_DocumentRelation{
						DocumentRelation: &usersactivity.CitizenDocumentRelation{
							Added:      true,
							DocumentId: rel.GetDocumentId(),
							Relation:   rel.GetRelation(),
						},
					},
				},
			); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}

			if rel.GetRelation() == documentsrelations.DocRelation_DOC_RELATION_MENTIONED {
				if err := s.notifyMentionedUser(
					ctx,
					lastId,
					userInfo.GetUserId(),
					rel.GetTargetUserId(),
				); err != nil {
					return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
				}
			}
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
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	canAccessUpdate, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !canAccessUpdate {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	canStatusUpdate, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	canEditUpdate, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}

	oldDoc, err := s.getDocument(ctx,
		tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		userInfo, true)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if oldDoc == nil {
		return nil, errorsdocuments.ErrPermissionDenied
	}

	// Document is closed and the update request isn't re-opening the document
	if oldDoc.GetMeta().GetClosed() && req.GetMeta().GetClosed() && !userInfo.GetJobAdmin() {
		return nil, errorsdocuments.ErrClosedDoc
	}

	// A document can only be switched to published once
	if !oldDoc.GetMeta().GetDraft() && oldDoc.GetMeta().GetDraft() != req.GetMeta().GetDraft() {
		// Allow a super user to change the draft state
		if !userInfo.GetJobAdmin() {
			req.GetMeta().Draft = oldDoc.GetMeta().GetDraft()
		}
	}

	var oldAccess *documentsaccess.DocumentAccess
	if req.GetAccess() != nil {
		oldAccess, err = s.store.GetDocumentAccess(ctx, oldDoc.GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	accessChanged := documentUpdateAccessChanged(oldAccess, req)
	statusChanged := documentUpdateStatusChanged(oldDoc, req)
	contentChanged := documentUpdateContentChanged(oldDoc, req)

	if accessChanged && !canAccessUpdate {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}
	if statusChanged && !canStatusUpdate && !canEditUpdate {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}
	if contentChanged && !canEditUpdate {
		return nil, errorsdocuments.ErrDocUpdateDenied
	}

	var tmpl *documentstemplates.Template
	if oldDoc.GetTemplateId() > 0 {
		var err error
		tmpl, err = s.store.GetTemplate(ctx, oldDoc.GetTemplateId(), false)
		if err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
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

	if contentChanged {
		extracted := req.GetContent().Extract()

		tDocument := table.FivenetDocuments
		stmt := tDocument.
			UPDATE(
				tDocument.CategoryID,
				tDocument.Title,
				tDocument.Summary,
				tDocument.WordCount,
				tDocument.FirstHeading,
				tDocument.ContentType,
				tDocument.ContentJSON,
				tDocument.ContentText,
				tDocument.Data,
			).
			SET(
				req.CategoryId,
				req.GetTitle(),
				extracted.GetSummary(DocSummaryLength),
				extracted.GetWordCount(),
				extracted.GetFirstHeading(),
				int32(content.ContentType_CONTENT_TYPE_TIPTAP_JSON),
				req.GetContent(),
				extracted.GetText(),
				req.GetData(),
			).
			WHERE(
				tDocument.ID.EQ(mysql.Int64(oldDoc.GetId())),
			).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}

		if tmpl != nil && tmpl.GetWorkflow() != nil {
			if err := s.store.UpsertWorkflowState(
				ctx,
				tx,
				oldDoc.GetId(),
				tmpl.GetWorkflow(),
			); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}
	}

	if statusChanged {
		tDocument := table.FivenetDocuments
		stmt := tDocument.
			UPDATE(
				tDocument.State,
				tDocument.Closed,
				tDocument.Draft,
				tDocument.Public,
			).
			SET(
				req.GetMeta().GetState(),
				req.GetMeta().GetClosed(),
				req.GetMeta().GetDraft(),
				req.GetMeta().GetPublic(),
			).
			WHERE(
				tDocument.ID.EQ(mysql.Int64(oldDoc.GetId())),
			).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	if accessChanged {
		if err := s.handleDocumentAccessChange(
			ctx,
			tx,
			oldDoc.GetId(),
			userInfo,
			req.GetAccess(),
			true,
		); err != nil {
			return nil, err
		}
	}

	if contentChanged || statusChanged {
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

		if contentChanged {
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
				diff.FilesChange = &documentsactivity.DocFilesChange{
					Added:   added,
					Deleted: deleted,
				}
			}
		}

		if diff.HasChanges() {
			if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
				DocumentId:   oldDoc.GetId(),
				ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_UPDATED,
				CreatorId:    &userInfo.UserId,
				CreatorJob:   userInfo.GetJob(),
				Data: &documentsactivity.DocActivityData{
					Data: &documentsactivity.DocActivityData_Updated{
						Updated: diff,
					},
				},
			}); err != nil {
				return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
			}
		}

		if oldDoc.GetMeta().GetDraft() != req.GetMeta().GetDraft() {
			if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
				DocumentId:   oldDoc.GetId(),
				ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_DRAFT_TOGGLED,
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

	if err := s.stats.RebuildDocumentMetrics(ctx, doc); err != nil {
		s.logger.Warn(
			"failed to rebuild document metrics after update",
			zap.Int64("document_id", req.GetDocumentId()),
			zap.Error(err),
		)
	}

	s.collabServer.SendTargetSaved(ctx, doc.GetId())

	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

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
	tmpl *documentstemplates.Template,
) error {
	apr := tmpl.GetApproval()
	if apr == nil || !apr.GetEnabled() {
		return nil
	}

	pol, err := s.store.GetApprovalPolicy(
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

	now := timestamp.New(time.Now().Truncate(time.Second))

	// Fill in and create new approval policy
	newPol := &documentsapproval.ApprovalPolicy{
		SnapshotDate: now,
	}
	if apr.GetPolicy() != nil {
		newPol.RuleKind = apr.GetPolicy().GetRuleKind()
		newPol.OnEditBehavior = apr.GetPolicy().GetOnEditBehavior()
		requiredCount := apr.GetPolicy().GetRequiredCount()

		switch newPol.GetRuleKind() {
		case documentsapproval.ApprovalRuleKind_APPROVAL_RULE_KIND_REQUIRE_ALL:
			requiredCount = 0
		case documentsapproval.ApprovalRuleKind_APPROVAL_RULE_KIND_QUORUM_ANY:
			if requiredCount <= 0 {
				//nolint:gosec // G115: there can't be more than math.MaxInt32 tasks due to API validation, so this cast is safe
				requiredCount = int32(len(apr.GetTasks()))
				// If required count isn't set, default to 1 approver required
				if requiredCount == 0 {
					requiredCount = 1
				}
			}
		}
		if requiredCount > 0 {
			newPol.RequiredCount = &requiredCount
		}

		newPol.SignatureRequired = apr.GetPolicy().GetSignatureRequired()
		newPol.SelfApproveAllowed = apr.GetPolicy().GetSelfApproveAllowed()
	}
	newPol.Default()

	if err = s.store.CreateApprovalPolicy(ctx, tx, doc.GetId(), newPol); err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	seeds := []*pbdocuments.ApprovalTaskSeed{}
	for _, task := range apr.GetTasks() {
		var dueAt *timestamp.Timestamp
		if task.GetDueInDays() > 0 {
			dueTime := now.AsTime().AddDate(0, 0, int(task.GetDueInDays()))
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

	if _, _, err := s.store.CreateApprovalTasks(
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
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetJobAdmin() {
		if !userInfo.GetJobAdmin() {
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
	if doc == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	// Require a reason if the document is not already deleted
	if doc.GetDeletedAt() == nil && req.Reason == nil {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	// Field Permission Check
	fields, err := permsdocuments.DocumentsService.DeleteDocument.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		doc.GetCreatorJob(),
		doc.GetCreator(),
	) {
		return nil, errorsdocuments.ErrDocDeleteDenied
	}

	var deletedAtTime *timestamp.Timestamp
	if doc.GetDeletedAt() == nil || !userInfo.GetJobAdmin() {
		deletedAtTime = timestamp.Now()
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)
	} else {
		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_RESTORED)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	tDocument := table.FivenetDocuments
	stmt := tDocument.
		UPDATE(
			tDocument.DeletedAt,
		).
		SET(
			tDocument.DeletedAt.SET(dbutils.TimestampToMySQL(deletedAtTime)),
		).
		WHERE(
			tDocument.ID.EQ(mysql.Int64(req.GetDocumentId())),
		).
		LIMIT(1)
	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := s.subjectAccess.RefreshTargetVisibility(ctx, tx, req.GetDocumentId()); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   req.GetDocumentId(),
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_DELETED,
		CreatorId:    &userInfo.UserId,
		CreatorJob:   userInfo.GetJob(),
		Reason:       req.Reason,
	}); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	return &pbdocuments.DeleteDocumentResponse{}, nil
}

func (s *Server) ToggleDocument(
	ctx context.Context,
	req *pbdocuments.ToggleDocumentRequest,
) (*pbdocuments.ToggleDocumentResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetJobAdmin() {
		if !userInfo.GetJobAdmin() {
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

	// Field Permission Check
	fields, err := permsdocuments.DocumentsService.ToggleDocument.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		doc.GetCreatorJob(),
		doc.GetCreator(),
	) {
		return nil, errorsdocuments.ErrDocToggleDenied
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	defer tx.Rollback()

	if err := s.store.ToggleDocument(
		ctx,
		tx,
		req.GetDocumentId(),
		doc.GetTemplateId(),
		req.GetClosed(),
		userInfo,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

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
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetJobAdmin() {
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
	if doc == nil {
		return nil, errorsdocuments.ErrFailedQuery
	}

	// Document must be created by the same job
	if doc.GetCreatorJob() != userInfo.GetJob() {
		return nil, errorsdocuments.ErrDocOwnerWrongJob
	}

	// If user is not a super user make sure they can only change owner to themselves
	if req.NewUserId == nil || !userInfo.GetJobAdmin() {
		req.NewUserId = &userInfo.UserId
	}

	// Field Permission Check
	fields, err := permsdocuments.DocumentsService.ChangeDocumentOwner.AccessTyped.Get(
		s.ps,
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(
		fields.StringList(),
		userInfo,
		doc.GetCreatorJob(),
		doc.GetCreator(),
	) {
		return nil, errorsdocuments.ErrDocOwnerFailed
	}

	tUsers := table.FivenetUser.AS("user_short")

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

	var newOwner usershort.UserShort
	if err := stmt.QueryContext(ctx, s.db, &newOwner); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	if newOwner.GetUserId() <= 0 {
		return nil, errorsdocuments.ErrFailedQuery
	}

	// Allow super users to transfer documents cross jobs
	if !userInfo.GetJobAdmin() {
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

	if err := s.store.UpdateDocumentOwner(
		ctx,
		tx,
		req.GetDocumentId(),
		userInfo,
		&newOwner,
	); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_DOCUMENT,
		Id:        &doc.Id,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	return &pbdocuments.ChangeDocumentOwnerResponse{}, nil
}
