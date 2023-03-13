package docstore

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/htmlsanitizer"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	u    = table.Users
	adt  = table.ArpanetDocumentsTemplates
	ad   = table.ArpanetDocuments.AS("document")
	adc  = table.ArpanetDocumentsComments
	adua = table.ArpanetDocumentsUserAccess
	adja = table.ArpanetDocumentsJobAccess
	dc   = table.ArpanetDocumentsCategories.AS("document_category")
)

type Server struct {
	DocStoreServiceServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) ListTemplates(ctx context.Context, req *ListTemplatesRequest) (*ListTemplatesResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := adt.SELECT(
		adt.ID,
		adt.Job,
		adt.JobGrade,
		adt.Title,
		adt.Description,
		adt.CreatorID,
	).
		FROM(adt).
		WHERE(
			jet.AND(
				adt.Job.EQ(jet.String(job)),
				adt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &ListTemplatesResponse{}
	if err := stmt.QueryContext(ctx, query.DB, &resp.Templates); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetTemplate(ctx context.Context, req *GetTemplateRequest) (*GetTemplateResponse, error) {
	_, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	stmt := adt.SELECT(
		adt.AllColumns,
	).
		FROM(adt).
		WHERE(
			jet.AND(
				adt.ID.EQ(jet.Uint64(req.TemplateId)),
				adt.Job.EQ(jet.String(job)),
				adt.JobGrade.LT_EQ(jet.Int32(jobGrade)),
			),
		)

	resp := &GetTemplateResponse{}
	if err := stmt.QueryContext(ctx, query.DB, &resp.Template); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	condition := jet.BoolExp(jet.Raw("MATCH(title) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.Search}))
	countStmt := s.getDocumentsQuery(
		condition,
		jet.ProjectionList{jet.COUNT(ad.ID).AS("total_count")},
		nil, userID, job, jobGrade)
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

	stmt := s.getDocumentsQuery(condition, nil, nil, userID, job, jobGrade)
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

	condition := jet.OR(
		ad.ID.EQ(jet.Uint64(req.DocumentId)),
	)

	countStmt := s.getDocumentsQuery(condition, jet.ProjectionList{jet.COUNT(ad.ID).AS("total_count")}, nil, userID, job, jobGrade)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, query.DB, &count); err != nil {
		return nil, err
	}

	resp := &GetDocumentResponse{
		Document:    &documents.Document{},
		JobsAccess:  []*documents.DocumentJobAccess{},
		UsersAccess: []*documents.DocumentUserAccess{},
	}

	stmt := s.getDocumentsQuery(condition, nil, nil, userID, job, jobGrade)
	if err := stmt.QueryContext(ctx, query.DB, &resp.Document); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	docAccess, err := s.GetDocumentAccess(ctx, &GetDocumentAccessRequest{
		DocumentId: resp.Document.Id,
	})
	if err != nil {
		return nil, err
	}

	resp.JobsAccess = docAccess.Jobs
	resp.UsersAccess = docAccess.Users

	return resp, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	userID, job, _ := auth.GetUserInfoFromContext(ctx)

	ad := table.ArpanetDocuments
	stmt := ad.INSERT(
		ad.Title,
		ad.Content,
		ad.ContentType,
		ad.Closed,
		ad.State,
		ad.Public,
		ad.CreatorID,
		ad.CategoryID,
	).VALUES(
		req.Title,
		htmlsanitizer.Sanitize(req.Content),
		documents.DOC_CONTENT_TYPE_HTML,
		req.Closed,
		req.State,
		req.Public,
		userID,
		job,
		req.CategoryId,
	)

	result, err := stmt.ExecContext(ctx, query.DB)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccess(ctx, DOC_ACCESS_UPDATE_MODE_REPLACE, uint64(lastID), req.JobsAccess, req.UsersAccess); err != nil {
		return nil, err
	}

	return &CreateDocumentResponse{
		Id: uint64(lastID),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserCanEditDocument(ctx, userID, job, jobGrade)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this document!")
	}

	stmt := ad.UPDATE(
		ad.Title,
		ad.Content,
		ad.Closed,
		ad.State,
		ad.Public,
	).
		SET(
			req.Title,
			req.Content,
			req.Closed,
			req.State,
			req.Public,
		).
		WHERE(ad.ID.EQ(jet.Uint64(req.DocumentId)))

	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	return &UpdateDocumentResponse{}, nil
}

func (s *Server) GetDocumentComments(ctx context.Context, req *GetDocumentCommentsRequest) (*GetDocumentCommentsResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document comments!")
	}

	condition := adc.DocumentID.EQ(jet.Uint64(req.DocumentID))
	countStmt := adc.SELECT(
		jet.COUNT(ad.ID).AS("total_count"),
	).
		FROM(
			adc,
		).
		WHERE(condition)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, query.DB, &count); err != nil {
		return nil, err
	}

	stmt := adc.SELECT(
		adc.ID,
		adc.Comment,
		adc.CreatorID,
		u.ID,
		u.Identifier,
		u.Job,
		u.JobGrade,
		u.Firstname,
		u.Lastname,
	).
		FROM(
			adc.
				LEFT_JOIN(u,
					adc.CreatorID.EQ(u.ID),
				),
		).
		WHERE(condition)

	resp := &GetDocumentCommentsResponse{}
	if err := stmt.QueryContext(ctx, query.DB, resp.Comments); err != nil {
		return nil, err
	}

	resp.TotalCount = count.TotalCount
	if req.Offset >= resp.TotalCount {
		resp.Offset = 0
	} else {
		resp.Offset = req.Offset
	}
	resp.End = resp.Offset + int64(len(resp.Comments))

	return resp, nil
}

func (s *Server) PostDocumentComment(ctx context.Context, req *PostDocumentCommentRequest) (*PostDocumentCommentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to post a comment on this document!")
	}

	// Clean comment from
	req.Comment.Comment = htmlsanitizer.StripTags(req.Comment.Comment)

	stmt := adc.INSERT(
		adc.DocumentID,
		adc.Comment,
		adc.CreatorID,
	).
		VALUES(
			req.Comment.DocumentId,
			req.Comment.Comment,
			userID,
		)

	resp := &PostDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	return resp, nil
}
func (s *Server) EditDocumentComment(ctx context.Context, req *EditDocumentCommentRequest) (*EditDocumentCommentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_VIEW)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this comment!")
	}

	stmt := adc.UPDATE().
		SET(
			adc.Comment.SET(jet.String(req.Comment.Comment)),
		).
		WHERE(
			adc.ID.EQ(jet.Uint64(req.Comment.Id)),
		)

	resp := &EditDocumentCommentResponse{}
	if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserCanEditDocument(ctx, userID, job, jobGrade)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document access!")
	}

	stmt := ad.SELECT(
		adua.AllColumns,
		adja.AllColumns,
	).
		FROM(
			ad.
				LEFT_JOIN(adua,
					ad.ID.EQ(adua.DocumentID)).
				LEFT_JOIN(adja,
					ad.ID.EQ(adua.DocumentID)),
		).
		WHERE(
			jet.AND(
				ad.ID.EQ(jet.Uint64(req.DocumentId)),
				jet.OR(
					jet.AND(
						adua.Access.IS_NOT_NULL(),
						adua.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
					jet.AND(
						adua.Access.IS_NULL(),
						adja.Access.IS_NOT_NULL(),
						adja.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
				),
			),
		)

	resp := &GetDocumentAccessResponse{}
	if err := stmt.QueryContext(ctx, query.DB, resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *SetDocumentAccessRequest) (*SetDocumentAccessResponse, error) {
	resp := &SetDocumentAccessResponse{}

	if err := s.handleDocumentAccess(ctx, req.Mode, req.DocumentId, req.Jobs, req.Users); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) handleDocumentAccess(ctx context.Context, mode DOC_ACCESS_UPDATE_MODE, documentID uint64, ja []*documents.DocumentJobAccess, ua []*documents.DocumentUserAccess) error {
	// Select existing job and user accesses
	var dest struct {
		jobs  []*documents.DocumentJobAccess
		users []*documents.DocumentUserAccess
	}
	selectStmt := jet.SELECT(
		adja.AllColumns,
		adua.AllColumns,
	).
		FROM(
			adja,
			adua,
		).
		WHERE(
			adja.DocumentID.EQ(jet.Uint64(documentID)),
		)

	fmt.Println(selectStmt.DebugSql())

	if err := selectStmt.QueryContext(ctx, query.DB, &dest); err != nil && errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// TODO add/update/remove for document access based on the current access in the database

	// Create accesses
	if len(ja) > 0 {
		for k := 0; k < len(ja); k++ {
			ja[k].DocumentId = documentID
		}

		// Create document job access
		stmt := adja.INSERT(
			adja.DocumentID,
			adja.Job,
			adja.MinimumGrade,
			adja.Access,
		).
			MODELS(ja)
		fmt.Println(stmt.DebugSql())
		if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
			return err
		}
	}

	if len(ua) > 0 {
		for k := 0; k < len(ua); k++ {
			ua[k].DocumentId = documentID
		}
		// Create document user access
		stmt := adua.INSERT(
			adua.DocumentID,
			adua.UserID,
			adua.Access,
		).
			MODELS(ua)
		if _, err := stmt.ExecContext(ctx, query.DB); err != nil {
			return err
		}
	}

	return nil
}
