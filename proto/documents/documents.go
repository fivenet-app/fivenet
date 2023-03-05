package documents

import (
	context "context"
	"errors"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/permissions"
	"github.com/galexrt/arpanet/query"
	"gorm.io/gorm"
)

func init() {
	permissions.RegisterPerms([]*permissions.Perm{
		{Key: "documents", Name: "View"},
		{Key: "documents", Name: "FindDocuments"},
		{Key: "documents", Name: "GetDocument"},
		{Key: "documents", Name: "CreateDocument"},
	})
}

type Server struct {
	DocumentsServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) prepareDocumentQuery(start query.IDocumentDo, user *model.User) query.IDocumentDo {
	d := query.Document
	dja := query.DocumentJobAccess
	dua := query.DocumentUserAccess

	if start == nil {
		start = d.Where()
	}
	return start.
		LeftJoin(dua, dua.DocumentID.EqCol(d.ID), dua.UserID.Eq(user.ID)).
		LeftJoin(dja, dja.DocumentID.EqCol(d.ID), dja.Name.Eq(user.Job), dja.MinimumGrade.Lte(user.JobGrade)).
		Where(
			d.Where(
				d.Where(
					d.Public.Is(true)).
					Or(d.CreatorID.Eq(user.ID)),
			).
				Or(
					d.Where(
						d.Where(
							dua.Access.IsNotNull(),
							dua.Access.Neq(model.BlockedAccessRole),
						),
					).
						Or(
							dja.Where(
								dua.Access.IsNull(),
								dja.Access.IsNotNull(),
								dja.Access.Neq(model.BlockedAccessRole),
							),
						),
				),
		).
		Order(d.CreatedAt.Desc()).
		Preload(
			d.JobAccess.On(dja.Name.Eq(user.Job)),
			d.UserAccess.On(dua.UserID.Eq(user.ID)),
		)
}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	resp := &FindDocumentsResponse{}

	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	documents, err := s.prepareDocumentQuery(nil, user).
		Find()
	if err != nil {
		return nil, err
	}

	_ = documents
	// TODO

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error) {
	resp := &GetDocumentResponse{}

	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	d := query.Document
	document, err := s.prepareDocumentQuery(d.Where(d.ID.Eq(uint(req.Id))), user).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = document
	// TODO

	return resp, nil
}

func (s *Server) CreateDocument(ctx context.Context, in *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	resp := &CreateDocumentResponse{}

	return resp, nil
}
