package docstore

import (
	context "context"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	docRef = table.ArpanetDocumentsReferences.AS("documentreference")
	docRel = table.ArpanetDocumentsRelations.AS("documentrelation")
)

func (s *Server) GetDocumentReferences(ctx context.Context, req *GetDocumentReferencesRequest) (*GetDocumentReferencesResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view this document!")
	}

	references, err := s.getDocumentReferences(ctx, req.DocumentId)
	if err != nil {
		return nil, err
	}

	return &GetDocumentReferencesResponse{
		References: references,
	}, nil
}

func (s *Server) GetDocumentRelations(ctx context.Context, req *GetDocumentRelationsRequest) (*GetDocumentRelationsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
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

func (s *Server) AddDocumentReferences(ctx context.Context, req *AddDocumentReferencesRequest) (*AddDocumentReferencesResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add references from this document!")
	}

	resp := &AddDocumentReferencesResponse{}

	// TODO

	return resp, nil
}
func (s *Server) RemoveDocumentReferences(ctx context.Context, req *RemoveDocumentReferencesRequest) (*RemoveDocumentReferencesResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
	}

	ids := make([]jet.Expression, len(req.RefIds))
	for i := 0; i < len(req.RefIds); i++ {
		ids[i] = jet.Uint64(req.RefIds[i])
	}

	stmt := docRef.DELETE().
		WHERE(
			docRef.ID.IN(ids...),
		).
		LIMIT(5)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	resp := &RemoveDocumentReferencesResponse{}

	return resp, nil
}
func (s *Server) AddDocumentRelations(ctx context.Context, req *AddDocumentRelationsRequest) (*AddDocumentRelationsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to add relations from this document!")
	}

	resp := &AddDocumentRelationsResponse{}

	// TODO

	return resp, nil
}
func (s *Server) RemoveDocumentRelations(ctx context.Context, req *RemoveDocumentRelationsRequest) (*RemoveDocumentRelationsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove relations from this document!")
	}

	ids := make([]jet.Expression, len(req.RelIds))
	for i := 0; i < len(req.RelIds); i++ {
		ids[i] = jet.Uint64(req.RelIds[i])
	}

	stmt := docRel.DELETE().
		WHERE(
			docRel.ID.IN(ids...),
		).
		LIMIT(5)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	resp := &RemoveDocumentRelationsResponse{}

	return resp, nil
}

func (s *Server) getDocumentReferences(ctx context.Context, documentID uint64) ([]*documents.DocumentReference, error) {
	sourceDoc := docs.AS("source_document")
	uCreator := u.AS("ref_creator")
	stmt := docRef.SELECT(
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
		dCategory.ID,
		dCategory.Name,
		dCategory.Description,
		uCreator.ID,
		uCreator.Identifier,
		uCreator.Job,
		uCreator.JobGrade,
		uCreator.Firstname,
		uCreator.Lastname,
	).
		FROM(
			docRef.
				LEFT_JOIN(sourceDoc,
					docRef.SourceDocumentID.EQ(sourceDoc.ID),
				).
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(sourceDoc.CategoryID),
				).
				LEFT_JOIN(uCreator,
					docRef.CreatorID.EQ(uCreator.ID),
				),
		).
		WHERE(
			docRef.TargetDocumentID.EQ(jet.Uint64(documentID)),
		)

	var dest struct {
		References []*documents.DocumentReference
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest.References, nil
}

func (s *Server) getDocumentRelations(ctx context.Context, documentID uint64) ([]*documents.DocumentRelation, error) {
	uSource := u.AS("source_user")
	uTarget := u.AS("target_user")
	stmt := docRel.SELECT(
		docRel.ID,
		docRel.CreatedAt,
		docRel.DocumentID,
		docRel.SourceUserID,
		docRel.Relation,
		docRel.TargetUserID,
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
					docRel.DocumentID.EQ(docs.ID),
				).
				LEFT_JOIN(uSource,
					docRel.SourceUserID.EQ(uSource.ID),
				).
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(docs.CategoryID),
				).
				LEFT_JOIN(uTarget,
					docRel.TargetUserID.EQ(uTarget.ID),
				),
		).
		WHERE(
			docRel.DocumentID.EQ(jet.Uint64(documentID)),
		)

	var dest struct {
		Relations []*documents.DocumentRelation
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest.Relations, nil
}
