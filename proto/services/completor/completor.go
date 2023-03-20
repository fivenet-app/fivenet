package completor

import (
	context "context"
	"database/sql"
	"errors"
	"strings"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/complhelper"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	us = table.Users.AS("usershort")
)

type Server struct {
	CompletorServiceServer

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

func (s *Server) CompleteCharNames(ctx context.Context, req *CompleteCharNamesRequest) (*CompleteCharNamesRespoonse, error) {
	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw("MATCH(firstname,lastname) AGAINST ($search IN NATURAL LANGUAGE MODE)", jet.RawArgs{"$search": req.Search}))
	} else {
		condition = jet.Bool(true)
	}

	stmt := us.
		SELECT(
			us.ID,
			us.Identifier,
			us.Firstname,
			us.Lastname,
			us.Job,
			us.JobGrade,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname")).
		FROM(us).
		WHERE(condition).
		LIMIT(15)

	resp := &CompleteCharNamesRespoonse{}
	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) CompleteJobNames(ctx context.Context, req *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error) {
	resp := &CompleteJobNamesResponse{}

	req.Search = strings.ToLower(req.Search)

	keys := s.c.Jobs.Keys()
	for i := 0; i < len(keys); i++ {
		job, ok := s.c.Jobs.Get(keys[i])
		if !ok {
			continue
		}

		if strings.HasPrefix(strings.ToLower(job.Label), req.Search) || strings.Contains(strings.ToLower(job.Label), req.Search) {
			resp.Jobs = append(resp.Jobs, job)
		} else if strings.HasPrefix(job.Name, req.Search) || strings.Contains(job.Name, req.Search) {
			resp.Jobs = append(resp.Jobs, job)
		}

		if i > 10 {
			break
		}
	}

	return resp, nil
}

func (s *Server) CompleteDocumentCategory(ctx context.Context, req *CompleteDocumentCategoryRequest) (*CompleteDocumentCategoryResponse, error) {
	userId := auth.GetUserIDFromContext(ctx)

	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, CompletorServicePermKey+"-CompleteDocumentCategory")
	if err != nil {
		return nil, err
	}

	resp := &CompleteDocumentCategoryResponse{}
	if len(jobs) == 0 {
		return resp, nil
	}

	req.Search = strings.ToLower(req.Search)

	for _, j := range jobs {
		c, ok := s.c.DocCategories.Get(j)
		if !ok {
			continue
		}

		for _, v := range c {
			if strings.HasPrefix(strings.ToLower(v.Name), req.Search) || strings.Contains(strings.ToLower(v.Name), req.Search) {
				resp.Categories = append(resp.Categories, v)
			}
		}
	}

	return resp, nil
}
