package documents

import (
	context "context"
	"errors"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/modelhelper"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

var (
	d   = table.ArpanetDocuments.AS("document")
	dua = table.ArpanetDocumentsUserAccess
	dja = table.ArpanetDocumentsJobAccess
)

type Server struct {
	DocumentsServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) getDocumentsQuery(where jet.BoolExpression, additionalColumns jet.ProjectionList, userID int32, job string, jobGrade int32) jet.SelectStatement {
	wheres := []jet.BoolExpression{jet.OR(
		jet.OR(
			d.Public.IS_TRUE(),
			d.CreatorID.EQ(jet.Int32(userID)),
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
	)}

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if additionalColumns == nil {
		q = d.SELECT(
			d.AllColumns,
		)
	} else {
		q = d.SELECT(
			d.AllColumns,
			additionalColumns...,
		)
	}

	return q.
		FROM(
			d.LEFT_JOIN(dua,
				dua.DocumentID.EQ(d.ID).
					AND(dua.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(dja,
					dja.DocumentID.EQ(d.ID).
						AND(dja.Name.EQ(jet.String(job))).
						AND(dja.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				),
		).WHERE(
		jet.AND(
			wheres...,
		),
	).
		ORDER_BY(d.CreatedAt.DESC())

}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "FindDocuments") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to list documents!")
	}

	resp := &FindDocumentsResponse{}
	stmt := s.getDocumentsQuery(nil, nil, userID, job, jobGrade)
	if err := stmt.QueryContext(ctx, query.DB, &resp.Documents); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "GetDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to get a document!")
	}

	resp := &GetDocumentResponse{
		Document: &Document{},
	}
	stmt := s.getDocumentsQuery(d.ID.EQ(jet.Uint64(req.Id)), nil, userID, job, jobGrade).
		LIMIT(1)
	if err := stmt.QueryContext(ctx, query.DB, resp.Document); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// TODO Load responses

	return resp, nil
}

func (s *Server) CreateOrEditDocument(ctx context.Context, in *CreateOrEditDocumentRequest) (*CreateOrEditDocumentResponse, error) {
	userID, _, _ := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "CreateDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to create/ edit a document!")
	}

	resp := &CreateOrEditDocumentResponse{}

	// TODO

	return resp, nil
}

// TODO add/update/remove for document access is needed
