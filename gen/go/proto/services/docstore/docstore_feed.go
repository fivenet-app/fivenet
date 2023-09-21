package docstore

import (
	context "context"
	"errors"
	"fmt"
	"strconv"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common"
	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tDocRef = table.FivenetDocumentsReferences.AS("documentreference")
	tDocRel = table.FivenetDocumentsRelations.AS("documentrelation")
)

func (s *Server) GetDocumentReferences(ctx context.Context, req *GetDocumentReferencesRequest) (*GetDocumentReferencesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view this document's references!")
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
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

	ids, err := s.checkIfUserHasAccessToDocIDs(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW, docIds...)
	if err != nil {
		return nil, err
	}

	if len(ids) == 0 {
		return resp, nil
	}

	dIds := make([]jet.Expression, len(ids))
	for i := 0; i < len(ids); i++ {
		dIds[i] = jet.Uint64(ids[i])
	}

	tSourceDoc := tDocs.AS("source_document")
	tTargetDoc := tDocs.AS("target_document")
	tRefCreator := tUsers.AS("ref_creator")
	tDCreator := tUsers.AS("creator")
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
			tDCreator.ID,
			tDCreator.Identifier,
			tDCreator.Job,
			tDCreator.JobGrade,
			tDCreator.Firstname,
			tDCreator.Lastname,
			tTargetDoc.ID,
			tTargetDoc.CreatedAt,
			tTargetDoc.UpdatedAt,
			tTargetDoc.CategoryID,
			tTargetDoc.Title,
			tTargetDoc.CreatorID,
			tTargetDoc.State,
			tTargetDoc.Closed,
			tRefCreator.ID,
			tRefCreator.Identifier,
			tRefCreator.Job,
			tRefCreator.JobGrade,
			tRefCreator.Firstname,
			tRefCreator.Lastname,
		).
		FROM(
			tDocRef.
				LEFT_JOIN(tSourceDoc,
					tDocRef.SourceDocumentID.EQ(tSourceDoc.ID),
				).
				LEFT_JOIN(tTargetDoc,
					tDocRef.TargetDocumentID.EQ(tTargetDoc.ID),
				).
				LEFT_JOIN(tDCreator,
					tSourceDoc.CreatorID.EQ(tDCreator.ID),
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(dest); i++ {
		if dest[i].Creator != nil {
			s.enricher.EnrichJobInfo(dest[i].Creator)
		}

		s.enricher.EnrichCategory(dest[i].SourceDocument)
		if dest[i].SourceDocument.Creator != nil {
			s.enricher.EnrichJobInfo(dest[i].SourceDocument.Creator)
		}

		s.enricher.EnrichCategory(dest[i].TargetDocument)
	}

	resp.References = dest

	return resp, nil
}

func (s *Server) GetDocumentRelations(ctx context.Context, req *GetDocumentRelationsRequest) (*GetDocumentRelationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view this document!")
	}

	relations, err := s.getDocumentRelations(ctx, req.DocumentId)
	if err != nil {
		return nil, err
	}

	return &GetDocumentRelationsResponse{
		Relations: relations,
	}, nil
}

func (s *Server) AddDocumentReference(ctx context.Context, req *AddDocumentReferenceRequest) (*AddDocumentReferenceResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentReference",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	if req.Reference.SourceDocumentId == req.Reference.TargetDocumentId {
		return nil, status.Error(codes.InvalidArgument, "You can't reference a document with itself!")
	}

	// Check if user has access to both documents
	check, err := s.checkIfUserHasAccessToDocs(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT,
		req.Reference.SourceDocumentId, req.Reference.TargetDocumentId)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add references from/to this document!")
	}

	req.Reference.CreatorId = &userInfo.UserId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &AddDocumentReferenceResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentReference(ctx context.Context, req *RemoveDocumentReferenceRequest) (*RemoveDocumentReferenceResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentReference",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

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
		return nil, err
	}

	check, err := s.checkIfUserHasAccessToDocs(ctx, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT, docIDs.Source, docIDs.Target)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
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
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &RemoveDocumentReferenceResponse{}, nil
}

