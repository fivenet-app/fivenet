package completor

import (
	context "context"
	"database/sql"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lfu"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/documents"
	"github.com/galexrt/arpanet/proto/resources/jobs"
	"github.com/galexrt/arpanet/query/arpanet/table"
)

var (
	j   = table.Jobs.AS("job")
	jg  = table.JobGrades.AS("job_grade")
	adc = table.ArpanetDocumentsCategories
)

type Server struct {
	CompletorServiceServer

	db *sql.DB
	p  perms.Permissions

	// Cache related
	cancel              context.CancelFunc
	jobsCache           *cache.Cache[string, *jobs.Job]
	docsCategoriesCache *cache.Cache[string, []*documents.DocumentCategory]
}

func NewServer(db *sql.DB, p perms.Permissions) *Server {
	ctx, cancel := context.WithCancel(context.Background())

	jobsCache := cache.NewContext(
		ctx,
		cache.AsLFU[string, *jobs.Job](lfu.WithCapacity(32)),
		cache.WithJanitorInterval[string, *jobs.Job](120*time.Second),
	)

	docsCategoriesCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*documents.DocumentCategory](lru.WithCapacity(32)),
	)

	s := &Server{
		db: db,
		p:  p,

		cancel:              cancel,
		jobsCache:           jobsCache,
		docsCategoriesCache: docsCategoriesCache,
	}

	s.refreshCache()

	return s
}

func (s *Server) refreshCache() error {
	if err := s.refreshJobsCache(); err != nil {
		return err
	}

	if err := s.refreshDocumentCategories(); err != nil {
		return err
	}

	return nil
}

func (s *Server) refreshJobsCache() error {
	var dest []*jobs.Job

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
	if err := stmt.Query(s.db, &dest); err != nil {
		return err
	}

	// Update cache
	for _, job := range dest {
		s.jobsCache.Set(strings.ToLower(job.Name), job)
	}

	return nil
}

func (s *Server) refreshDocumentCategories() error {
	var dest []*documents.DocumentCategory

	stmt := adc.SELECT(
		adc.AllColumns,
	).
		FROM(adc).
		GROUP_BY(adc.Job).
		ORDER_BY(adc.Name.ASC())
	if err := stmt.Query(s.db, &dest); err != nil {
		return err
	}

	categoriesPerJob := map[string][]*documents.DocumentCategory{}
	for _, c := range dest {
		if _, ok := categoriesPerJob[c.Job]; !ok {
			categoriesPerJob[c.Job] = []*documents.DocumentCategory{}
		}
		categoriesPerJob[c.Job] = append(categoriesPerJob[c.Job], c)
	}

	// Update cache
	for job, cs := range categoriesPerJob {
		s.docsCategoriesCache.Set(job, cs)
	}

	return nil
}

// TODO use Bleve search in the future
func (s *Server) CompleteJobNames(ctx context.Context, req *CompleteJobNamesRequest) (*CompleteJobNamesResponse, error) {
	resp := &CompleteJobNamesResponse{}
	keys := s.jobsCache.Keys()
	for i := 0; i < len(keys); i++ {
		job, ok := s.jobsCache.Get(keys[i])
		if !ok {
			continue
		}

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
		if strings.HasPrefix(g.Label, req.Search) || strings.Contains(g.Label, req.Search) {
			resp.Grades = append(resp.Grades, g)
		}
	}

	return resp, nil
}

func (s *Server) CompleteDocumentCategory(ctx context.Context, req *CompleteDocumentCategoryRequest) (*CompleteDocumentCategoryResponse, error) {
	userID := auth.GetUserIDFromContext(ctx)

	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userID, CompletorServicePermKey+"-CompleteDocumentCategory")
	if err != nil {
		return nil, err
	}

	resp := &CompleteDocumentCategoryResponse{}
	if len(jobs) == 0 {
		return resp, nil
	}

	for _, j := range jobs {
		c, ok := s.docsCategoriesCache.Get(j)
		if !ok {
			continue
		}

		for _, v := range c {
			if strings.HasPrefix(v.Name, req.Search) || strings.Contains(v.Name, req.Search) {
				resp.Categories = append(resp.Categories, v)
			}
		}
	}

	return resp, nil
}
