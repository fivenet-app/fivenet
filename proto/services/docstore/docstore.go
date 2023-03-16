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
	u         = table.Users
	adt       = table.ArpanetDocumentsTemplates
	ad        = table.ArpanetDocuments.AS("document")
	adc       = table.ArpanetDocumentsComments
	adua      = table.ArpanetDocumentsUserAccess.AS("user_access")
	adja      = table.ArpanetDocumentsJobAccess.AS("job_access")
	dCategory = table.ArpanetDocumentsCategories.AS("category")
	adref     = table.ArpanetDocumentsReferences.AS("documentreference")
	adrel     = table.ArpanetDocumentsRelations.AS("documentrelation")
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
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw("MATCH(title) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.Search}))
	} else {
		condition = jet.Bool(true)
	}

	countStmt := s.getDocumentsQuery(
		condition,
		jet.ProjectionList{jet.COUNT(ad.ID).AS("total_count")},
		nil, userID, job, jobGrade)
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

	stmt := s.getDocumentsQuery(condition, nil, nil, userID, job, jobGrade)
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
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)

	condition := jet.OR(
		ad.ID.EQ(jet.Uint64(req.DocumentId)),
	)

	countStmt := s.getDocumentsQuery(condition, jet.ProjectionList{jet.COUNT(ad.ID).AS("total_count")}, nil, userID, job, jobGrade)
	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &GetDocumentResponse{
		Document: &documents.Document{},
		Access:   &documents.DocumentAccess{},
	}

	stmt := s.getDocumentsQuery(condition, nil, nil, userID, job, jobGrade)
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

	result, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.handleDocumentAccess(ctx, DOC_ACCESS_UPDATE_MODE_REPLACE, uint64(lastID), req.Access); err != nil {
		return nil, err
	}

	return &CreateDocumentResponse{
		Id: uint64(lastID),
	}, nil
}

func (s *Server) UpdateDocument(ctx context.Context, req *UpdateDocumentRequest) (*UpdateDocumentResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	check, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_EDIT)
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

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &UpdateDocumentResponse{}, nil
}

func (s *Server) GetDocumentAccess(ctx context.Context, req *GetDocumentAccessRequest) (*GetDocumentAccessResponse, error) {
	userID, job, jobGrade := auth.GetUserInfoFromContext(ctx)
	ok, err := s.checkIfUserHasAccessToDoc(ctx, userID, job, jobGrade, documents.DOC_ACCESS_ACCESS)
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
					ad.ID.EQ(adja.DocumentID)),
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

	if err := s.handleDocumentAccess(ctx, req.Mode, req.DocumentId, req.Access); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) handleDocumentAccess(ctx context.Context, mode DOC_ACCESS_UPDATE_MODE, documentID uint64, access *documents.DocumentAccess) error {
	userID := auth.GetUserIDFromContext(ctx)

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

	if err := selectStmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// TODO add/update/remove for document access based on the current access in the database

	ja := access.Jobs
	// Create accesses
	if len(ja) > 0 {
		for k := 0; k < len(ja); k++ {
			ja[k].DocumentId = documentID
			ja[k].CreatorId = userID
		}

		// Create document job access
		stmt := adja.INSERT(
			adja.DocumentID,
			adja.Job,
			adja.MinimumGrade,
			adja.Access,
			adja.CreatorID,
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
			ua[k].CreatorId = userID
		}
		// Create document user access
		stmt := adua.INSERT(
			adua.DocumentID,
			adua.UserID,
			adua.Access,
			adua.CreatorID,
		).
			MODELS(ua)

		fmt.Println(stmt.DebugSql())

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}
