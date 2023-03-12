package documents

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	database "github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/microcosm-cc/bluemonday"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "documents", Name: "View"},
		{Key: "documents", Name: "FindDocuments"},
		{Key: "documents", Name: "GetDocument"},
		{Key: "documents", Name: "CreateDocument"},
		{Key: "documents", Name: "UpdateDocument"},
		{Key: "documents", Name: "CompleteCategories", PerJob: true},
	})
}

var (
	u   = table.Users
	d   = table.ArpanetDocuments.AS("document")
	dua = table.ArpanetDocumentsUserAccess
	dja = table.ArpanetDocumentsJobAccess
	dt  = table.ArpanetDocumentsTemplates
)

type Server struct {
	DocumentsServiceServer

	p *bluemonday.Policy
}

func NewServer() *Server {
	return &Server{
		p: bluemonday.UGCPolicy(),
	}
}

func (s *Server) getDocumentsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, additionalColumns jet.ProjectionList, userID int32, job string, jobGrade int32) jet.SelectStatement {
	wheres := []jet.BoolExpression{jet.OR(
		jet.OR(
			d.Public.IS_TRUE(),
			d.CreatorID.EQ(jet.Int32(userID)),
		),
		jet.OR(
			jet.AND(
				dua.Access.IS_NOT_NULL(),
				dua.Access.NOT_EQ(jet.Int32(int32(documents.DOCUMENT_ACCESS_BLOCKED))),
			),
			jet.AND(
				dua.Access.IS_NULL(),
				dja.Access.IS_NOT_NULL(),
				dja.Access.NOT_EQ(jet.Int32(int32(documents.DOCUMENT_ACCESS_BLOCKED))),
			),
		),
	)}

	if where != nil {
		wheres = append(wheres, where)
	}

	u := u.AS("creator")
	var q jet.SelectStatement
	if onlyColumns != nil {
		q = d.SELECT(
			onlyColumns,
		)
	} else {
		if additionalColumns == nil {
			q = d.SELECT(
				d.AllColumns,
				u.ID,
				u.Identifier,
				u.Job,
				u.JobGrade,
				u.Firstname,
				u.Lastname,
			)
		} else {
			additionalColumns = append(jet.ProjectionList{
				u.ID,
				u.Identifier,
				u.Job,
				u.JobGrade,
				u.Firstname,
				u.Lastname,
			}, additionalColumns)
			q = d.SELECT(
				d.AllColumns,
				additionalColumns...,
			)
		}
	}

	return q.
		FROM(
			d.LEFT_JOIN(dua,
				dua.DocumentID.EQ(d.ID).
					AND(dua.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(dja,
					dja.DocumentID.EQ(d.ID).
						AND(dja.Job.EQ(jet.String(job))).
						AND(dja.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				).
				LEFT_JOIN(u, u.ID.EQ(jet.Int32(userID))),
		).WHERE(
		jet.AND(
			wheres...,
		),
	).
		ORDER_BY(d.CreatedAt.DESC()).
		LIMIT(database.DefaultPageLimit)

}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "FindDocuments") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to list documents!")
	}

	countStmt := s.getDocumentsQuery(d.ResponseID.IS_NULL(), jet.ProjectionList{jet.COUNT(d.ID).AS("total_count")}, nil, userID, job, jobGrade)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, query.DB, &count); err != nil {
		return nil, err
	}

	resp := &FindDocumentsResponse{
		Offset:     req.Offset,
		TotalCount: count.TotalCount,
		End:        0,
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.getDocumentsQuery(d.ResponseID.IS_NULL(), nil, nil, userID, job, jobGrade)
	if err := stmt.QueryContext(ctx, query.DB, &resp.Documents); err != nil {
		return nil, err
	}

	resp.TotalCount = count.TotalCount
	if req.Offset >= resp.TotalCount {
		resp.Offset = 0
	} else {
		resp.Offset = req.Offset
	}
	resp.End = resp.Offset + int64(len(resp.Documents))

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "GetDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to get a document!")
	}

	resp := &GetDocumentResponse{
		Document: &documents.Document{},
	}
	stmt := s.getDocumentsQuery(jet.AND(
		d.ResponseID.IS_NULL(),
		d.ID.EQ(jet.Uint64(req.Id)),
	), nil, nil, userID, job, jobGrade).
		LIMIT(1)
	if err := stmt.QueryContext(ctx, query.DB, resp.Document); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return resp, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	userID, job, _ := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "CreateDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to create documents!")
	}

	stmt := d.INSERT(
		d.Title,
		d.Content,
		d.ContentType,
		d.Closed,
		d.State,
		d.Public,
		d.CreatorID,
		d.CreatorJob,
	).VALUES(
		req.Title,
		s.p.Sanitize(req.Content),
		documents.DOCUMENT_CONTENT_TYPE_HTML.String(),
		req.Closed,
		req.State,
		req.Public,
		userID,
		job,
	)

	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccess(req.JobsAccess, req.UsersAccess); err != nil {
		return nil, err
	}

	return &CreateDocumentResponse{}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "UpdateDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit documents!")
	}

	checkStmt := d.SELECT(
		d.ID,
	).
		FROM(
			d.LEFT_JOIN(dua,
				dua.DocumentID.EQ(d.ID).
					AND(dua.UserID.EQ(jet.Int32(userID)))).
				LEFT_JOIN(dja,
					dja.DocumentID.EQ(d.ID).
						AND(dja.Job.EQ(jet.String(job))).
						AND(dja.MinimumGrade.LT_EQ(jet.Int32(jobGrade))),
				).
				LEFT_JOIN(u, u.ID.EQ(jet.Int32(userID))),
		).WHERE(
		jet.OR(
			d.CreatorID.EQ(jet.Int32(userID)),
			jet.AND(
				dua.Access.IS_NOT_NULL(),
				dua.Access.IN(
					jet.Int32(int32(documents.DOCUMENT_ACCESS_EDIT)),
					jet.Int32(int32(documents.DOCUMENT_ACCESS_ADMIN)),
				),
			),
			jet.AND(
				dua.Access.IS_NULL(),
				dja.Access.IS_NOT_NULL(),
				dja.Access.IN(
					jet.Int32(int32(documents.DOCUMENT_ACCESS_EDIT)),
					jet.Int32(int32(documents.DOCUMENT_ACCESS_LEADER)),
					jet.Int32(int32(documents.DOCUMENT_ACCESS_ADMIN)),
				),
			),
		),
	).
		LIMIT(1)

	var dest struct {
		ID uint64 `alias:"document.id"`
	}
	if err := checkStmt.QueryContext(ctx, query.DB, &dest); err != nil {
		return nil, err
	}

	if dest.ID <= 0 {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this document!")
	}

	resp := &UpdateDocumentResponse{}

	stmt := d.UPDATE(
		d.Title,
		d.Content,
		d.Closed,
		d.State,
		d.Public,
	).
		SET(
			req.Title,
			req.Content,
			req.Closed,
			req.State,
			req.Public,
		).
		WHERE(d.ID.EQ(jet.Uint64(req.Id)))

	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "CreateDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to list/ get document templates!")
	}

	resp := &ListTemplatesResponse{}
	stmt := dt.SELECT(
		dt.ID,
		dt.Job,
		dt.JobGrade,
		dt.Title,
		dt.Description,
		dt.CreatorID,
	).
		FROM(dt).
		WHERE(
			jet.AND(
				dt.Job.EQ(jet.String(job)),
				dt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	if err := stmt.QueryContext(ctx, query.DB, &resp.Templates); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	if !perms.P.CanID(userID, "documents", "CreateDocument") {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to list/ get document templates!")
	}

	resp := &GetTemplateResponse{}
	stmt := dt.SELECT(
		dt.AllColumns,
	).
		FROM(dt).
		WHERE(
			jet.AND(
				dt.ID.EQ(jet.Uint64(req.Id)),
				dt.Job.EQ(jet.String(job)),
				dt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	if err := stmt.QueryContext(ctx, query.DB, &resp.Template); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	resp := &GetDocumentAccessResponse{}

	stmt := d.SELECT(
		dua.AllColumns,
		dja.AllColumns,
	).
		FROM(d).
		WHERE(
			jet.AND(
				d.ID.EQ(jet.Uint64(req.Id)),
				jet.OR(
					jet.AND(
						dua.Access.IS_NOT_NULL(),
						dua.Access.NOT_EQ(jet.Int32(int32(documents.DOCUMENT_ACCESS_BLOCKED))),
					),
					jet.AND(
						dua.Access.IS_NULL(),
						dja.Access.IS_NOT_NULL(),
						dja.Access.NOT_EQ(jet.Int32(int32(documents.DOCUMENT_ACCESS_BLOCKED))),
					),
				),
			),
		)

	if err := stmt.QueryContext(ctx, query.DB, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *SetDocumentAccessRequest) (*SetDocumentAccessResponse, error) {
	resp := &SetDocumentAccessResponse{}

	// TODO add/update/remove for document access is needed

	return resp, nil
}

func (s *Server) handleDocumentAccess(ja []*documents.DocumentJobAccess, ua []*documents.DocumentUserAccess) error {

	// TODO need to create job and user access

	return nil
}
