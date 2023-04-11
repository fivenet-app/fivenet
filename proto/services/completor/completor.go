package completor

import (
	context "context"
	"database/sql"
	"errors"
	"strings"

	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	user = table.Users.AS("usershort")
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

func (s *Server) CompleteCharNames(ctx context.Context, req *CompleteCharNamesRequest) (*CompleteCharNamesRespoonse, error) {
	req.Search = strings.ToLower(strings.TrimSpace(req.Search))

	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw("MATCH(firstname,lastname) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.Search}))
	} else {
		condition = jet.Bool(true)
	}

	stmt := user.
		SELECT(
			user.ID,
			user.Identifier,
			user.Firstname,
			user.Lastname,
			user.Job,
			user.JobGrade,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(user).
		WHERE(condition).
		LIMIT(15)

	resp := &CompleteCharNamesRespoonse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, FailedSearchErr
		}
	}

	return resp, nil
}

func (s *Server) CompleteJobNames(ctx context.Context, req *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error) {
	resp := &CompleteJobNamesResponse{}

	var search string
	if req.Search == nil || req.CurrentJob {
		_, job, _ := auth.GetUserInfoFromContext(ctx)
		search = job
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

func (s *Server) CompleteDocumentCategory(ctx context.Context, req *CompleteDocumentCategoryRequest) (*CompleteDocumentCategoryResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, CompletorServicePermKey+"-CompleteDocumentCategory")
	if err != nil {
		return nil, FailedSearchErr
	}

	resp := &CompleteDocumentCategoryResponse{}
	if len(jobs) == 0 {
		return resp, nil
	}

	resp.Categories, err = s.data.GetSearcher().
		SearchDocumentCategories(ctx, req.Search, jobs)
	if err != nil {
		return nil, FailedSearchErr
	}

	return resp, nil
}
