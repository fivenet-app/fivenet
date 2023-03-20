package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/complhelper"
	"github.com/galexrt/arpanet/pkg/htmlsanitizer"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	u           = table.Users
	docs        = table.ArpanetDocuments.AS("document")
	dComments   = table.ArpanetDocumentsComments
	dUserAccess = table.ArpanetDocumentsUserAccess.AS("user_access")
	dJobAccess  = table.ArpanetDocumentsJobAccess.AS("job_access")
	dCategory   = table.ArpanetDocumentsCategories.AS("category")
)

type Server struct {
	DocStoreServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *complhelper.Completor
}

func NewServer(db *sql.DB, p perms.Permissions, c *complhelper.Completor) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
	}
}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw(
			"MATCH(title) AGAINST ($search IN NATURAL LANGUAGE MODE)",
			jet.RawArgs{"$search": req.Search},
		))
	} else {
		condition = jet.Bool(true)
	}

	countStmt := s.getDocumentsQuery(
		condition,
		jet.ProjectionList{jet.COUNT(docs.ID).AS("total_count")},
		nil, userId, job, jobGrade)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
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

	stmt := s.getDocumentsQuery(condition, nil, nil, userId, job, jobGrade)
	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
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
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	condition := jet.OR(
		docs.ID.EQ(jet.Uint64(req.DocumentId)),
	)

	countStmt := s.getDocumentsQuery(condition, jet.ProjectionList{jet.COUNT(docs.ID).AS("total_count")}, nil, userId, job, jobGrade)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &GetDocumentResponse{
		Document: &documents.Document{},
	}

	stmt := s.getDocumentsQuery(condition, nil, nil, userId, job, jobGrade)
	if err := stmt.QueryContext(ctx, s.db, resp.Document); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	docAccess, err := s.GetDocumentAccess(ctx, &GetDocumentAccessRequest{
		DocumentId: resp.Document.Id,
	})
	if err != nil {
		return nil, err
	}

	resp.Access = docAccess.Access

	return resp, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	docs := table.ArpanetDocuments
	stmt := docs.
		INSERT(
			docs.CategoryID,
			docs.Title,
			docs.Content,
			docs.ContentType,
			docs.Data,
			docs.CreatorID,
			docs.CreatorJob,
			docs.State,
			docs.Closed,
			docs.Public,
		).
		VALUES(
			req.CategoryId,
			req.Title,
			htmlsanitizer.Sanitize(req.Content),
			documents.DOC_CONTENT_TYPE_HTML,
			req.Data,
			userId,
			job,
			req.State,
			req.Closed,
			req.Public,
		)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccessChanges(ctx, DOC_ACCESS_UPDATE_MODE_UPDATE, uint64(lastId), req.Access); err != nil {
		return nil, err
	}

	return &CreateDocumentResponse{
		DocumentId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, err
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this document!")
	}

	stmt := docs.
		UPDATE(
			docs.Title,
			docs.Content,
			docs.Closed,
			docs.State,
			docs.Public,
		).
		SET(
			req.Title,
			req.Content,
			req.Closed,
			req.State,
			req.Public,
		).
		WHERE(docs.ID.EQ(jet.Uint64(req.DocumentId)))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccessChanges(ctx, DOC_ACCESS_UPDATE_MODE_UPDATE, req.DocumentId, req.Access); err != nil {
		return nil, err
	}

	return &UpdateDocumentResponse{
		DocumentId: req.DocumentId,
	}, nil
}
