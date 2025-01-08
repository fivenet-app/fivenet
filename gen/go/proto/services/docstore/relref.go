package docstore

import (
	context "context"
	"errors"
	"fmt"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsdocstore "github.com/fivenet-app/fivenet/gen/go/proto/services/docstore/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tDocRef = table.FivenetDocumentsReferences.AS("documentreference")
	tDocRel = table.FivenetDocumentsRelations.AS("documentrelation")
)

func (s *Server) GetDocumentReferences(ctx context.Context, req *GetDocumentReferencesRequest) (*GetDocumentReferencesResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRefsViewDenied
	}

	resp := &GetDocumentReferencesResponse{}

	var docsIds []struct {
		Source *uint64
		Target *uint64
	}
	idStmt := tDocRef.
		SELECT(
			tDocRef.SourceDocumentID.AS("source"),
			tDocRef.TargetDocumentID.AS("target"),
		).
		FROM(
			tDocRef,
		).
		WHERE(
			jet.AND(
				tDocRef.DeletedAt.IS_NULL(),
				jet.OR(
					tDocRef.SourceDocumentID.EQ(jet.Uint64(req.DocumentId)),
					tDocRef.TargetDocumentID.EQ(jet.Uint64(req.DocumentId)),
				),
			),
		)

	if err := idStmt.QueryContext(ctx, s.db, &docsIds); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	if len(docsIds) == 0 {
		return resp, nil
	}

	var docIds []uint64
	for _, v := range docsIds {
		if v.Source != nil {
			docIds = append(docIds, *v.Source)
		}
		if v.Target != nil {
			docIds = append(docIds, *v.Target)
		}
	}

	ids, err := s.access.CanUserAccessTargetIDs(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW, docIds...)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if len(ids) == 0 {
		return resp, nil
	}

	dIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		dIds[i] = jet.Uint64(ids[i])
	}

	tSourceDoc := tDocument.AS("source_document")
	tTargetDoc := tDocument.AS("target_document")
	tCreator := tables.Users().AS("creator")
	tRefCreator := tCreator.AS("ref_creator")

	stmt := tDocRef.
		SELECT(
			tDocRef.ID,
			tDocRef.CreatedAt,
			tDocRef.SourceDocumentID,
			tDocRef.Reference,
			tDocRef.TargetDocumentID,
			tDocRef.CreatorID,
			tSourceDoc.ID,
			tSourceDoc.CreatedAt,
			tSourceDoc.UpdatedAt,
			tSourceDoc.CategoryID,
			tSourceDoc.Title,
			tSourceDoc.CreatorID,
			tSourceDoc.State,
			tSourceDoc.Closed,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tTargetDoc.ID,
			tTargetDoc.CreatedAt,
			tTargetDoc.UpdatedAt,
			tTargetDoc.CategoryID,
			tTargetDoc.Title,
			tTargetDoc.CreatorID,
			tTargetDoc.State,
			tTargetDoc.Closed,
			tRefCreator.ID,
			tRefCreator.Job,
			tRefCreator.JobGrade,
			tRefCreator.Firstname,
			tRefCreator.Lastname,
			tRefCreator.Dateofbirth,
		).
		FROM(
			tDocRef.
				LEFT_JOIN(tSourceDoc,
					tDocRef.SourceDocumentID.EQ(tSourceDoc.ID),
				).
				LEFT_JOIN(tTargetDoc,
					tDocRef.TargetDocumentID.EQ(tTargetDoc.ID),
				).
				LEFT_JOIN(tCreator,
					tSourceDoc.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tRefCreator,
					tDocRef.CreatorID.EQ(tRefCreator.ID),
				),
		).
		WHERE(
			jet.AND(
				tDocRef.DeletedAt.IS_NULL(),
				jet.OR(
					tDocRef.SourceDocumentID.EQ(jet.Uint64(req.DocumentId)),
					tDocRef.TargetDocumentID.EQ(jet.Uint64(req.DocumentId)),
				),
			),
		).
		ORDER_BY(
			tDocRef.CreatedAt.DESC(),
		).
		LIMIT(25)

	var dest []*documents.DocumentReference
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(dest); i++ {
		if dest[i].Creator != nil {
			jobInfoFn(dest[i].Creator)
		}

		s.enricher.EnrichCategory(dest[i].SourceDocument)
		if dest[i].SourceDocument.Creator != nil {
			jobInfoFn(dest[i].SourceDocument.Creator)
		}

		s.enricher.EnrichCategory(dest[i].TargetDocument)
	}

	resp.References = dest

	return resp, nil
}

