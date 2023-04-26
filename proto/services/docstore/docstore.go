package docstore

import (
	context "context"
	"database/sql"
	"errors"

	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/htmlsanitizer"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	database "github.com/galexrt/fivenet/proto/resources/common/database"
	"github.com/galexrt/fivenet/proto/resources/documents"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DocsDefaultPageLimit  = 10
	DocShortContentLength = 200
)

var (
	user        = table.Users
	uCreator    = user.AS("creator")
	docs        = table.FivenetDocuments.AS("document")
	dUserAccess = table.FivenetDocumentsUserAccess.AS("user_access")
	dJobAccess  = table.FivenetDocumentsJobAccess.AS("job_access")
)

var (
	FailedQueryErr = status.Error(codes.Internal, "Failed to get/create/update documents!")
)

type Server struct {
	DocStoreServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *mstlystcdata.Enricher
	a  audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, aud audit.IAuditer) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
		a:  aud,
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
		jet.ProjectionList{jet.COUNT(docs.ID).AS("datacount.totalcount")},
		-1, userId, job, jobGrade)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, FailedQueryErr
	}

	resp := &FindDocumentsResponse{
		Pagination: database.EmptyPaginationResponseWithPageSize(req.Pagination.Offset, DocsDefaultPageLimit),
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := s.getDocumentsQuery(condition, nil,
		DocShortContentLength, userId, job, jobGrade).
		OFFSET(req.Pagination.Offset)

	if err := stmt.QueryContext(ctx, s.db, &resp.Documents); err != nil {
		return nil, FailedQueryErr
	}

	for i := 0; i < len(resp.Documents); i++ {
		s.c.EnrichJobInfo(resp.Documents[i].Creator)
	}

	database.PaginationHelperWithPageSize(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Documents),
		DocsDefaultPageLimit,
	)

	return resp, nil
}

func (s *Server) GetDocument(ctx context.Context, req *GetDocumentRequest) (*GetDocumentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	condition := docs.ID.EQ(jet.Uint64(req.DocumentId))

	countStmt := s.getDocumentsQuery(condition,
		jet.ProjectionList{jet.COUNT(docs.ID).AS("datacount.totalcount")},
		-1, userId, job, jobGrade)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, FailedQueryErr
	}

	resp := &GetDocumentResponse{}

	var err error
	resp.Document, err = s.getDocument(ctx, condition, userId, job, jobGrade)
	if err != nil {
		return nil, err
	}

	docAccess, err := s.GetDocumentAccess(ctx, &GetDocumentAccessRequest{
		DocumentId: resp.Document.Id,
	})
	if err != nil {
		return nil, FailedQueryErr
	}

	resp.Access = docAccess.Access

	return resp, nil
}

func (s *Server) getDocument(ctx context.Context, condition jet.BoolExpression, userId int32, job string, jobGrade int32) (*documents.Document, error) {
	var doc documents.Document

	stmt := s.getDocumentsQuery(condition, nil, -1, userId, job, jobGrade)

	if err := stmt.QueryContext(ctx, s.db, &doc); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, FailedQueryErr
		}
	}

	s.c.EnrichJobInfo(doc.Creator)

	return &doc, nil
}

func (s *Server) CreateDocument(ctx context.Context, req *CreateDocumentRequest) (*CreateDocumentResponse, error) {
	userId, job, _ := auth.GetUserInfoFromContext(ctx)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FailedQueryErr
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	docs := table.FivenetDocuments
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

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, FailedQueryErr
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, FailedQueryErr
	}

	if err := s.handleDocumentAccessChanges(ctx, tx, DOC_ACCESS_UPDATE_MODE_UPDATE, uint64(lastId), req.Access); err != nil {
		return nil, FailedQueryErr
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	return &CreateDocumentResponse{
		DocumentId: uint64(lastId),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, FailedQueryErr
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to edit this document!")
	}

	doc, err := s.getDocument(ctx, docs.ID.EQ(jet.Uint64(req.DocumentId)), userId, job, jobGrade)
	if err != nil {
		return nil, err
	}

	// Either the document is closed and the update request isn't re-opening the document
	if doc.GetClosed() && req.Closed != nil && !*req.Closed {
		return nil, status.Error(codes.Canceled, "Document is closed and can't be edited!")
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, FailedQueryErr
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
		WHERE(
			docs.ID.EQ(jet.Uint64(req.DocumentId)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, FailedQueryErr
	}

	if err := s.handleDocumentAccessChanges(ctx, tx, DOC_ACCESS_UPDATE_MODE_UPDATE, req.DocumentId, req.Access); err != nil {
		return nil, FailedQueryErr
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, FailedQueryErr
	}

	return &UpdateDocumentResponse{
		DocumentId: req.DocumentId,
	}, nil
}

func (s *Server) DeleteDocument(ctx context.Context, req *DeleteDocumentRequest) (*DeleteDocumentResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_EDIT)
	if err != nil {
		return nil, FailedQueryErr
	}
	if !check {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to delete this document!")
	}

	stmt := docs.
		UPDATE(
			docs.DeletedAt,
		).
		SET(
			docs.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			docs.ID.EQ(jet.Uint64(req.DocumentId)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return nil, nil
}
