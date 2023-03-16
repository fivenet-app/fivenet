package docstore

import (
	context "context"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/proto/resources/documents"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetDocumentReferences(ctx context.Context, req *GetDocumentReferencesRequest) (*GetDocumentReferencesResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
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
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
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
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
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
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove references from this document!")
	}

	resp := &RemoveDocumentReferencesResponse{}

	// TODO

	return resp, nil
}
func (s *Server) AddDocumentRelations(ctx context.Context, req *AddDocumentRelationsRequest) (*AddDocumentRelationsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
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
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to remove relations from this document!")
	}

	resp := &RemoveDocumentRelationsResponse{}

	// TODO

	return resp, nil
}

func (s *Server) getDocumentReferences(ctx context.Context, documentID uint64) ([]*documents.DocumentReference, error) {
	sourceDoc := ad.AS("source_document")
	uCreator := u.AS("ref_creator")
	stmt := adref.SELECT(
		adref.ID,
		adref.CreatedAt,
		adref.SourceDocumentID,
		adref.Reference,
		adref.TargetDocumentID,
		adref.CreatorID,
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
			adref.
				LEFT_JOIN(sourceDoc,
					adref.SourceDocumentID.EQ(sourceDoc.ID),
				).
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(sourceDoc.CategoryID),
				).
				LEFT_JOIN(uCreator,
					adref.CreatorID.EQ(uCreator.ID),
				),
		).
		WHERE(
			adref.TargetDocumentID.EQ(jet.Uint64(documentID)),
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
	stmt := adrel.SELECT(
		adrel.ID,
		adrel.CreatedAt,
		adrel.DocumentID,
		adrel.SourceUserID,
		adrel.Relation,
		adrel.TargetUserID,
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
			adrel.
				LEFT_JOIN(ad,
					adrel.DocumentID.EQ(ad.ID),
				).
				LEFT_JOIN(uSource,
					adrel.SourceUserID.EQ(uSource.ID),
				).
				LEFT_JOIN(dCategory,
					dCategory.ID.EQ(ad.CategoryID),
				).
				LEFT_JOIN(uTarget,
					adrel.TargetUserID.EQ(uTarget.ID),
				),
		).
		WHERE(
			adrel.DocumentID.EQ(jet.Uint64(documentID)),
		)

	var dest struct {
		Relations []*documents.DocumentRelation
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest.Relations, nil
}
