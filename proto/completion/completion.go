package completion

import (
	context "context"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/table"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		{Key: "dispatches", Name: "View"},
	})
}

var (
	j  = table.Jobs.AS("job")
	jg = table.JobGrades.AS("job_grade")
)

type Server struct {
	CompletionServiceServer

	cancel context.CancelFunc

	jobsCache *cache.Cache[string, *Job]
}

func NewServer() *Server {
	ctx, cancel := context.WithCancel(context.Background())

	jobsCache := cache.NewContext(
		ctx,
		cache.AsLFU[string, *Job](lfu.WithCapacity(32)),
		cache.WithJanitorInterval[string, *Job](120*time.Second),
	)

	s := &Server{
		cancel:    cancel,
		jobsCache: jobsCache,
	}

	s.refreshCache()

	return s
}

func (s *Server) refreshCache() error {
	var dest []*Job

	stmt := j.SELECT(
		j.Name,
		j.Label,
		jg.JobName.AS("job_name"),
		jg.Grade,
		jg.Label,
	).FROM(
		j.LEFT_JOIN(jg, jg.JobName.EQ(j.Name)),
	).
		ORDER_BY(
			j.Name.ASC(),
			jg.Grade.ASC(),
		)

	if err := stmt.Query(query.DB, &dest); err != nil {
		return err
	}

	// Update cache
	for _, job := range dest {
		s.jobsCache.Set(strings.ToLower(job.Name), job)
	}

	return nil
}

func (s *Server) CompleteJobNames(ctx context.Context, req *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error) {
	resp := &CompleteJobNamesResponse{}

	keys := s.jobsCache.Keys()

	for i := 0; i < len(keys); i++ {
		job, ok := s.jobsCache.Get(keys[i])
		if !ok {
			continue
		}

		// TODO use Bleve search in the future
		if strings.HasPrefix(job.Name, req.Search) || strings.Contains(job.Name, req.Search) {
			resp.Jobs = append(resp.Jobs, job)
		}
	}

	return resp, nil
}
func (s *Server) CompleteJobGrades(ctx context.Context, req *CompleteJobGradesRequest) (*CompleteJobGradesResponse, error) {
	resp := &CompleteJobGradesResponse{}

	job, ok := s.jobsCache.Get(strings.ToLower(req.Job))
	if !ok {
		return resp, nil
	}

	for _, g := range job.Grades {
		// TODO use Bleve search in the future
		if strings.HasPrefix(g.Label, req.Search) || strings.Contains(g.Label, req.Search) {
			resp.Grades = append(resp.Grades, g)
		}
	}

	return resp, nil
}
