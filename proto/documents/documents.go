package documents

import (
	context "context"
	"errors"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/modelhelper"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/common"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"gorm.io/gorm"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
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

func (s *Server) getDocumentsQuery(where jet.BoolExpression, user *common.ShortUser) jet.SelectStatement {
	d := table.ArpanetDocuments
	dua := table.ArpanetDocumentsUserAccess
	dja := table.ArpanetDocumentsJobAccess

	return d.SELECT(
		d.AllColumns,
		dja.AllColumns,
	).
		FROM(
			d.LEFT_JOIN(dua,
				dua.DocumentID.EQ(d.ID).
					AND(dua.UserID.EQ(jet.Int32(user.UserID)))),
			d.LEFT_JOIN(dja,
				dja.DocumentID.EQ(d.ID).
					AND(dja.Name.EQ(jet.String(user.Job))).
					AND(dja.MinimumGrade.LT_EQ(jet.Int32(user.JobGrade))),
			),
		).WHERE(
		jet.AND(
			jet.OR(
				jet.OR(
					d.Public.IS_TRUE(),
					d.CreatorID.EQ(jet.Int32(user.UserID)),
				),
				jet.OR(
					jet.AND(
						dua.Access.IS_NOT_NULL(),
						dua.Access.NOT_EQ(jet.String(modelhelper.BlockedAccessRole)),
					),
					jet.AND(
						dua.Access.IS_NULL(),
						dja.Access.IS_NOT_NULL(),
						dja.Access.NOT_EQ(jet.String(modelhelper.BlockedAccessRole)),
					),
				),
			),
			where,
		),
	).
		ORDER_BY(d.CreatedAt.DESC())

}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	resp := &FindDocumentsResponse{}

	user, err := auth.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var documents []struct {
		model.ArpanetDocuments

		model.ArpanetDocumentsJobAccess
		model.ArpanetDocumentsUserAccess
	}
	if err := s.getDocumentsQuery(nil, user).QueryContext(ctx, query.DB, &documents); err != nil {
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

	d := table.ArpanetDocuments
	var documents []model.ArpanetDocuments

	if err := s.getDocumentsQuery(d.ID.EQ(jet.Uint64(req.Id)), user).
		QueryContext(ctx, query.DB, &documents); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	_ = documents
	// TODO

	return resp, nil
}

func (s *Server) CreateDocument(ctx context.Context, in *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	resp := &CreateDocumentResponse{}

	return resp, nil
}