func (s *Server) GetDocumentRelations(ctx context.Context, req *GetDocumentRelationsRequest) (*GetDocumentRelationsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.DocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRelsViewDenied
	}

	relations, err := s.getDocumentRelations(ctx, userInfo, req.DocumentId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	return &GetDocumentRelationsResponse{
		Relations: relations,
	}, nil
}

func (s *Server) AddDocumentReference(ctx context.Context, req *AddDocumentReferenceRequest) (*AddDocumentReferenceResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.source_document_id", int64(req.Reference.SourceDocumentId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.target_document_id", int64(req.Reference.TargetDocumentId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentReference",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.Reference.SourceDocumentId == req.Reference.TargetDocumentId {
		return nil, errorsdocstore.ErrFeedRefSelf
	}

	// Check if user has access to both documents
	check, err := s.access.CanUserAccessTargets(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT,
		req.Reference.SourceDocumentId, req.Reference.TargetDocumentId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRefAddDenied
	}

	req.Reference.CreatorId = &userInfo.UserId

	docRef := table.FivenetDocumentsReferences
	stmt := docRef.
		INSERT(
			docRef.SourceDocumentID,
			docRef.Reference,
			docRef.TargetDocumentID,
			docRef.CreatorID,
		).
		VALUES(
			req.Reference.SourceDocumentId,
			req.Reference.Reference,
			req.Reference.TargetDocumentId,
			req.Reference.CreatorId,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &AddDocumentReferenceResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentReference(ctx context.Context, req *RemoveDocumentReferenceRequest) (*RemoveDocumentReferenceResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.reference_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentReference",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	var docIDs struct {
		Source uint64
		Target uint64
	}

	// Get document IDs of reference entry
	docsStmt := tDocRef.
		SELECT(
			tDocRef.SourceDocumentID.AS("source"),
			tDocRef.TargetDocumentID.AS("target"),
		).
		FROM(tDocRef).
		WHERE(tDocRef.ID.EQ(jet.Uint64(req.Id))).
		LIMIT(1)

	if err := docsStmt.QueryContext(ctx, s.db, &docIDs); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	check, err := s.access.CanUserAccessTargets(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT, docIDs.Source, docIDs.Target)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRefRemoveDenied
	}

	stmt := tDocRef.
		UPDATE(
			tDocRef.DeletedAt,
		).
		SET(
			tDocRef.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tDocRef.ID.EQ(jet.Uint64(req.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &RemoveDocumentReferenceResponse{}, nil
}

func (s *Server) AddDocumentRelation(ctx context.Context, req *AddDocumentRelationRequest) (*AddDocumentRelationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Relation.DocumentId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int("fivenet.docstore.source_user_id", int(req.Relation.SourceUserId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int("fivenet.docstore.target_user_id", int(req.Relation.TargetUserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentRelation",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Relation.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRelAddDenied
	}

	req.Relation.SourceUserId = userInfo.UserId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDocRel := table.FivenetDocumentsRelations
	stmt := tDocRel.
		INSERT(
			tDocRel.DocumentID,
			tDocRel.SourceUserID,
			tDocRel.Relation,
			tDocRel.TargetUserID,
		).
		VALUES(
			req.Relation.DocumentId,
			req.Relation.SourceUserId,
			req.Relation.Relation,
			req.Relation.TargetUserId,
		)

	var lastId int64

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		stmt := tDocRel.
			SELECT(
				tDocRel.ID.AS("id"),
			).
			FROM(tDocRel).
			WHERE(jet.AND(
				tDocRel.DocumentID.EQ(jet.Uint64(req.Relation.DocumentId)),
				tDocRel.Relation.EQ(jet.Int16(int16(req.Relation.Relation))),
				tDocRel.TargetUserID.EQ(jet.Int32(req.Relation.TargetUserId)),
			)).
			LIMIT(1)

		var dest struct {
			ID int64
		}
		if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		lastId = dest.ID
	} else {
		lastId, err = result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		// Only mention users when the relation has been created and not been "duplicated"
		if err := s.addUserActivity(ctx, tx,
			userInfo.UserId, req.Relation.TargetUserId, users.UserActivityType_USER_ACTIVITY_TYPE_MENTIONED, "DocStore.Relation", "",
			strconv.Itoa(int(req.Relation.DocumentId)), req.Relation.Relation.String()); err != nil {
			return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
		}

		if req.Relation.Relation == documents.DocRelation_DOC_RELATION_MENTIONED {
			if err := s.notifyMentionedUser(ctx, req.Relation.DocumentId, userInfo.UserId, req.Relation.TargetUserId); err != nil {
				return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &AddDocumentRelationResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentRelation(ctx context.Context, req *RemoveDocumentRelationRequest) (*RemoveDocumentRelationResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.docstore.id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentRelation",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	var docID struct {
		ID uint64
	}

	// Get document IDs of reference entry
	docsStmt := tDocRel.
		SELECT(
			tDocRel.DocumentID.AS("id"),
		).
		FROM(tDocRel).
		WHERE(tDocRel.ID.EQ(jet.Uint64(req.Id))).
		LIMIT(1)

	if err := docsStmt.QueryContext(ctx, s.db, &docID); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	check, err := s.access.CanUserAccessTarget(ctx, docID.ID, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return nil, errorsdocstore.ErrFeedRelRemoveDenied
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tDocRel.
		UPDATE(
			tDocRel.DeletedAt,
		).
		SET(
			tDocRel.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tDocRel.ID.EQ(jet.Uint64(req.Id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	rel, err := s.getDocumentRelation(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	if err := s.addUserActivity(ctx, tx,
		userInfo.UserId, rel.TargetUserId, users.UserActivityType_USER_ACTIVITY_TYPE_MENTIONED, "DocStore.Relation",
		strconv.Itoa(int(docID.ID)), "", rel.Relation.String()); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsdocstore.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &RemoveDocumentRelationResponse{}, nil
}

func (s *Server) getDocumentRelation(ctx context.Context, id uint64) (*documents.DocumentRelation, error) {
	stmt := tDocRel.
		SELECT(
			tDocRel.ID,
			tDocRel.CreatedAt,
			tDocRel.DocumentID,
			tDocRel.SourceUserID,
			tDocRel.Relation,
			tDocRel.TargetUserID,
		).
		FROM(
			tDocRel,
		).
		WHERE(
			tDocRel.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest documents.DocumentRelation
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) getDocumentRelations(ctx context.Context, userInfo *userinfo.UserInfo, documentId uint64) ([]*documents.DocumentRelation, error) {
	tSourceUser := tables.Users().AS("source_user")
	tTargetUser := tSourceUser.AS("target_user")

	stmt := tDocRel.
		SELECT(
			tDocRel.ID,
			tDocRel.CreatedAt,
			tDocRel.DocumentID,
			tDocRel.SourceUserID,
			tDocRel.Relation,
			tDocRel.TargetUserID,
			tDocument.ID,
			tDocument.CreatedAt,
			tDocument.UpdatedAt,
			tDocument.CategoryID,
			tDocument.Title,
			tDocument.CreatorID,
			tDocument.State,
			tDocument.Closed,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			tDCategory.Color,
			tDCategory.Icon,
			tSourceUser.ID,
			tSourceUser.Job,
			tSourceUser.JobGrade,
			tSourceUser.Firstname,
			tSourceUser.Lastname,
			tSourceUser.Dateofbirth,
			tTargetUser.ID,
			tTargetUser.Job,
			tTargetUser.JobGrade,
			tTargetUser.Firstname,
			tTargetUser.Lastname,
			tTargetUser.Dateofbirth,
		).
		FROM(
			tDocRel.
				LEFT_JOIN(tDocument,
					tDocument.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDCategory,
					tDocument.CategoryID.EQ(tDCategory.ID),
				).
				LEFT_JOIN(tSourceUser,
					tSourceUser.ID.EQ(tDocRel.SourceUserID),
				).
				LEFT_JOIN(tTargetUser,
					tTargetUser.ID.EQ(tDocRel.TargetUserID),
				),
		).
		WHERE(
			jet.AND(
				tDocRel.DocumentID.EQ(jet.Uint64(documentId)),
				tDocRel.DeletedAt.IS_NULL(),
			),
		).
		ORDER_BY(
			tDocRel.CreatedAt.DESC(),
		).
		LIMIT(25)

	var dest []*documents.DocumentRelation
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(dest); i++ {
		if dest[i].SourceUser != nil {
			jobInfoFn(dest[i].SourceUser)
		}
		if dest[i].TargetUser != nil {
			jobInfoFn(dest[i].TargetUser)
		}
	}

	return dest, nil
}

func (s *Server) notifyMentionedUser(ctx context.Context, documentId uint64, sourceUserId int32, targetUserId int32) error {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, targetUserId)
	if err != nil {
		return err
	}

	// Make sure target user has access to document
	check, err := s.access.CanUserAccessTarget(ctx, documentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	doc, err := s.getDocument(ctx, tDocument.ID.EQ(jet.Uint64(documentId)), userInfo, false)
	if err != nil {
		return err
	}
	if doc == nil {
		return nil
	}

	not := &notifications.Notification{
		UserId: targetUserId,
		Title: &common.TranslateItem{
			Key: "notifications.docstore.document_relation_mentioned.title",
		},
		Content: &common.TranslateItem{
			Key:        "notifications.docstore.document_relation_mentioned.content",
			Parameters: map[string]string{"title": doc.Title},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", doc.Id),
			},
			CausedBy: &users.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	if err := s.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