func (s *Server) AddDocumentRelation(ctx context.Context, req *AddDocumentRelationRequest) (*AddDocumentRelationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentRelation",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Relation.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add relation from/to this document!")
	}

	req.Relation.SourceUserId = userInfo.UserId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	docRel := table.FivenetDocumentsRelations
	stmt := docRel.
		INSERT(
			docRel.DocumentID,
			docRel.SourceUserID,
			docRel.Relation,
			docRel.TargetUserID,
		).
		VALUES(
			req.Relation.DocumentId,
			req.Relation.SourceUserId,
			req.Relation.Relation,
			req.Relation.TargetUserId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.addUserActivity(ctx, tx,
		userInfo.UserId, req.Relation.TargetUserId, users.UserActivityType_USER_ACTIVITY_TYPE_MENTIONED, "DocStore.Relation", "",
		strconv.Itoa(int(lastId)), req.Relation.Relation.String()); err != nil {
		return nil, ErrFailedQuery
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	if req.Relation.Relation == documents.DocRelation_DOC_RELATION_MENTIONED {
		s.notifyUser(ctx, uint64(lastId), userInfo.UserId, req.Relation.TargetUserId)
	}

	return &AddDocumentRelationResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentRelation(ctx context.Context, req *RemoveDocumentRelationRequest) (*RemoveDocumentRelationResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentRelation",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

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
		return nil, err
	}

	check, err := s.checkIfUserHasAccessToDoc(ctx, docID.ID, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, err
	}
	if !check && !userInfo.SuperUser {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
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

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, ErrFailedQuery
	}

	rel, err := s.getDocumentRelation(ctx, req.Id)
	if err != nil {
		return nil, ErrFailedQuery
	}

	if err := s.addUserActivity(ctx, tx,
		userInfo.UserId, rel.TargetUserId, users.UserActivityType_USER_ACTIVITY_TYPE_MENTIONED, "DocStore.Relation",
		strconv.Itoa(int(docID.ID)), "", rel.Relation.String()); err != nil {
		return nil, ErrFailedQuery
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, ErrFailedQuery
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return &dest, nil
}

func (s *Server) getDocumentRelations(ctx context.Context, documentId uint64) ([]*documents.DocumentRelation, error) {
	uSource := tUsers.AS("source_user")
	uTarget := tUsers.AS("target_user")
	stmt := tDocRel.
		SELECT(
			tDocRel.ID,
			tDocRel.CreatedAt,
			tDocRel.DocumentID,
			tDocRel.SourceUserID,
			tDocRel.Relation,
			tDocRel.TargetUserID,
			tDocs.ID,
			tDocs.CreatedAt,
			tDocs.UpdatedAt,
			tDocs.CategoryID,
			tDocs.Title,
			tDocs.CreatorID,
			tDocs.State,
			tDocs.Closed,
			tDCategory.ID,
			tDCategory.Name,
			tDCategory.Description,
			uSource.ID,
			uSource.Identifier,
			uSource.Job,
			uSource.JobGrade,
			uSource.Firstname,
			uSource.Lastname,
			uTarget.ID,
			uTarget.Identifier,
			uTarget.Job,
			uTarget.JobGrade,
			uTarget.Firstname,
			uTarget.Lastname,
		).
		FROM(
			tDocRel.
				LEFT_JOIN(tDocs,
					tDocs.ID.EQ(tDocRel.DocumentID),
				).
				LEFT_JOIN(tDCategory,
					tDocs.CategoryID.EQ(tDCategory.ID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(tDocRel.SourceUserID),
				).
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(tDocRel.TargetUserID),
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
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(dest); i++ {
		if dest[i].SourceUser != nil {
			s.enricher.EnrichJobInfo(dest[i].SourceUser)
		}
		if dest[i].TargetUser != nil {
			s.enricher.EnrichJobInfo(dest[i].TargetUser)
		}
	}

	return dest, nil
}

func (s *Server) notifyUser(ctx context.Context, documentId uint64, sourceUserId int32, targetUserId int32) {
	userInfo, err := s.ui.GetUserInfoWithoutAccountId(ctx, targetUserId)
	if err != nil {
		return
	}

	doc, err := s.getDocument(ctx, tDocs.ID.EQ(jet.Uint64(documentId)), userInfo)
	if err != nil {
		return
	}
	if doc == nil {
		return
	}

	// TODO add source user as `CausedBy` to `Notification.Data``
	nType := string(notifi.InfoType)
	not := &notifications.Notification{
		UserId: targetUserId,
		Title: &common.TranslateItem{
			Key: "notifications.notifi.document_relation_mentioned.title",
		},
		Content: &common.TranslateItem{
			Key:        "notifications.notifi.document_relation_mentioned.content",
			Parameters: []string{doc.Title},
		},
		Type:     &nType,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", documentId),
			},
			CausedBy: &users.UserShort{
				UserId: sourceUserId,
			},
		},
	}
	s.notif.NotifyUser(ctx, not)
}
