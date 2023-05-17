package completor

import (
	context "context"
	"database/sql"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tUsers            = table.Users.AS("usershort")
	tDocumentCategory = table.FivenetDocumentsCategories.AS("document_category")
)

var (
	FailedSearchErr = status.Error(codes.Internal, "Failed to complete/ search the data!")
)

type Server struct {
	CompletorServiceServer

	db   *sql.DB
	p    perms.Permissions
	data *mstlystcdata.Cache
}

func NewServer(db *sql.DB, p perms.Permissions, data *mstlystcdata.Cache) *Server {
	s := &Server{
		db:   db,
		p:    p,
		data: data,
	}

	return s
}

func (s *Server) CompleteCitizens(ctx context.Context, req *CompleteCitizensRequest) (*CompleteCitizensRespoonse, error) {
	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	condition := jet.Bool(true)
	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
		condition = jet.CONCAT(tUsers.Firstname, jet.String(" "), tUsers.Lastname).
			LIKE(jet.String(req.Search))
	}

	stmt := tUsers.
		SELECT(
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(
			tUsers.Lastname.DESC(),
		).
		LIMIT(15)

	resp := &CompleteCitizensRespoonse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, FailedSearchErr
		}
	}

	return resp, nil
}

func (s *Server) CompleteJobs(ctx context.Context, req *CompleteJobsRequest) (*CompleteJobsResponse, error) {
	resp := &CompleteJobsResponse{}

	var search string
	if req.Search == nil || req.CurrentJob {
		userInfo := auth.MustGetUserInfoFromContext(ctx)
		search = userInfo.Job
	} else {
		search = *req.Search
	}

	var err error
	resp.Jobs, err = s.data.GetSearcher().SearchJobs(ctx, search, req.ExactMatch)
	if err != nil {
		return nil, FailedSearchErr
	}

	return resp, nil
}

func (s *Server) CompleteDocumentCategories(ctx context.Context, req *CompleteDocumentCategoriesRequest) (*CompleteDocumentCategoriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	jobsAttr, err := s.p.Attr(userInfo, CompletorServicePerm, CompletorServiceCompleteDocumentCategoriesPerm, CompletorServiceCompleteDocumentCategoriesJobsPermField)
	if err != nil {
		return nil, FailedSearchErr
	}
	var jobs perms.StringList
	if jobsAttr != nil {
		jobs = jobsAttr.([]string)
	}

	if len(jobs) == 0 {
		jobs = append(jobs, userInfo.Job)
	}

	req.Search = strings.TrimSpace(req.Search)
	req.Search = strings.ReplaceAll(req.Search, "%", "")
	req.Search = strings.ReplaceAll(req.Search, " ", "%")

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	condition := tDocumentCategory.Job.IN(jobsExp...)
	if req.Search != "" {
		req.Search = "%" + req.Search + "%"
		condition = condition.AND(
			tDocumentCategory.Name.LIKE(jet.String(req.Search)),
		)
	}

	stmt := tDocumentCategory.
		SELECT(
			tDocumentCategory.ID,
			tDocumentCategory.Name,
			tDocumentCategory.Description,
			tDocumentCategory.Job,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(tDocumentCategory).
		WHERE(condition).
		ORDER_BY(
			tDocumentCategory.Name.DESC(),
		).
		LIMIT(15)

	resp := &CompleteDocumentCategoriesResponse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Categories); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, FailedSearchErr
		}
	}

	return resp, nil
}
