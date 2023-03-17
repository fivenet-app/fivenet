package docstore

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/galexrt/arpanet/pkg/auth"
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
	dTemplates  = table.ArpanetDocumentsTemplates
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
}

func NewServer(db *sql.DB, p perms.Permissions) *Server {
	return &Server{
		db: db,
		p:  p,
	}
}

func (s *Server) FindDocuments(ctx context.Context, req *FindDocumentsRequest) (*FindDocumentsResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw("MATCH(title) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.Search}))
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
		Access:   &documents.DocumentAccess{},
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
	stmt := docs.INSERT(
		docs.Title,
		docs.Content,
		docs.ContentType,
		docs.Closed,
		docs.State,
		docs.Public,
		docs.CreatorID,
		docs.CategoryID,
	).VALUES(
		req.Title,
		htmlsanitizer.Sanitize(req.Content),
		documents.DOC_CONTENT_TYPE_HTML,
		req.Closed,
		req.State,
		req.Public,
		userId,
		job,
		req.CategoryId,
	)

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccessChanges(ctx, DOC_ACCESS_UPDATE_MODE_REPLACE, uint64(lastId), req.Access); err != nil {
		return nil, err
	}

	return &CreateDocumentResponse{
		Id: uint64(lastId),
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

	stmt := docs.UPDATE(
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

	return &UpdateDocumentResponse{}, nil
}

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	userId, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, req.DocumentId, userId, job, jobGrade, false, documents.DOC_ACCESS_ACCESS)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, status.Error(codes.PermissionDenied, "You don't have permission to view document access!")
	}

	stmt := docs.SELECT(
		dUserAccess.AllColumns,
		dJobAccess.AllColumns,
	).
		FROM(
			docs.
				LEFT_JOIN(dUserAccess,
					docs.ID.EQ(dUserAccess.DocumentID)).
				LEFT_JOIN(dJobAccess,
					docs.ID.EQ(dJobAccess.DocumentID)),
		).
		WHERE(
			jet.AND(
				docs.ID.EQ(jet.Uint64(req.DocumentId)),
				jet.OR(
					jet.AND(
						dUserAccess.Access.IS_NOT_NULL(),
						dUserAccess.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
					jet.AND(
						dUserAccess.Access.IS_NULL(),
						dJobAccess.Access.IS_NOT_NULL(),
						dJobAccess.Access.NOT_EQ(jet.Int32(int32(documents.DOC_ACCESS_BLOCKED))),
					),
				),
			),
		)

	resp := &GetDocumentAccessResponse{
		Access: &documents.DocumentAccess{},
	}
	if err := stmt.QueryContext(ctx, s.db, resp.Access); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) SetDocumentAccess(ctx context.Context, req *SetDocumentAccessRequest) (*SetDocumentAccessResponse, error) {
	resp := &SetDocumentAccessResponse{}

	if err := s.handleDocumentAccessChanges(ctx, req.Mode, req.DocumentId, req.Access); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) handleDocumentAccessChanges(ctx context.Context, mode DOC_ACCESS_UPDATE_MODE, documentID uint64, access *documents.DocumentAccess) error {
	userId := auth.GetUserIDFromContext(ctx)

	// Select existing job and user accesses
	var dest struct {
		jobs  []*documents.DocumentJobAccess
		users []*documents.DocumentUserAccess
	}
	selectStmt := jet.SELECT(
		dJobAccess.AllColumns,
		dUserAccess.AllColumns,
	).
		FROM(
			dJobAccess,
			dUserAccess,
		).
		WHERE(
			dJobAccess.DocumentID.EQ(jet.Uint64(documentID)),
		)

	if err := selectStmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// TODO add/update/remove for document access based on the current access in the database
	_ = dest.jobs
	_ = dest.users

	switch mode {
	case DOC_ACCESS_UPDATE_MODE_ADD:
		ja := access.Jobs
		// Create accesses
		if len(ja) > 0 {
			for k := 0; k < len(ja); k++ {
				ja[k].DocumentId = documentID
				ja[k].CreatorId = userId
			}

			// Create document job access
			stmt := dJobAccess.INSERT(
				dJobAccess.DocumentID,
				dJobAccess.Job,
				dJobAccess.MinimumGrade,
				dJobAccess.Access,
				dJobAccess.CreatorID,
			).
				MODELS(ja)

			fmt.Println(stmt.DebugSql())

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				return err
			}
		}

		ua := access.Users
		if len(ua) > 0 {
			for k := 0; k < len(ua); k++ {
				ua[k].DocumentId = documentID
				ua[k].CreatorId = userId
			}
			// Create document user access
			stmt := dUserAccess.INSERT(
				dUserAccess.DocumentID,
				dUserAccess.UserID,
				dUserAccess.Access,
				dUserAccess.CreatorID,
			).
				MODELS(ua)

			fmt.Println(stmt.DebugSql())

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				return err
			}
		}

	case DOC_ACCESS_UPDATE_MODE_REPLACE:
		// TODO

	case DOC_ACCESS_UPDATE_MODE_DELETE:
		// TODO

	case DOC_ACCESS_UPDATE_MODE_CLEAR:
		jobAccessStmt := dJobAccess.DELETE().
			WHERE(dJobAccess.DocumentID.EQ(jet.Uint64(documentID)))

		if _, err := jobAccessStmt.ExecContext(ctx, s.db); err != nil {
			return err
		}

		userAccessStmt := dUserAccess.DELETE().
			WHERE(dUserAccess.DocumentID.EQ(jet.Uint64(documentID)))

		if _, err := userAccessStmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}
