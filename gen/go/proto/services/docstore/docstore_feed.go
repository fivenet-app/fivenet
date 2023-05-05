package docstore

import (
	context "context"
	"errors"

	"github.com/galexrt/fivenet/gen/go/proto/resources/documents"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	docRef = table.FivenetDocumentsReferences.AS("documentreference")
	docRel = table.FivenetDocumentsRelations.AS("documentrelation")
)

func (s *Server) GetDocumentReferences(ctx context.Context, req *GetDocumentReferencesRequest) (*GetDocumentReferencesResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, true, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view this document's references!")
	}

	resp := &GetDocumentReferencesResponse{}

	var docsIds []struct {
		Source *uint64
		Target *uint64
	}
	idStmt := docRef.
		SELECT(
			docRef.SourceDocumentID.AS("source"),
			docRef.TargetDocumentID.AS("target"),
		).
		FROM(
			docRef,
		).
		WHERE(
			jet.AND(
				docRef.DeletedAt.IS_NULL(),
				jet.OR(
					docRef.SourceDocumentID.EQ(jet.Uint64(req.DocumentId)),
					docRef.TargetDocumentID.EQ(jet.Uint64(req.DocumentId)),
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

	ids, err := s.checkIfUserHasAccessToDocIDs(ctx, userId, job, jobGrade, true, documents.DOC_ACCESS_VIEW, docIds...)
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

	sourceDoc := docs.AS("source_document")
	targetDoc := docs.AS("target_document")
	refCreator := user.AS("ref_creator")
	dCreator := user.AS("creator")
	stmt := docRef.
		SELECT(
			docRef.ID,
			docRef.CreatedAt,
			docRef.SourceDocumentID,
			docRef.Reference,
			docRef.TargetDocumentID,
			docRef.CreatorID,
			sourceDoc.ID,
			sourceDoc.CreatedAt,
			sourceDoc.UpdatedAt,
			sourceDoc.CategoryID,
			sourceDoc.Title,
			sourceDoc.CreatorID,
			sourceDoc.State,
			sourceDoc.Closed,
			dCreator.ID,
			dCreator.Identifier,
			dCreator.Job,
			dCreator.JobGrade,
			dCreator.Firstname,
			dCreator.Lastname,
			targetDoc.ID,
			targetDoc.CreatedAt,
			targetDoc.UpdatedAt,
			targetDoc.CategoryID,
			targetDoc.Title,
			targetDoc.CreatorID,
			targetDoc.State,
			targetDoc.Closed,
			refCreator.ID,
			refCreator.Identifier,
			refCreator.Job,
			refCreator.JobGrade,
			refCreator.Firstname,
			refCreator.Lastname,
		).
		FROM(
			docRef.
				LEFT_JOIN(sourceDoc,
					docRef.SourceDocumentID.EQ(sourceDoc.ID),
				).
				LEFT_JOIN(targetDoc,
					docRef.TargetDocumentID.EQ(targetDoc.ID),
				).
				LEFT_JOIN(dCreator,
					sourceDoc.CreatorID.EQ(dCreator.ID),
				).
				LEFT_JOIN(refCreator,
					docRef.CreatorID.EQ(refCreator.ID),
				),
		).
		WHERE(
			jet.AND(
				docRef.DeletedAt.IS_NULL(),
				jet.OR(
					docRef.SourceDocumentID.EQ(jet.Uint64(req.DocumentId)),
					docRef.TargetDocumentID.EQ(jet.Uint64(req.DocumentId)),
				),
			),
		).
		ORDER_BY(
			docRef.CreatedAt.DESC(),
		).
		LIMIT(25)

	var dest []*documents.DocumentReference
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(dest); i++ {
		s.c.EnrichJobInfo(dest[i].Creator)
		s.c.EnrichDocumentCategory(dest[i].SourceDocument)
		s.c.EnrichJobInfo(dest[i].SourceDocument.Creator)
		s.c.EnrichDocumentCategory(dest[i].TargetDocument)
	}

	resp.References = dest

	return resp, nil
}

func (s *Server) GetDocumentRelations(ctx context.Context, req *GetDocumentRelationsRequest) (*GetDocumentRelationsResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, true, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
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
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentReference",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if req.Reference.SourceDocumentId == req.Reference.TargetDocumentId {
		return nil, status.Error(codes.InvalidArgument, "You can't reference a document with itself!")
	}

	// Check if user has access to both documents
	check, err := s.checkIfUserHasAccessToDocs(ctx, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT,
		req.Reference.SourceDocumentId, req.Reference.TargetDocumentId)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add references from/to this document!")
	}

	req.Reference.CreatorId = userId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FailedQueryErr
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
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return &AddDocumentReferenceResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentReference(ctx context.Context, req *RemoveDocumentReferenceRequest) (*RemoveDocumentReferenceResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentReference",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	var docIDs struct {
		Source uint64
		Target uint64
	}

	// Get document IDs of reference entry
	docsStmt := docRef.
		SELECT(
			docRef.SourceDocumentID.AS("source"),
			docRef.TargetDocumentID.AS("target"),
		).
		FROM(docRef).
		WHERE(docRef.ID.EQ(jet.Uint64(req.Id))).
		LIMIT(1)

	if err := docsStmt.QueryContext(ctx, s.db, &docIDs); err != nil {
		return nil, err
	}

	check, err := s.checkIfUserHasAccessToDocs(ctx, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT, docIDs.Source, docIDs.Target)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
	}

	stmt := docRef.
		UPDATE(
			docRef.DeletedAt,
		).
		SET(
			docRef.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			docRef.ID.EQ(jet.Uint64(req.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &RemoveDocumentReferenceResponse{}, nil
}

func (s *Server) AddDocumentRelation(ctx context.Context, req *AddDocumentRelationRequest) (*AddDocumentRelationResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "AddDocumentRelation",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	check, err := s.checkIfUserHasAccessToDoc(ctx, req.Relation.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add relation from/to this document!")
	}

	req.Relation.SourceUserId = userId

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FailedQueryErr
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

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	s.notifyUser(ctx, uint64(lastId), req.Relation.TargetUserId)

	return &AddDocumentRelationResponse{
		Id: uint64(lastId),
	}, nil
}

func (s *Server) RemoveDocumentRelation(ctx context.Context, req *RemoveDocumentRelationRequest) (*RemoveDocumentRelationResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: DocStoreService_ServiceDesc.ServiceName,
		Method:  "RemoveDocumentRelation",
		UserID:  userId,
		UserJob: job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	var docID struct {
		ID uint64
	}

	// Get document IDs of reference entry
	docsStmt := docRel.
		SELECT(
			docRel.DocumentID.AS("id"),
		).
		FROM(docRel).
		WHERE(docRel.ID.EQ(jet.Uint64(req.Id))).
		LIMIT(1)

	if err := docsStmt.QueryContext(ctx, s.db, &docID); err != nil {
		return nil, err
	}

	check, err := s.checkIfUserHasAccessToDoc(ctx, docID.ID, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
	}

	stmt := docRel.
		UPDATE(
			docRel.DeletedAt,
		).
		SET(
			docRel.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			docRel.ID.EQ(jet.Uint64(req.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	return &RemoveDocumentRelationResponse{}, nil
}

func (s *Server) getDocumentRelations(ctx context.Context, documentId uint64) ([]*documents.DocumentRelation, error) {
	uSource := user.AS("source_user")
	uTarget := user.AS("target_user")
	stmt := docRel.
		SELECT(
			docRel.ID,
			docRel.CreatedAt,
			docRel.DocumentID,
			docRel.SourceUserID,
			docRel.Relation,
			docRel.TargetUserID,
			docs.ID,
			docs.CreatedAt,
			docs.UpdatedAt,
			docs.CategoryID,
			docs.Title,
			docs.CreatorID,
			docs.State,
			docs.Closed,
			dCategory.ID,
			dCategory.Name,
			dCategory.Description,
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
			docRel.
				LEFT_JOIN(docs,
					docs.ID.EQ(docRel.DocumentID),
				).
				LEFT_JOIN(dCategory,
					docs.CategoryID.EQ(dCategory.ID),
				).
				LEFT_JOIN(uSource,
					uSource.ID.EQ(docRel.SourceUserID),
				).
				LEFT_JOIN(uTarget,
					uTarget.ID.EQ(docRel.TargetUserID),
				),
		).
		WHERE(
			jet.AND(
				docRel.DocumentID.EQ(jet.Uint64(documentId)),
				docRel.DeletedAt.IS_NULL(),
			),
		).
		ORDER_BY(
			docRel.CreatedAt.DESC(),
		).
		LIMIT(25)

	var dest []*documents.DocumentRelation
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	for i := 0; i < len(dest); i++ {
		s.c.EnrichJobInfo(dest[i].SourceUser)
		s.c.EnrichJobInfo(dest[i].TargetUser)
	}

	return dest, nil
}

func (s *Server) notifyUser(ctx context.Context, documentId uint64, targetUserId int32) {
	// check, err := s.checkIfUserHasAccessToDoc(ctx, documentId, targetUserId, job, jobGrade, false, documents.DOC_ACCESS_VIEW)
	// if err != nil {
	// 	return
	// }
	// if !check {
	// 	return
	// }

	// s.n.Add(targetUserId, "TITLE", "CONTENT", "info")
}
